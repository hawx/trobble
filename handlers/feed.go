package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/gorilla/feeds"
	"hawx.me/code/mux"
	"hawx.me/code/trobble/data"
)

func Feed(db *data.Database, title, url, name, homepage string) http.Handler {
	return mux.Accept{
		"application/activity+json": ActivityStream(db, title, url, name, homepage),
		"*/*": Rss(db, title, url, name),
	}
}

func ActivityStream(db *data.Database, title, rooturl, name, homepage string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/activity+json")

		stream := OrderedCollection{
			Context:      "https://www.w3.org/ns/activitystreams",
			Name:         title,
			Type:         "OrderedCollection",
			Id:           rooturl + "feed",
			OrderedItems: []Activity{},
		}

		for _, played := range db.RecentlyPlayed() {
			href := rooturl + path.Join("artist", url.QueryEscape(played.Artist), url.QueryEscape(played.Album), url.QueryEscape(played.Track))

			stream.OrderedItems = append(stream.OrderedItems, Activity{
				Type:      "Listen",
				Name:      played.Track,
				Summary:   "Listened to " + played.Track + " by " + played.Artist,
				Published: time.Unix(played.Timestamp, 0),
				Actor: Actor{
					Type: "Person",
					Name: name,
					Id:   homepage,
				},
				Object: Object{
					Type: "Link",
					Name: played.Track + " by " + played.Artist,
					Href: href,
					Id:   href,
				},
			})
		}

		if err := json.NewEncoder(w).Encode(stream); err != nil {
			log.Println(err)
			w.WriteHeader(500)
		}
	})
}

type Object struct {
	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
	Href string `json:"href,omitempty"`
	Id   string `json:"id,omitempty"`
}

type OrderedCollection struct {
	Context      string     `json:"@context,omitempty"`
	Type         string     `json:"type,omitempty"`
	Name         string     `json:"name,omitempty"`
	Id           string     `json:"id,omitempty"`
	OrderedItems []Activity `json:"orderedItems,omitempty"`
}

type Activity struct {
	Type      string    `json:"type,omitempty"`
	Name      string    `json:"name,omitempty"`
	Summary   string    `json:"summary,omitempty"`
	Published time.Time `json:"published,omitempty"`
	Actor     Actor     `json:"actor,omitempty"`
	Object    Object    `json:"object,omitempty"`
}

type Actor struct {
	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
	Id   string `json:"id,omitempty"`
}

func Rss(db *data.Database, title, url, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		feed := &feeds.Feed{
			Title:   title,
			Link:    &feeds.Link{Href: url},
			Created: time.Now(),
			Author:  &feeds.Author{Name: name},
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
