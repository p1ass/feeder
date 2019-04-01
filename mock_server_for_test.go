package feeder

import (
	"io"
	"net/http"
	"net/http/httptest"
)

type Response struct {
	Path, Query, ContentType, Body string
}

func NewMockServer(response *Response) *httptest.Server {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", response.ContentType)
		_, err := io.WriteString(w, response.Body)
		if err != nil {
			panic(err)
		}
	}
	return httptest.NewServer(http.HandlerFunc(handler))
}
