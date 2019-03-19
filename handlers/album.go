package handlers

import (
	"html/template"
	"net/http"
	"net/url"

	"hawx.me/code/route"
	"hawx.me/code/trobble/data"
)

type albumCtx struct {
	Title      string
	TotalPlays int
	NowPlaying *data.Playing

	Artist   string
	Album    string
	Plays    []int
	MaxPlays int
	Tracks   []data.Track
}

type albumHandler struct {
	db        *data.Database
	title     string
	templates *template.Template
}

func Album(db *data.Database, title string, templates *template.Template) http.Handler {
	return &albumHandler{db, title, templates}
}

func (h albumHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")

	artist, err := url.QueryUnescape(route.Vars(r)["artist"])
	if err != nil {
		http.NotFound(w, r)
		return
	}

	album, err := url.QueryUnescape(route.Vars(r)["album"])
	if err != nil {
		http.NotFound(w, r)
		return
	}

	ctx := albumCtx{}
	ctx.Title = h.title
	ctx.TotalPlays = h.db.TotalPlays()
	if nowPlaying, ok := h.db.GetNowPlaying(); ok {
		ctx.NowPlaying = nowPlaying
	}

	ctx.Artist = artist
	ctx.Album = album
	ctx.Tracks = h.db.AlbumTracks(artist, album)

	if len(ctx.Tracks) == 0 {
		http.NotFound(w, r)
		return
	}

	calcSegment := func(play data.PlayCount) int {
		return play.Year*12 + int(play.Month)
	}

	plays := h.db.AlbumPlays(artist, album)
	lastPlay := plays[len(plays)-1]
	lastSegment := calcSegment(lastPlay)

	ctx.Plays = make([]int, playsLength)

	for i := len(plays) - 1; i >= 0; i-- {
		play := plays[i]

		j := calcSegment(play) - 1 - lastSegment + playsLength
		if j < 0 {
			break
		}

		ctx.Plays[j] = play.Count
		if play.Count > ctx.MaxPlays {
			ctx.MaxPlays = play.Count
		}
	}

	h.templates.ExecuteTemplate(w, "album.gotmpl", ctx)
}
