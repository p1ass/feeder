package feeder

import (
	"log"
	"sync"
	"time"

	ogp "github.com/otiai10/opengraph"
)

// Deprecated: Fetcher is replaced by Crawler
type Fetcher interface {
	Fetch() (*Items, error)
}

// Crawler is interface for crawling
type Crawler interface {
	Fetch() (*Items, error)
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

// Items represents slice of item
type Items struct {
	Items []*Item
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
	Items       Items
	Copyright   string
	Image       *Image
}

// Add adds item to Items
func (items *Items) Add(i *Items) {
	items.Items = append(items.Items, i.Items...)
}

// Crawl is function that crawls all site using goroutine.
// func Crawl(fetchers ...Fetcher) *Items is deprecated
func Crawl(crawlers ...Crawler) *Items {
	items := &Items{}
	mutex := sync.Mutex{}
	wg := sync.WaitGroup{}

	for _, f := range crawlers {
		wg.Add(1)
		go func(f Crawler) {
			i, err := f.Fetch()
			if err != nil {
				log.Fatal(err)
			} else {
				mutex.Lock()
				items.Add(i)
				mutex.Unlock()
			}
			wg.Done()
		}(f)
	}
	wg.Wait()

	fetchOGP(items)

	return items
}

func fetchOGP(items *Items) *Items {
	wg := sync.WaitGroup{}

	for _, i := range items.Items {
		wg.Add(1)
		i := i
		go func() {
			if i.Enclosure == nil || i.Enclosure.URL == "" {
				og, err := ogp.Fetch(i.Link.Href)
				if err != nil {
					log.Printf("Failed to fetch ogp. %#v", err)
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
			wg.Done()
		}()
	}
	wg.Wait()

	return items
}
