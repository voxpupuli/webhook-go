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

type Queue struct {
	Items   []*QueueItem
	wg      sync.WaitGroup
	jobChan chan *QueueItem
}

type QueueItem struct {
	Id          uuid.UUID
	Name        string
	CommandType string
	AddedAt     time.Time
	StartedAt   time.Time
	FinishedAt  time.Time
	Command     []string
	Response    interface{}
	State       string
}

var queue = &Queue{}

func GetQueueItems() []*QueueItem {
	return queue.Items
}

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

func trimItems() {
	conf := config.GetConfig()

	if conf.Server.Queue.MaxHistoryItems == 0 {
		return
	}

	if len(queue.Items) > conf.Server.Queue.MaxHistoryItems {
		queue.Items = queue.Items[:conf.Server.Queue.MaxHistoryItems]
	}
}

func Work() error {
	conf := config.GetConfig()
	log.Printf("start queue with %d jobs", conf.Server.Queue.MaxConcurrentJobs)

	queue.jobChan = make(chan *QueueItem, conf.Server.Queue.MaxConcurrentJobs)
	queue.wg.Add(1)
	queue.Items = []*QueueItem{}

	go worker()
	return nil
}

func Dispose() {
	close(queue.jobChan)
}

func queueJob(command *QueueItem) bool {
	select {
	case queue.jobChan <- command:
		return true
	default:
		return false
	}
}

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
				conn.PostMessage(http.StatusInternalServerError, job.Name)
			}
			job.State = "failed"
			continue
		}

		job.State = "success"
	}
}
