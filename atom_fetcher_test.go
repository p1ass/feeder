package feeder

import (
	"github.com/kr/pretty"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestAtomFetch(t *testing.T) {
	// Set up mock server
	xmlFile, err := os.Open("test.atom")
	if err != nil {
		t.Fatal("Failed to open test atom feed file.")
	}
	bytes, _ := ioutil.ReadAll(xmlFile)
	response := &response{
		path:        "/feed",
		contentType: "application/xml",
		body:        string(bytes),
	}
	server := newMockServer(response)
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
