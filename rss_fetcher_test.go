package feeder_test

import (
	"github.com/kr/pretty"
	"github.com/naoki-kishi/feeder"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestRSSFetch(t *testing.T) {
	// Set up mock server
	xmlFile, err := os.Open("rss_test.xml")
	if err != nil {
		t.Fatal("Failed to open test rss feed file.")
	}
	bytes, _ := ioutil.ReadAll(xmlFile)
	response := &feeder.Response{
		Path:        "/rss",
		ContentType: "application/xml",
		Body:        string(bytes),
	}
	server := feeder.NewMockServer(response)
	defer server.Close()

	publishedString := "2019-01-01T00:00:00+09:00"
	published, _ := time.Parse(time.RFC3339, publishedString)
	expected := &feeder.Items{
		[]*feeder.Item{{
			Title: "title",
			Link: &feeder.Link{
				Href: "http://example.com",
				Rel:  "",
			},
			Source: nil,
			Author: &feeder.Author{
				Name: "name",
			},
			Description: "summary_content",
			Id:          "id",
			Updated:     nil,
			Created:     &published,
			Enclosure: &feeder.Enclosure{
				Url:    "http://example.com/image.png",
				Type:   "image/png",
				Length: "0",
			},
			Content: "",
		}}}

	fetcher := feeder.NewRSSFetcher(server.URL + "/rss")
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
