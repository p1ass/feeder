package feeder

import (
	"log"
	"sync"
	"time"
)

// Fetcher is ...
type Fetcher interface {
	Fetch() (*Items, error)
}

type Link struct {
	Href, Rel, Type, Length string
}

type Author struct {
	Name, Email string
}

type Image struct {
	Url, Title, Link string
	Width, Height    int
}

type Enclosure struct {
	Url, Length, Type string
}

type Item struct {
	Title       string
	Link        *Link
	Source      *Link
	Author      *Author
	Description string
	Id          string
	Updated     *time.Time
	Created     *time.Time
	Enclosure   *Enclosure
	Content     string
}

type Items struct {
	items []*Item
}

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

func (items *Items) Add(i *Items) {
	items.items = append(items.items, i.items...)
}

func Crawl(fetchers ...Fetcher) *Items {
	items := &Items{}
	mutex := sync.Mutex{}
	wg := sync.WaitGroup{}

	for _, f := range fetchers {
		wg.Add(1)
		go func() {
			i, err := f.Fetch()
			if err != nil {
				log.Fatal(err)
			} else {
				mutex.Lock()
				items.Add(i)
				mutex.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()

	return items
}
