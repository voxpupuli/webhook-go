package rcserver

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
)

var mockRocketChat *MockRocketChat

type Attachment struct {
	Color string `json:"fallback"`
	Text  string `json:"text"`
}

type MockRocketChat struct {
	Server   *httptest.Server
	Received struct {
		Attachment []Attachment
	}
}

func New() *MockRocketChat {
	mockRocketChat = &MockRocketChat{Server: mockServer()}
	return mockRocketChat
}

func mockServer() *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc("api/v1/chat.postMessage", handlePostMessage)

	return httptest.NewServer(handler)
}

func parseAttachment(data string) []Attachment {
	a := make([]Attachment, 0)
	json.Unmarshal([]byte(data), &a)

	return a
}

func handlePostMessage(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	kvs := strings.Split(string(body), "&")

	m := make(map[string]string)

	for _, s := range kvs {
		kv := strings.Split(s, "=")
		s, err := url.QueryUnescape(kv[1])
		if err != nil {
			m[kv[0]] = kv[1]
		} else {
			m[kv[0]] = s
		}
	}

	mockRocketChat.Received.Attachment = parseAttachment(m["attachments"])

	const response = `{
	  "ts": 0000,
	  "channel": "%s",
	  "message": {
		"alias": "",
		"msg": "%s",
		"parseUrls": true,
		"groupable": false,
		"ts": "2016-12-14T20:56:05.117Z",
		"u": {
		  "_id": "y65tAmHs93aDChMWu",
		  "username": "graywolf336"
		},
		"rid": "GENERAL",
		"_updatedAt": "2016-12-14T20:56:05.119Z",
		"_id": "jC9chsFddTvsbFQG7"
	  },
	  "success": true
	}`

	s := fmt.Sprintf(response, m["channel"], m["msg"])
	_, _ = w.Write([]byte(s))
}
