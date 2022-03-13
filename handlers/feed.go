package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/feeds"
	"hawx.me/code/route"
	"hawx.me/code/trobble/data"
)

func Feed(db *data.Database, title, url string) route.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		feed := &feeds.Feed{
			Title:   title,
			Link:    &feeds.Link{Href: url},
			Created: time.Now(),
		}

		recent, err := db.RecentlyPlayed()
		if err != nil {
			return err
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

		return feed.WriteRss(w)
	}
}
