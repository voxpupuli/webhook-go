package queue

import (
	"errors"
	"net/http"
	"slices"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
	"github.com/voxpupuli/webhook-go/config"
	"github.com/voxpupuli/webhook-go/lib/helpers"
)

// Queue holds the items to be processed and manages job concurrency.

type Queue struct {
	items   []*QueueItem    // History of queued items, guarded by mu
	jobChan chan *QueueItem // Channel for jobs to be processed, guarded by ctrlMu
	wg      sync.WaitGroup  // WaitGroup tracking the worker goroutine

	mu        sync.RWMutex // Guards items and all mutations of QueueItem fields
	ctrlMu    sync.Mutex   // Guards sends/closes on jobChan and the started/disposed flags
	lifecycle sync.Mutex   // Serializes Work() and Dispose() so they can't interleave

	started  bool // Whether a worker goroutine is currently running
	disposed bool // Whether the queue has been shut down
}

// QueueItem represents a task to be executed, with metadata like ID, command, and state.

type QueueItem struct {
	Id          uuid.UUID   // Unique identifier for the job
	Name        string      // Descriptive name of the job
	CommandType string      // Type of command being executed
	AddedAt     time.Time   // Timestamp when the job was added to the queue
	StartedAt   time.Time   // Timestamp when the job execution started
	FinishedAt  time.Time   // Timestamp when the job execution finished
	Command     []string    // Command to be executed
	Response    interface{} // Response from the executed command
	State       string      // Current state of the job (e.g., added, success, failed)
}

var (
	queue = &Queue{}

	errQueueFull    = errors.New("queue is full")
	errShuttingDown = errors.New("queue is shutting down")
)

// GetQueueItems returns a snapshot of all items currently in the queue.
// The returned items are detached copies, so they can safely be read or
// marshaled while the worker keeps updating the live items.
func GetQueueItems() []*QueueItem {
	queue.mu.RLock()
	defer queue.mu.RUnlock()

	items := make([]*QueueItem, len(queue.items))
	for i, item := range queue.items {
		items[i] = copyItem(item)
	}

	return items
}

// AddToQueue adds a new command to the queue and triggers job processing.
// Returns a snapshot of the newly created QueueItem as it was at queuing
// time, or an error if the job cannot be queued. The live item keeps being
// mutated by the worker; use GetQueueItems to follow its progress.
func AddToQueue(commandType string, name string, command []string) (*QueueItem, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	queueItem := QueueItem{
		Id:          id,
		AddedAt:     time.Now(),
		Command:     command,
		CommandType: commandType,
		Name:        name,
		State:       "added",
	}

	queue.mu.Lock()
	queue.items = append(queue.items, &queueItem)
	queue.mu.Unlock()

	// The snapshot returned to the caller is taken before the job is handed
	// to the worker, so it always reflects the state at queuing time and
	// can never race with the worker's updates to the live item.
	result := snapshot(&queueItem)

	if err := queueJob(&queueItem); err != nil {
		state := "queue full"
		if errors.Is(err, errShuttingDown) {
			state = "shutting down"
		}

		setState(&queueItem, state)
		trimItems()

		// The job never reached the worker, so it is safe to take another
		// snapshot that reflects why it wasn't queued.
		return snapshot(&queueItem), err
	}

	trimItems()

	return result, nil
}

// snapshot returns a detached copy of item, taken under the queue's read
// lock so it can't race with concurrent updates to the live item.
func snapshot(item *QueueItem) *QueueItem {
	queue.mu.RLock()
	defer queue.mu.RUnlock()

	return copyItem(item)
}

// copyItem returns a detached copy of item. The caller must hold mu or
// otherwise ensure the item is not being mutated concurrently.
func copyItem(item *QueueItem) *QueueItem {
	cp := *item
	cp.Command = slices.Clone(item.Command)

	return &cp
}

