package feeder

import (
	"github.com/kr/pretty"
	"reflect"
	"testing"
	"time"
)

type mockFetcher struct {
}

func (f *mockFetcher) Fetch() (*Items, error) {
	publishedString := "2019-01-01T00:00:00+09:00"
	published, _ := time.Parse(time.RFC3339, publishedString)
	return &Items{[]*Item{{
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
	}}}, nil
}

func TestCrawl(t *testing.T) {
	publishedString := "2019-01-01T00:00:00+09:00"
	published, _ := time.Parse(time.RFC3339, publishedString)

	item := &Item{
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
	}
	expected := &Items{[]*Item{item, item}}

	fetcher1 := &mockFetcher{}
	fetcher2 := &mockFetcher{}
	items := Crawl(fetcher1, fetcher2)

	if !reflect.DeepEqual(*expected, *items) {
		diffs := pretty.Diff(*expected, *items)
		t.Log(pretty.Println(diffs))
		t.Error("Crawl does not match.")

	}
}
