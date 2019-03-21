package feeder

import (
	"encoding/xml"
	"github.com/naoki-kishi/feeds"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"time"
)

type rssFetcher struct {
	URL string
}

//NewRSSFetcher is ...
func NewRSSFetcher(url string) Fetcher {
	return &rssFetcher{URL: url}
}

func (cli *rssFetcher) Fetch() (*Items, error) {
	resp, err := http.Get(cli.URL)
	if err != nil {
		log.Fatal(err)
		return nil, errors.Wrap(err, "Failed to get response from rss.")
	}

	var rss feeds.RssFeedXml
	err = xml.NewDecoder(resp.Body).Decode(&rss)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to decode response body.")
	}

	items := []*Item{}

	for _, i := range rss.Channel.Items {
		item, err := convertRssItemToItem(i)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to convert RSSItem to Item.")
		}
		items = append(items, item)
	}
	return &Items{items}, nil
}

func convertRssItemToItem(i *feeds.RssItem) (*Item, error) {
	t, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", i.PubDate)
	if err != nil {
		return nil, errors.Wrap(err, "Parse Error")
	}

	item := &Item{
		Title:       i.Title,
		Link:        &Link{Href: i.Link},
		Description: i.Description,
		Id:          i.Guid,
		Created:     t,
		Source:      &Link{Href: i.Source},
	}

	if i.Author != "" {
		item.Author = &Author{Name: i.Author}
	}

	if i.Content != nil {
		item.Content = i.Content.Content
	}

	if i.Enclosure != nil {
		item.Enclosure = &Enclosure{
			Url:    i.Enclosure.Url,
			Length: i.Enclosure.Length,
			Type:   i.Enclosure.Type}
	}

	return item, nil
}