// trimItems ensures the queue doesn't exceed the configured maximum history
// size. When the limit is exceeded the oldest items are dropped so the most
// recent history is retained.
func trimItems() {
	conf := config.GetConfig()

	max := conf.Server.Queue.MaxHistoryItems
	if max <= 0 {
		return
	}

	queue.mu.Lock()
	defer queue.mu.Unlock()

	if len(queue.items) > max {
		queue.items = append(queue.items[:0], queue.items[len(queue.items)-max:]...)
	}
}

// Work initializes the job queue and starts the worker process.
// It is idempotent: calling it again while a worker is already running
// is a no-op instead of replacing the job channel and leaking a worker.
func Work() error {
	conf := config.GetConfig()

	queue.lifecycle.Lock()
	defer queue.lifecycle.Unlock()

	queue.ctrlMu.Lock()
	if queue.started {
		queue.ctrlMu.Unlock()
		return nil
	}
	queue.started = true
	queue.disposed = false
	queue.jobChan = make(chan *QueueItem, conf.Server.Queue.MaxConcurrentJobs)
	jobChan := queue.jobChan
	queue.ctrlMu.Unlock()

	queue.mu.Lock()
	queue.items = nil
	queue.mu.Unlock()

	queue.wg.Add(1)

	log.Printf("start queue with %d jobs", conf.Server.Queue.MaxConcurrentJobs)

	go worker(jobChan)

	return nil
}

// Dispose shuts down the queue: it stops accepting new jobs, closes the job
// channel, and waits for the worker to finish all in-flight and buffered
// jobs. It is safe to call Dispose multiple times; subsequent calls just
// re-check that everything has drained.
func Dispose() {
	queue.lifecycle.Lock()
	defer queue.lifecycle.Unlock()

	queue.ctrlMu.Lock()
	if !queue.started {
		queue.ctrlMu.Unlock()
		queue.wg.Wait()
		return
	}
	queue.started = false
	queue.disposed = true
	close(queue.jobChan)
	queue.ctrlMu.Unlock()

	queue.wg.Wait()
}

// queueJob attempts to add a job to the job channel.
// The send happens under ctrlMu so it can never race with Dispose closing
// the channel. Returns an error if the queue is full or shutting down.
func queueJob(command *QueueItem) error {
	queue.ctrlMu.Lock()
	defer queue.ctrlMu.Unlock()

	if queue.disposed {
		return errShuttingDown
	}

	select {
	case queue.jobChan <- command:
		return nil
	default:
		return errQueueFull
	}
}

// setState updates a job's state under the queue lock.
func setState(job *QueueItem, state string) {
	queue.mu.Lock()
	defer queue.mu.Unlock()

	job.State = state
}

// markStarted records when a job's execution started, under the queue lock.
func markStarted(job *QueueItem) {
	queue.mu.Lock()
	defer queue.mu.Unlock()

	job.StartedAt = time.Now()
}

// markFinished records a job's result and final state, under the queue lock.
func markFinished(job *QueueItem, res interface{}, state string) {
	queue.mu.Lock()
	defer queue.mu.Unlock()

	job.Response = res
	job.FinishedAt = time.Now()
	job.State = state
}

// worker processes jobs from the queue and executes the associated commands.
func worker(jobChan chan *QueueItem) {
	defer queue.wg.Done()

	log.Println("Worker is waiting for jobs")

	conf := config.GetConfig()
	conn := helpers.ChatopsSetup()

	for job := range jobChan {
		log.Println("Worker picked Job", job.Id)
		markStarted(job)

		res, err := helpers.Execute(job.Command)

		if err != nil {
			log.Errorf("failed to execute local command `%s` with error: `%s` `%s`", job.Command, err, res)

			if conf.ChatOps.Enabled {
				conn.PostMessage(http.StatusInternalServerError, job.Name, res)
			}
			markFinished(job, res, "failed")
			continue
		}

		conn.PostMessage(http.StatusAccepted, job.Name, res)
		markFinished(job, res, "success")
	}
}
