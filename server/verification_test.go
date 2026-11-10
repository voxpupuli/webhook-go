package server

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/voxpupuli/webhook-go/config"
	"github.com/voxpupuli/webhook-go/lib/queue"
)

const verificationTestSecret = "github-test-secret"

// initVerificationConfig loads a config with the queue enabled (so
// successful deploys never run r10k) and a GitHub verification secret.
func initVerificationConfig(t *testing.T) {
	t.Helper()

	configYAML := `---
r10k:
  command_path: /nonexistent/test/webhook-go-r10k
server:
  queue:
    enabled: true
    max_concurrent_jobs: 10
    max_history_items: 20
secrets:
  github: ` + verificationTestSecret + "\n"

	path := filepath.Join(t.TempDir(), "webhook.yml")
	if err := os.WriteFile(path, []byte(configYAML), 0o600); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	config.Init(&path)
}

// TestWebhookVerificationRoute verifies through the full router that a
// configured GitHub secret rejects unsigned webhooks with 401 while
// correctly signed ones are still processed, and that webhooks claiming
// an unconfigured provider cannot bypass verification.
func TestWebhookVerificationRoute(t *testing.T) {
	initVerificationConfig(t)

	if err := queue.Work(); err != nil {
		t.Fatalf("failed to start the queue: %v", err)
	}

	router := NewRouter()

	payload, err := os.ReadFile("../lib/parsers/json/github/push_not_deleted.json")
	if err != nil {
		t.Fatal(err)
	}

	mac := hmac.New(sha256.New, []byte(verificationTestSecret))
	_, _ = mac.Write(payload)
	signature := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	t.Run("unsigned webhook rejected with 401", func(t *testing.T) {
		for _, endpoint := range []string{"/api/v1/r10k/environment", "/api/v1/r10k/module"} {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(payload))
			req.Header.Add("X-GitHub-Event", "push")
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusUnauthorized, w.Code)

			var body map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
				t.Fatalf("failed to unmarshal response body: %v", err)
			}
			assert.Equal(t, "Webhook verification failed", body["message"])
		}
	})

	t.Run("webhook with invalid signature rejected with 401", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/r10k/environment", bytes.NewReader(payload))
		req.Header.Add("X-GitHub-Event", "push")
		req.Header.Add("X-Hub-Signature-256", "sha256="+hex.EncodeToString([]byte("not a valid signature")))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("spoofed webhook from a provider without a secret rejected with 401", func(t *testing.T) {
		// Only a GitHub secret is configured. The provider is detected
		// from caller-controlled headers, so an unsigned request merely
		// claiming to come from Gitea must not slip past verification.
		giteaPayload, err := os.ReadFile("../lib/parsers/json/gitea/push.json")
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/r10k/environment", bytes.NewReader(giteaPayload))
		req.Header.Add("X-Gitea-Event", "push")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var body map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
			t.Fatalf("failed to unmarshal response body: %v", err)
		}
		assert.Equal(t, "Webhook verification required: no secret configured for this provider", body["message"])
	})

	t.Run("signed webhook accepted", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/r10k/environment", bytes.NewReader(payload))
		req.Header.Add("X-GitHub-Event", "push")
		req.Header.Add("X-Hub-Signature-256", signature)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusAccepted, w.Code)
	})
}
