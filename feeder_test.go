package feeder_test

import (
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/kr/pretty"
	"github.com/p1ass/feeder"
)

type mockFetcher struct {
	Id string
}

func (f *mockFetcher) Crawl() ([]*feeder.Item, error) {
	sleepTime, _ := strconv.Atoi(f.Id)
	time.Sleep(time.Second * time.Duration(sleepTime))

	publishedString := "2019-01-01T00:00:00+09:00"
	published, _ := time.Parse(time.RFC3339, publishedString)
	return []*feeder.Item{{
		Title: "title",
		Link: &feeder.Link{
			Href: "http://ogp.me",
			Rel:  "",
		},
		Source: nil,
		Author: &feeder.Author{
			Name: "name",
		},
		Description: "summary_content",
		ID:          f.Id,
		Updated:     nil,
		Created:     &published,
		Content:     "",
	}}, nil
}

func TestCrawl(t *testing.T) {
	publishedString := "2019-01-01T00:00:00+09:00"
	published, _ := time.Parse(time.RFC3339, publishedString)

	expected := []*feeder.Item{{
		Title: "title",
		Link: &feeder.Link{
			Href: "http://ogp.me",
			Rel:  "",
		},
		Source: nil,
		Author: &feeder.Author{
			Name: "name",
		},
		Description: "summary_content",
		ID:          "1",
		Updated:     nil,
		Created:     &published,
		Content:     "",
	}, {
		Title: "title",
		Link: &feeder.Link{
			Href: "http://ogp.me",
			Rel:  "",
		},
		Source: nil,
		Author: &feeder.Author{
			Name: "name",
		},
		Description: "summary_content",
		ID:          "2",
		Updated:     nil,
		Created:     &published,
		Content:     "",
	}}

	crawler1 := &mockFetcher{Id: "1"}
	crawler2 := &mockFetcher{Id: "2"}
	items, err := feeder.Crawl(crawler1, crawler2)
	if err != nil {
		t.Errorf("Crawl() shoud return nil error error = %v", err)
		return
	}

	if !reflect.DeepEqual(&expected, &items) {
		diffs := pretty.Diff(expected, items)
		t.Log(pretty.Println(diffs))
		t.Error("Crawl does not match.")
	}
}
