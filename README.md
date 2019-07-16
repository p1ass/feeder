<img src="image/feeder_logo.png" style="width:400px">

feeder is a RSS or JSON feeds generator from multiple RSS, Atom, Qiita, and so on

## Getting started

### Install
```bash
go get -u github.com/naoki-kishi/feeder
```

### Examples
```go
import "github.com/naoki-kishi/feeder"

func fetch(){
	rssFetcher := feeder.NewRSSFetcher("https://example.com/rss")
	qiitaFetcher := feeder.NewQiitaFetcher("https://qiita.com/api/v2/users/plus_kyoto/items")

	// Fetch data using goroutine.
	items := feeder.Crawl(rssFetcher, qiitaFetcher)

	feed := &feeder.Feed{
		Title:       "My feeds",
		Link:        &feeder.Link{Href: "https://example.com/feed"},
		Description: "My feeds.",
		Author:      &feeder.Author{
			Name: "naoki-kishi",
			Email: "naoki-kishi@example.com"},
		Created:     time.Now(),
		Items:       items,
	}

	json, err := feed.ToJSON() // json is a `string`
	rss, err := feed.ToRSS() // rss is a `string`
	atom, err := feed.ToAtom() // atom is a `string`
}

```

## Advanced usages

### Implement original `Fetcher`
You can create a original fetcher by implementing `feeder.Fetcher`.
```go
type Fetcher interface {
	Fetch() (*Items, error)
}
```

This is a example of Qiita API(`GET /api/v2/users/:user_id/items`).

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

type qiitaFetcher struct {
	URL string
}

func (fetcher *qiitaFetcher) Fetch() (*feeder.Items, error) {
	resp, err := http.Get(fetcher.URL)
	if err != nil {
		log.Fatal(err)
		return nil, errors.Wrap(err, "Failed to get response from qiita.")
	}

	var qiita []*qiitaResponse
	err = json.NewDecoder(resp.Body).Decode(&qiita)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to decode response body.")
	}

	items := []*feeder.Item{}
	for _, i := range qiita {
		items = append(items, convertQiitaToItem(i))
	}
	return &feeder.Items{items}, nil
}

func convertQiitaToItem(q *qiitaResponse) *feeder.Item {
	length := utf8string.NewString(q.Body).RuneCount()
	maxLength := 200
	if length < 200 {
		maxLength = length
	}

	i := &feeder.Item{
		Title:       q.Title,
		Link:        &feeder.Link{Href: q.URL},
		Created:     q.CreatedAt,
		Id:          q.ID,
		Description: utf8string.NewString(q.Body).Slice(0, maxLength),
	}

	if q.User != nil {
		i.Author = &feeder.Author{
			Name: q.User.ID,
		}
	}
	return i
}
```
