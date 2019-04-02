package feeder_test

import (
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/kr/pretty"
	"github.com/naoki-kishi/feeder"
)

type mockFetcher struct {
	Id string
}

func (f *mockFetcher) Fetch() (*feeder.Items, error) {
	sleepTime, _ := strconv.Atoi(f.Id)
	time.Sleep(time.Second * time.Duration(sleepTime))

	publishedString := "2019-01-01T00:00:00+09:00"
	published, _ := time.Parse(time.RFC3339, publishedString)
	return &feeder.Items{[]*feeder.Item{{
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
		Id:          f.Id,
		Updated:     nil,
		Created:     &published,
		Content:     "",
	}}}, nil
}

func TestCrawl(t *testing.T) {
	publishedString := "2019-01-01T00:00:00+09:00"
	published, _ := time.Parse(time.RFC3339, publishedString)

	expected := &feeder.Items{[]*feeder.Item{{
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
		Id:          "1",
		Updated:     nil,
		Created:     &published,
		Enclosure: &feeder.Enclosure{
			Url:    "http://ogp.me/logo.png",
			Type:   "image/png",
			Length: "0",
		},
		Content: "",
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
		Id:          "2",
		Updated:     nil,
		Created:     &published,
		Enclosure: &feeder.Enclosure{
			Url:    "http://ogp.me/logo.png",
			Type:   "image/png",
			Length: "0",
		}, Content: "",
	}}}

	fetcher1 := &mockFetcher{Id: "1"}
	fetcher2 := &mockFetcher{Id: "2"}
	items := feeder.Crawl(fetcher1, fetcher2)

	if !reflect.DeepEqual(*expected, *items) {
		diffs := pretty.Diff(*expected, *items)
		t.Log(pretty.Println(diffs))
		t.Error("Crawl does not match.")

	}
}

func TestItems_limitDescription(t *testing.T) {
	type fields struct {
		Items []*feeder.Item
	}
	type args struct {
		limit int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "longer than limit",
			fields: fields{Items: []*feeder.Item{
				{Description: strings.Repeat("a", 300)},
			}},
			args: args{limit: 200},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items := &feeder.Items{
				Items: tt.fields.Items,
			}
			feeder.ExportItemsLimitDescription(items, tt.args.limit)

			for _, i := range items.Items {
				got := utf8.RuneCountInString(i.Description)
				if tt.args.limit < got {
					t.Errorf("Exceed string length limit. limit=%d got=%d", tt.args.limit, got)
				}
			}
		})
	}
}
