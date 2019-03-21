package feeder

import (
	"encoding/json"
	"fmt"
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
}

// QiitaFetcher is ...
type qiitaFetcher struct {
	UserName string
}

//NewQiitaFetcher is ...
func NewQiitaFetcher(userName string) Fetcher {
	return &qiitaFetcher{UserName: userName}
}

// Fetch is ...
func (cli *qiitaFetcher) Fetch() (*Items, error) {
	resp, err := http.Get(fmt.Sprintf("https://qiita.com/api/v2/users/%s/items", cli.UserName))
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
	i := &Item{
		Title:       q.Title,
		Link:        &Link{Href: q.URL},
		Created:     *q.CreatedAt,
		Id:          q.ID,
		Description: utf8string.NewString(q.Body).Slice(0, 200),
	}
	return i
}
