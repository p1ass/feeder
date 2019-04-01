package feeder

import (
	"github.com/frozzare/go-ogp"
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

	Id        string
	Updated   *time.Time
	Created   *time.Time
	Enclosure *Enclosure
	Content   string
}

type Items struct {
	Items []*Item
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
	items.Items = append(items.Items, i.Items...)
}

func Crawl(fetchers ...Fetcher) *Items {
	items := &Items{}
	mutex := sync.Mutex{}
	wg := sync.WaitGroup{}

	for _, f := range fetchers {
		wg.Add(1)
		go func(f Fetcher) {
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
			if i.Enclosure == nil || i.Enclosure.Url == "" {
				list := ogp.Fetch(i.Link.Href)

				if ogpLink, ok := list["image"]; ok {
					i.Enclosure = &Enclosure{}
					i.Enclosure.Url = ogpLink.(string)

					if imageType, ok := list["image:type"]; ok {
						i.Enclosure.Type = imageType.(string)
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
