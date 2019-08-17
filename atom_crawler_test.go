package feeder_test

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/kr/pretty"
	"github.com/p1ass/feeder"
)

func TestAtomFetch(t *testing.T) {
	// Set up mock server
	xmlFile, err := os.Open("atom_test.xml")
	if err != nil {
		t.Fatal("Failed to open test atom feed file.")
	}
	bytes, _ := ioutil.ReadAll(xmlFile)
	response := &feeder.Response{
		Path:        "/feed",
		ContentType: "application/xml",
		Body:        string(bytes),
	}
	server := feeder.NewMockServer(response)
	defer server.Close()

	updatedString := "2019-01-02T00:00:00+09:00"
	updated, _ := time.Parse(time.RFC3339, updatedString)
	publishedString := "2019-01-01T00:00:00+09:00"
	published, _ := time.Parse(time.RFC3339, publishedString)
	expected := &feeder.Items{
		[]*feeder.Item{{
			Title: "title",
			Link: &feeder.Link{
				Href: "http://example.com",
				Rel:  "alternate",
			},
			Source: nil,
			Author: &feeder.Author{
				Name:  "name",
				Email: "email@example.com",
			},
			Description: "summary_content",
			ID:          "id",
			Updated:     &updated,
			Created:     &published,
			Enclosure: &feeder.Enclosure{
				URL:    "http://example.com/image.png",
				Type:   "image/png",
				Length: "0",
			},
			Content: "content",
		}}}

	fetcher := feeder.NewAtomCrawler(server.URL + "/feed")
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
