package feeder

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type qiitaResponse struct {
	CreatedAt    *time.Time `json:"created_at"`
	Title        string     `json:"title"`
	URL          string     `json:"url"`
	RenderedBody string     `json:"rendered_body"`
	ID           string     `json:"id"`
	User         *qiitaUser `json:"user"`
}

type qiitaUser struct {
	ID string `json:"id"`
}

// QiitaFetcher is ...
type qiitaFetcher struct {
	URL string
}

func NewQiitaCrawler(url string) Crawler {
	return &qiitaFetcher{URL: url}
}

// Deprecated: Use NewAtomCrawler instead of NewQiitaFetcher
func NewQiitaFetcher(url string) Fetcher {
	return &qiitaFetcher{URL: url}
}

// Fetch is ...
func (fetcher *qiitaFetcher) Fetch() (*Items, error) {
	resp, err := http.Get(fetcher.URL)
	if err != nil {
		log.Fatal(err)
		return nil, errors.Wrap(err, "Failed to get response from qiita.")
	}
	defer resp.Body.Close()

	var qiita []*qiitaResponse
	err = json.NewDecoder(resp.Body).Decode(&qiita)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to decode response body.")
	}

	items := []*Item{}

	for _, i := range qiita {
		items = append(items, convertQiitaToItem(i))
	}
	return &Items{items}, nil
}

func convertQiitaToItem(q *qiitaResponse) *Item {
	i := &Item{
		Title:       q.Title,
		Link:        &Link{Href: q.URL},
		Created:     q.CreatedAt,
		Id:          q.ID,
		Description: q.RenderedBody,
	}

	if q.User != nil {
		i.Author = &Author{
			Name: q.User.ID,
		}
	}
	return i
}
