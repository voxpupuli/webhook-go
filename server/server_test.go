package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/voxpupuli/webhook-go/config"
	"github.com/voxpupuli/webhook-go/lib/queue"
)

// TestPingRoute ensures the /health endpoint returns HTTP 200 with "running" message.
func TestPingRoute(t *testing.T) {
	router := NewRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)                                  // Check response status code
	assert.Equal(t, "{\"message\":\"running\"}", w.Body.String()) // Check response body
}

// TestQueue verifies if a payload is correctly processed and added to the queue.
func TestQueue(t *testing.T) {
	// Initialize the config with the webhook queue config file
	mCfg := "../lib/helpers/yaml/webhook.queue.yaml"
	config.Init(&mCfg)

	// Start the queue worker
	queue.Work()

	router := NewRouter()

	// Open the test payload file
	payloadFile, err := os.Open("../lib/parsers/json/github/push.json")
	if err != nil {
		t.Fatal(err) // Fail if unable to open the file
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/r10k/environment", payloadFile)
	req.Header.Add("X-GitHub-Event", "push") // Set GitHub event header

	router.ServeHTTP(w, req)

	var queueItem queue.QueueItem
	err = json.Unmarshal(w.Body.Bytes(), &queueItem)
	if err != nil {
		t.Fatal(err) // Fail if JSON unmarshaling fails
	}

	assert.Equal(t, 202, w.Code)                  // Ensure 202 Accepted
	assert.Equal(t, "simple-tag", queueItem.Name) // Ensure correct queue item name
	assert.Equal(t, "added", queueItem.State)     // Ensure correct queue item state
}
