<img src="image/feeder_logo.png" style="width:400px">

`feeder` is the RSS, Atom and JSON feed generator from multiple RSS, Atom, and any entries you want.

## Getting started

### Install
```bash
go get -u github.com/p1ass/feeder
```

### Examples
```go
import "github.com/p1ass/feeder"

func fetch(){
	rssCrawler := feeder.NewRSSCrawler("https://example.com/rss")
	qiitaCrawler := feeder.NewQiitaCrawler("https://qiita.com/api/v2/users/plus_kyoto/items")

	// Crawl data using goroutine.
	items := feeder.Crawl(rssCrawler, qiitaCrawler)

	feed := &feeder.Feed{
		Title:       "My feeds",
		Link:        &feeder.Link{Href: "https://example.com/feed"},
		Description: "My feeds.",
		Author:      &feeder.Author{
			Name: "p1ass",
			Email: "p1ass@example.com"},
		Created:     time.Now(),
		Items:       items,
	}

	json, err := feed.ToJSON() // json is a `string`
	rss, err := feed.ToRSS() // rss is a `string`
	atom, err := feed.ToAtom() // atom is a `string`
}

```

## Advanced usages

### Implement original `Crawler`
You can create a original crawler by implementing `feeder.Crawler`.
```go
type Crawler interface {
	Fetch() (*Items, error)
}
```

This is an example of Qiita API(`GET /api/v2/users/:user_id/items`).

[Qiita API v2 documentation - Qiita:Developer](https://qiita.com/api/v2/docs)
```go

type qiitaResponse struct {
	CreatedAt *time.Time `json:"created_at"`
	Title     string     `json:"title"`
	URL       string     `json:"url"`
	Body      string     `json:"body"`
	ID        string     `json:"id"`
	User      *qiitaUser `json:"user"`
}

type qiitaUser struct {
	ID string `json:"id"`
}

type qiitaCrawler struct {
	URL string
}

func (crawler *qiitaCrawler) Fetch() (*feeder.Items, error) {
	resp, err := http.Get(crawler.URL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get response from qiita.")
	}

	var qiita []*qiitaResponse
	err = json.NewDecoder(resp.Body).Decode(&qiita)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode response body.")
	}

	items := []*feeder.Item{}
	for _, i := range qiita {
		items = append(items, convertQiitaToItem(i))
	}
	return &feeder.Items{items}, nil
}

func convertQiitaToItem(q *qiitaResponse) *feeder.Item {

	i := &feeder.Item{
		Title:       q.Title,
		Link:        &feeder.Link{Href: q.URL},
		Created:     q.CreatedAt,
		Id:          q.ID,
		Description: q.Body,
	}

	if q.User != nil {
		i.Author = &feeder.Author{
			Name: q.User.ID,
		}
	}
	return i
}
```
