package feeder

import (
	"github.com/kr/pretty"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestQiitaFetch(t *testing.T) {
	// Set up mock server
	jsonFile, err := os.Open("qiita_test.json")
	if err != nil {
		t.Fatal("Failed to open test rss feed file.")
	}
	bytes, _ := ioutil.ReadAll(jsonFile)
	response := &response{
		path:        "/qiita",
		contentType: "application/json",
		body:        string(bytes),
	}
	server := newMockServer(response)
	defer server.Close()

	publishedString := "2019-01-01T00:00:00+09:00"
	published, _ := time.Parse(time.RFC3339, publishedString)
	expected := &Items{
		[]*Item{{
			Title: "title",
			Link: &Link{
				Href: "http://example.com",
				Rel:  "",
			},
			Source: nil,
			Author: &Author{
				Name: "name",
			},
			Description: "summary_content",
			Id:          "id",
			Updated:     nil,
			Created:     &published,
			Content:     "",
		}}}

	fetcher := NewQiitaFetcher(server.URL + "/qiita")
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
