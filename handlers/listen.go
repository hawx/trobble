package handlers

import (
	"html/template"
	"log"
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

func Listen(db *data.Database, title string, templates *template.Template) http.Handler {
	return &listenHandler{db, title, templates}
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
	ctx.AlbumArtist = foundTrack.AlbumArtist
	ctx.Album = foundTrack.Album
	ctx.Track = foundTrack.Track
	ctx.Timestamp = foundTrack.Timestamp

	if err := h.templates.ExecuteTemplate(w, "listen.gotmpl", ctx); err != nil {
		log.Println(err)
	}
}
