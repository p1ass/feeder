package feeder

import (
	ogp "github.com/otiai10/opengraph"
	"golang.org/x/exp/utf8string"
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

	items.limitDescription(200)

	fetchOGP(items)

	return items
}

func (items *Items) limitDescription(limit int) {
	for _, i := range items.Items {
		length := utf8string.NewString(i.Description).RuneCount()
		maxLength := limit
		if length < limit {
			maxLength = length
		}
		i.Description = utf8string.NewString(i.Description).Slice(0, maxLength)
	}
}

func fetchOGP(items *Items) *Items {
	wg := sync.WaitGroup{}

	for _, i := range items.Items {
		wg.Add(1)
		i := i
		go func() {
			if i.Enclosure == nil || i.Enclosure.Url == "" {
				og, err := ogp.Fetch(i.Link.Href)
				if err != nil {
					log.Printf("Failed to fetch ogp. %#v", err)
				}

				if len(og.Image) > 0 {
					image := og.Image[0]
					i.Enclosure = &Enclosure{}
					i.Enclosure.Url = image.URL

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
