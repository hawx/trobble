package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/feeds"
	"github.com/hawx/trobble/data"
)

func Feed(db *data.Database, title, url string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		feed := &feeds.Feed{
			Title:   title,
			Link:    &feeds.Link{Href: url},
			Created: time.Now(),
		}

		for _, played := range db.RecentlyPlayed() {
			feed.Items = append(feed.Items, &feeds.Item{
				Title:       played.Track,
				Description: played.Artist,
				Link:        &feeds.Link{Href: url},
				Created:     time.Unix(played.Timestamp, 0),
			})
		}

		w.Header().Add("Content-Type", "application/rss+xml")

		err := feed.WriteRss(w)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}
	})
}
