package handlers

import (
	"html/template"
	"net/http"
	"net/url"

	"hawx.me/code/route"
	"hawx.me/code/trobble/data"
)

type listenCtx struct {
	Title      string
	TotalPlays int
	NowPlaying *data.Playing

	Artist      string
	AlbumArtist string
	Album       string
	Track       string
	Timestamp   int64
}

type listenHandler struct {
	db        *data.Database
	title     string
	templates *template.Template
}

func Listen(db *data.Database, title string, templates *template.Template) route.Handler {
	return &listenHandler{db, title, templates}
}

func (h listenHandler) ServeErrorHTTP(w http.ResponseWriter, r *http.Request) error {
	w.Header().Add("Content-Type", "text/html")

	vars := route.Vars(r)

	timestamp, err := url.QueryUnescape(vars["timestamp"])
	if err != nil {
		return ErrNotFound
	}

	foundTrack, ok := h.db.Scrobble(timestamp)
	if !ok {
		return ErrNotFound
	}

	ctx := listenCtx{}
	ctx.Title = h.title
	ctx.TotalPlays = h.db.TotalPlays()
	if nowPlaying, ok := h.db.GetNowPlaying(); ok {
		ctx.NowPlaying = nowPlaying
	}

	ctx.Artist = foundTrack.Artist
	ctx.AlbumArtist = foundTrack.AlbumArtist
	ctx.Album = foundTrack.Album
	ctx.Track = foundTrack.Track
	ctx.Timestamp = foundTrack.Timestamp

	return h.templates.ExecuteTemplate(w, "listen.gotmpl", ctx)
}
