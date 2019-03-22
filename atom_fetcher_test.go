package feeder

import (
	"github.com/kr/pretty"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"
)

type Response struct {
	path, query, contentType, body string
}

func TestAtomFetch(t *testing.T) {
	// Set up mock server
	xmlFile, err := os.Open("test.atom")
	if err != nil {
		t.Fatal("Failed to open test atom feed file.")
	}
	bytes, _ := ioutil.ReadAll(xmlFile)
	response := &Response{
		path:        "/feed",
		contentType: "application/xml",
		body:        string(bytes),
	}
	handler := func(w http.ResponseWriter, r *http.Request) {
		// Send response.
		w.Header().Set("Content-Type", response.contentType)
		_, err := io.WriteString(w, response.body)
		if err != nil {
			t.Fatal(err)
		}
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	updatedString := "2019-01-02T00:00:00+09:00"
	updated, _ := time.Parse(time.RFC3339, updatedString)
	publishedString := "2019-01-01T00:00:00+09:00"
	published, _ := time.Parse(time.RFC3339, publishedString)
	expected := &Items{
		[]*Item{{
			Title: "title",
			Link: &Link{
				Href: "http://example.com",
				Rel:  "alternate",
			},
			Source: nil,
			Author: &Author{
				Name:  "name",
				Email: "email@example.com",
			},
			Description: "summary_content",
			Id:          "id",
			Updated:     updated,
			Created:     published,
			Enclosure: &Enclosure{
				Url:    "http://example.com/image.png",
				Type:   "image/png",
				Length: "0",
			},
			Content: "content",
		}}}

	fetcher := NewAtomFetcher(server.URL + "/feed")
	got, err := fetcher.Fetch()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(*expected, *got) {
		diffs := pretty.Diff(*expected, *got)
		t.Log(pretty.Println(diffs))
		t.Error("Failed to convert AtomEntry to Item.")

	}
}
