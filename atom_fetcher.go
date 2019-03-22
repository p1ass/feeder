package feeder

import (
	"encoding/xml"
	"fmt"
	"github.com/naoki-kishi/feeds"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"time"
)

type atomFetcher struct {
	URL string
}

//NewAtomFetcher is ...
func NewAtomFetcher(url string) Fetcher {
	return &atomFetcher{URL: url}
}

// Fetch is ...
func (cli *atomFetcher) Fetch() (*Items, error) {
	resp, err := http.Get(cli.URL)
	if err != nil {
		log.Fatal(err)
		return nil, errors.Wrap(err, "Failed to get response from rss.")
	}

	var atom feeds.AtomFeed
	err = xml.NewDecoder(resp.Body).Decode(&atom)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to decode response body.")
	}

	items := []*Item{}

	for _, e := range atom.Entries {
		item, err := convertAtomEntryToItem(e)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to convert RSSItem to Item.")
		}
		items = append(items, item)
	}
	return &Items{items}, nil
}

func convertAtomEntryToItem(e *feeds.AtomEntry) (*Item, error) {
	p, err := time.Parse(time.RFC3339, e.Published)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Failed to parse published time. published=%v", e.Published))
	}
	u, err := time.Parse(time.RFC3339, e.Updated)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Failed to parse updated time. updated=%v", e.Updated))
	}

	i := &Item{
		Title:       e.Title,
		Description: e.Summary.Content,
		Id:          e.Id,
		Created:     p,
		Updated:     u,
	}

	var name, email string
	if e.Author != nil {
		name, email = e.Author.Name, e.Author.Email
	}
	if len(name) > 0 || len(email) > 0 {
		i.Author = &Author{
			Name:  e.Author.Name,
			Email: e.Author.Email,
		}
	}

	if e.Content != nil {
		i.Content = e.Content.Content
	}

	for _, link := range e.Links {
		if link.Rel == "enclosure" {
			i.Enclosure = &Enclosure{
				Url:    link.Href,
				Length: link.Length,
				Type:   link.Type,
			}
		} else {
			i.Link = &Link{
				Href:   link.Href,
				Rel:    link.Rel,
				Type:   link.Type,
				Length: link.Length,
			}
		}
	}
	return i, nil
}
