package feeder

import (
	"encoding/xml"
	"log"
	"net/http"
	"time"

	"github.com/p1ass/feeds"
	"github.com/pkg/errors"
)

type rssCrawler struct {
	URL string
}

// NewRSSCrawler returns rSSCrawler
func NewRSSCrawler(url string) Crawler {
	return &rssCrawler{URL: url}
}

// Deprecated: Use NewAtomCrawler instead of NewRSSFetcher
func NewRSSFetcher(url string) Fetcher {
	return &rssCrawler{URL: url}
}

// Fetch fetches entries from rss feed
func (fetcher *rssCrawler) Fetch() (*Items, error) {
	resp, err := http.Get(fetcher.URL)
	if err != nil {
		log.Println(err)
		return nil, errors.Wrap(err, "Failed to get response from rss.")
	}
	defer resp.Body.Close()

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
		ID:          i.Guid,
		Created:     &t,
	}

	if i.Author != "" {
		item.Author = &Author{Name: i.Author}
	}

	if i.Content != nil {
		item.Content = i.Content.Content
	}

	if i.Enclosure != nil {
		item.Enclosure = &Enclosure{
			URL:    i.Enclosure.Url,
			Length: i.Enclosure.Length,
			Type:   i.Enclosure.Type}
	}

	return item, nil
}
