package queue

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
	"github.com/voxpupuli/webhook-go/config"
	"github.com/voxpupuli/webhook-go/lib/helpers"
)

// Queue holds the items to be processed and manages job concurrency.

type Queue struct {
	Items   []*QueueItem    // List of items in the queue
	wg      sync.WaitGroup  // WaitGroup to synchronize goroutines
	jobChan chan *QueueItem // Channel for jobs to be processed
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

var queue = &Queue{}

// GetQueueItems returns all items currently in the queue.
func GetQueueItems() []*QueueItem {
	return queue.Items
}

// AddToQueue adds a new command to the queue and triggers job processing.
// Returns the newly created QueueItem or an error if the job cannot be queued.
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

	queue.Items = append(queue.Items, &queueItem)

	if !queueJob(&queueItem) {
		queueItem.State = "queue full"
		return &queueItem, fmt.Errorf("queue is full")
	}

	trimItems()

	return &queueItem, nil
}

// trimItems ensures the queue doesn't exceed the configured maximum history size.
func trimItems() {
	conf := config.GetConfig()

	if conf.Server.Queue.MaxHistoryItems == 0 {
		return
	}

	if len(queue.Items) > conf.Server.Queue.MaxHistoryItems {
		queue.Items = queue.Items[:conf.Server.Queue.MaxHistoryItems]
	}
}

// Work initializes the job queue and starts the worker process.
func Work() error {
	conf := config.GetConfig()
	log.Printf("start queue with %d jobs", conf.Server.Queue.MaxConcurrentJobs)

	queue.jobChan = make(chan *QueueItem, conf.Server.Queue.MaxConcurrentJobs)
	queue.wg.Add(1)
	queue.Items = []*QueueItem{}

	go worker()
	return nil
}

// Dispose shuts down the job channel and signals the worker to stop.
func Dispose() {
	close(queue.jobChan)
}

// queueJob attempts to add a job to the queue.
// Returns true if the job was successfully queued, false otherwise.
func queueJob(command *QueueItem) bool {
	select {
	case queue.jobChan <- command:
		return true
	default:
		return false
	}
}

// worker processes jobs from the queue and executes the associated commands.
func worker() {
	defer queue.wg.Done()

	log.Println("Worker is waiting for jobs")

	conf := config.GetConfig()
	conn := helpers.ChatopsSetup()

	for job := range queue.jobChan {
		log.Println("Worker picked Job", job.Id)
		job.StartedAt = time.Now()

		res, err := helpers.Execute(job.Command)
		job.Response = res

		job.FinishedAt = time.Now()
		if err != nil {
			log.Errorf("failed to execute local command `%s` with error: `%s` `%s`", job.Command, err, res)

			if conf.ChatOps.Enabled {
				conn.PostMessage(http.StatusInternalServerError, job.Name, res)
			}
			job.State = "failed"
			continue
		}

		conn.PostMessage(http.StatusAccepted, job.Name, res)
		job.State = "success"
	}
}
