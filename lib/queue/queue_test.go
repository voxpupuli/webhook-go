package queue

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/voxpupuli/webhook-go/config"
)

// initQueueConfig loads the queue-enabled test configuration
// (server.queue.enabled: true, max_concurrent_jobs: 10,
// max_history_items: 20) used by all queue tests.
func initQueueConfig(t *testing.T) {
	t.Helper()

	mCfg := "../helpers/yaml/webhook.queue.yaml"
	config.Init(&mCfg)
}

// TestAddToQueueConcurrent hammers the queue from many goroutines at once:
// concurrent AddToQueue calls, concurrent reads of the queue history, and
// marshaling of the item returned by AddToQueue while the worker processes
// it. With the old implementation this triggered the race detector on the
// unsynchronized Items slice and the worker's unsynchronized QueueItem
// field writes.
func TestAddToQueueConcurrent(t *testing.T) {
	initQueueConfig(t)

	if err := Work(); err != nil {
		t.Fatalf("failed to start queue: %v", err)
	}
	defer Dispose()

	var wg sync.WaitGroup

	// Marshal the item returned by AddToQueue while the worker mutates the
	// live item. The returned item must be a detached snapshot, otherwise
	// this races with the worker's updates.
	queued := make(chan *QueueItem, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()

		item, err := AddToQueue("module", "race-live-item", []string{"sleep", "0.5"})
		queued <- item
		if err != nil {
			t.Errorf("failed to add job: %v", err)
			return
		}
		if item.State != "added" {
			t.Errorf("expected state %q on the returned snapshot, got %q", "added", item.State)
		}

		// Keep marshaling the returned item until the worker finished the
		// job (or a deadline passed), guaranteeing temporal overlap with
		// the worker's writes to the live item.
		deadline := time.Now().Add(2 * time.Second)
		for time.Now().Before(deadline) {
			if _, err := json.Marshal(item); err != nil {
				t.Errorf("failed to marshal queue item: %v", err)
				return
			}
			for _, it := range GetQueueItems() {
				if it.Id == item.Id && (it.State == "success" || it.State == "failed") {
					return
				}
			}
		}
	}()

	// Wait until the long-running job is queued so the adder goroutines
	// below can't starve it out of the bounded job channel.
	if item := <-queued; item == nil {
		t.Fatal("expected the race-live-item job to be queued")
	}

	// Add jobs from many goroutines concurrently; with the old
	// implementation the unsynchronized appends to Items raced each other.
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			if _, err := AddToQueue("module", fmt.Sprintf("race-test-%d", i), []string{"sleep", "0.1"}); err != nil {
				// A full queue is expected here; the job is still recorded
				// in the history.
				return
			}
		}(i)
	}

	// Read and marshal the queue history concurrently with the additions
	// and the worker's item updates.
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for j := 0; j < 400; j++ {
				if _, err := json.Marshal(GetQueueItems()); err != nil {
					t.Errorf("failed to marshal queue items: %v", err)
					return
				}
				time.Sleep(5 * time.Millisecond)
			}
		}()
	}

	wg.Wait()
}

// TestTrimItemsKeepsNewest verifies that trimming the history keeps the
// most recent items instead of the oldest ones.
func TestTrimItemsKeepsNewest(t *testing.T) {
	initQueueConfig(t)

	if err := Work(); err != nil {
		t.Fatalf("failed to start queue: %v", err)
	}
	defer Dispose()

	// webhook.queue.yaml sets max_history_items to 20.
	const total = 30
	const maxHistoryItems = 20

	for i := 0; i < total; i++ {
		// A full queue is fine: the item is still recorded in the history.
		_, _ = AddToQueue("module", fmt.Sprintf("trim-test-%02d", i), []string{"true"})
	}

	items := GetQueueItems()
	if len(items) != maxHistoryItems {
		t.Fatalf("expected %d items after trimming, got %d", maxHistoryItems, len(items))
	}
	if items[0].Name != "trim-test-10" {
		t.Errorf("expected the oldest retained item to be %q, got %q", "trim-test-10", items[0].Name)
	}
	if items[len(items)-1].Name != "trim-test-29" {
		t.Errorf("expected the newest item to be %q, got %q", "trim-test-29", items[len(items)-1].Name)
	}
}

// TestWorkIsIdempotent verifies that a second Work call neither replaces
// the job channel (leaking the first worker goroutine) nor resets the
// queue history.
func TestWorkIsIdempotent(t *testing.T) {
	initQueueConfig(t)

	if err := Work(); err != nil {
		t.Fatalf("failed to start queue: %v", err)
	}

	item, err := AddToQueue("module", "work-idempotent", []string{"true"})
	if err != nil {
		t.Fatalf("failed to add job: %v", err)
	}

	if err := Work(); err != nil {
		t.Fatalf("failed to call Work a second time: %v", err)
	}

	items := GetQueueItems()
	if len(items) != 1 {
		t.Fatalf("expected the history to be preserved, got %d items", len(items))
	}
	if items[0].Id != item.Id {
		t.Errorf("expected the queued job %s to still be in the history, got %s", item.Id, items[0].Id)
	}

	Dispose()
}

// TestDisposeDrainsInFlightJobs verifies that Dispose waits for jobs that
// are still being processed (or still buffered) to finish.
func TestDisposeDrainsInFlightJobs(t *testing.T) {
	initQueueConfig(t)

	if err := Work(); err != nil {
		t.Fatalf("failed to start queue: %v", err)
	}

	item, err := AddToQueue("module", "drain-test", []string{"sleep", "1"})
	if err != nil {
		t.Fatalf("failed to add job: %v", err)
	}

	// Dispose must block until the in-flight job has finished.
	Dispose()

	var final *QueueItem
	for _, it := range GetQueueItems() {
		if it.Id == item.Id {
			final = it
		}
	}
	if final == nil {
		t.Fatal("expected the job to be recorded in the history")
	}
	if final.State != "success" {
		t.Errorf("expected state %q after Dispose drained the queue, got %q", "success", final.State)
	}
	if final.FinishedAt.IsZero() {
		t.Error("expected FinishedAt to be set after Dispose drained the queue")
	}
}

// TestAddToQueueAfterDisposeReturnsError verifies that adding a job after
// Dispose returns a clean error (and a state snapshot) instead of
// panicking with a send on a closed channel.
func TestAddToQueueAfterDisposeReturnsError(t *testing.T) {
	initQueueConfig(t)

	if err := Work(); err != nil {
		t.Fatalf("failed to start queue: %v", err)
	}

	Dispose()

	// A second Dispose must also be safe.
	Dispose()

	item, err := AddToQueue("module", "after-dispose", []string{"true"})
	if err == nil {
		t.Fatal("expected an error when adding a job to a disposed queue")
	}
	if item == nil {
		t.Fatal("expected a queue item snapshot even when the job was not queued")
	}
	if item.State != "shutting down" {
		t.Errorf("expected state %q when adding a job to a disposed queue, got %q", "shutting down", item.State)
	}
}
