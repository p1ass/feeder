package feeder

import (
	"io"
	"net/http"
	"net/http/httptest"
)

type response struct {
	path, query, contentType, body string
}

func newMockServer(response *response) *httptest.Server {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", response.contentType)
		_, err := io.WriteString(w, response.body)
		if err != nil {
			panic(err)
		}
	}
	return httptest.NewServer(http.HandlerFunc(handler))
}
