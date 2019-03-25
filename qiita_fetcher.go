package feeder

import (
	"encoding/json"
	"github.com/pkg/errors"
	"golang.org/x/exp/utf8string"
	"log"
	"net/http"
	"time"
)

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

// QiitaFetcher is ...
type qiitaFetcher struct {
	URL string
}

//NewQiitaFetcher is ...
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
	length := utf8string.NewString(q.Body).RuneCount()
	maxLength := 200
	if length < 200 {
		maxLength = length
	}

	i := &Item{
		Title:       q.Title,
		Link:        &Link{Href: q.URL},
		Created:     q.CreatedAt,
		Id:          q.ID,
		Description: utf8string.NewString(q.Body).Slice(0, maxLength),
	}

	if q.User != nil {
		i.Author = &Author{
			Name: q.User.ID,
		}
	}
	return i
}
