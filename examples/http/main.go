package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/p1ass/feeder"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		blogCrawler := feeder.NewRSSCrawler("https://blog.p1ass.com/index.xml")
		items, err := feeder.Crawl(blogCrawler)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		feed := feeder.Feed{
			Title:       "Example feed title",
			Description: "Example feed description",
			Author: &feeder.Author{
				Name:  "p1ass",
				Email: "concat@p1ass.com",
			},
			Updated: time.Now(),
			Created: time.Now(),
			Items:   items,
		}

		reader, err := feed.ToJSONReader()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if _, err := io.Copy(w, reader); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		r.Header.Set("Content-Type", "application/json; charset=utf-8")
		return
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
