package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/feeds"
	"hawx.me/code/trobble/data"
)

func Feed(db *data.Database, title, url string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		feed := &feeds.Feed{
			Title:   title,
			Link:    &feeds.Link{Href: url},
			Created: time.Now(),
		}

		recent, err := db.RecentlyPlayed()
		if err != nil {
			log.Println(err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		for _, played := range recent {
			feed.Items = append(feed.Items, &feeds.Item{
				Title:       played.Track,
				Description: played.Artist,
				Link:        &feeds.Link{Href: url + "listen/" + strconv.FormatInt(played.Timestamp, 10)},
				Created:     time.Unix(played.Timestamp, 0),
			})
		}

		w.Header().Add("Content-Type", "application/rss+xml")

		if err := feed.WriteRss(w); err != nil {
			log.Println(err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	})
}
