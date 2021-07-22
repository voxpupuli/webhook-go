package parsers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

type Header struct {
	Name  string
	Value string
}

func getGinContext(filename string, headers []Header) (*gin.Context, *httptest.ResponseRecorder, error) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	rawjson, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}

	req := &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header),
		Body:   rawjson,
		Method: "POST",
	}

	for _, header := range headers {
		req.Header.Add(header.Name, header.Value)
	}

	c.Request = req

	return c, w, nil
}
