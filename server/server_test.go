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

func TestPingRoute(t *testing.T) {
	router := NewRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"message\":\"running\"}", w.Body.String())
}

func TestQueue(t *testing.T) {
	mCfg := "../lib/helpers/yaml/webhook.queue.yaml"
	config.Init(&mCfg)

	queue.Work()

	router := NewRouter()

	payloadFile, err := os.Open("../lib/parsers/json/github/push.json")
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/r10k/environment", payloadFile)
	req.Header.Add("X-GitHub-Event", "push")

	router.ServeHTTP(w, req)

	var queueItem queue.QueueItem
	err = json.Unmarshal(w.Body.Bytes(), &queueItem)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 202, w.Code)
	assert.Equal(t, "simple-tag", queueItem.Name)
	assert.Equal(t, "added", queueItem.State)
}
