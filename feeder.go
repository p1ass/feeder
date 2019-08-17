package feeder

import (
	"sync"
	"time"

	"github.com/pkg/errors"

	ogp "github.com/otiai10/opengraph"
	"golang.org/x/sync/errgroup"
)

// Crawler is interface for crawling
type Crawler interface {
	Crawl() ([]*Item, error)
}

// Link represents http link
type Link struct {
	Href, Rel, Type, Length string
}

// Author represents entry author
type Author struct {
	Name, Email string
}

// Image represents image
type Image struct {
	URL, Title, Link string
	Width, Height    int
}

// Enclosure represents og link
type Enclosure struct {
	URL, Length, Type string
}

// Item represents a entry
type Item struct {
	Title       string
	Link        *Link
	Source      *Link
	Author      *Author
	Description string

	ID        string
	Updated   *time.Time
	Created   *time.Time
	Enclosure *Enclosure
	Content   string
}

// Feed represents rss feed or atom feed
type Feed struct {
	Title       string
	Link        *Link
	Description string
	Author      *Author
	Updated     time.Time
	Created     time.Time
	Id          string
	Subtitle    string
	Items       []*Item
	Copyright   string
	Image       *Image
}

// Crawl is function that crawls all site using goroutine.
// func Crawl(crawlers ...Fetcher) *Items is deprecated
func Crawl(crawlers ...Crawler) ([]*Item, error) {
	items := []*Item{}
	mutex := sync.Mutex{}

	eg := errgroup.Group{}
	for _, f := range crawlers {
		f := f
		eg.Go(func() error {
			i, err := f.Crawl()
			if err != nil {
				return err
			} else {
				mutex.Lock()
				items = append(items, i...)
				mutex.Unlock()
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, errors.Wrap(err, "failed to crawl items")
	}

	items, err := fetchOGP(items)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch ogp")
	}

	return items, nil
}

func fetchOGP(items []*Item) ([]*Item, error) {
	eg := errgroup.Group{}

	for _, i := range items {
		i := i
		eg.Go(func() error {
			if i.Enclosure == nil || i.Enclosure.URL == "" {
				og, err := ogp.Fetch(i.Link.Href)
				if err != nil {
					return err
				}

				if len(og.Image) > 0 {
					image := og.Image[0]
					i.Enclosure = &Enclosure{}
					i.Enclosure.URL = image.URL

					if image.Type != "" {
						i.Enclosure.Type = image.Type
					} else {
						i.Enclosure.Type = "image/png"
					}
					i.Enclosure.Length = "0"
				}
			}
			return nil
		})
	}

	return items, nil
}
