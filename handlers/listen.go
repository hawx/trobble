package handlers

import (
	"net/http"
	"net/url"

	"hawx.me/code/route"
	"hawx.me/code/trobble/data"
	"hawx.me/code/trobble/views"
)

type listenCtx struct {
	Title      string
	TotalPlays int
	NowPlaying *data.Playing

	Artist    string
	Album     string
	Track     string
	Timestamp int64
}

type listenHandler struct {
	db    *data.Database
	title string
}

func Listen(db *data.Database, title string) http.Handler {
	return &listenHandler{db, title}
}

func (h listenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")

	vars := route.Vars(r)

	timestamp, err := url.QueryUnescape(vars["timestamp"])
	if err != nil {
		http.NotFound(w, r)
		return
	}

	foundTrack, ok := h.db.Scrobble(timestamp)
	if !ok {
		http.NotFound(w, r)
		return
	}

	ctx := listenCtx{}
	ctx.Title = h.title
	ctx.TotalPlays = h.db.TotalPlays()
	if nowPlaying, ok := h.db.GetNowPlaying(); ok {
		ctx.NowPlaying = nowPlaying
	}

	ctx.Artist = foundTrack.Artist
	ctx.Album = foundTrack.Album
	ctx.Track = foundTrack.Track
	ctx.Timestamp = foundTrack.Timestamp

	views.Listen.Execute(w, ctx)
}
