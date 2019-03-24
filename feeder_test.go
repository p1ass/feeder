package feeder

import (
	"github.com/kr/pretty"
	"reflect"
	"strconv"
	"testing"
	"time"
)

type mockFetcher struct {
	Id string
}

func (f *mockFetcher) Fetch() (*Items, error) {
	sleepTime, _ := strconv.Atoi(f.Id)
	time.Sleep(time.Second * time.Duration(sleepTime))

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
		Id:          f.Id,
		Updated:     nil,
		Created:     &published,
		Content:     "",
	}}}, nil
}

func TestCrawl(t *testing.T) {
	publishedString := "2019-01-01T00:00:00+09:00"
	published, _ := time.Parse(time.RFC3339, publishedString)

	expected := &Items{[]*Item{{
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
		Id:          "1",
		Updated:     nil,
		Created:     &published,
		Content:     "",
	}, &Item{
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
		Id:          "2",
		Updated:     nil,
		Created:     &published,
		Content:     "",
	}}}

	fetcher1 := &mockFetcher{Id: "1"}
	fetcher2 := &mockFetcher{Id: "2"}
	items := Crawl(fetcher1, fetcher2)

	if !reflect.DeepEqual(*expected, *items) {
		diffs := pretty.Diff(*expected, *items)
		t.Log(pretty.Println(diffs))
		t.Error("Crawl does not match.")

	}
}
