package handlers

import (
	"html/template"
	"log"
	"net/http"
	"net/url"

	"hawx.me/code/route"
	"hawx.me/code/trobble/data"
)

type trackCtx struct {
	Title      string
	TotalPlays int
	NowPlaying *data.Playing

	Artist     string
	Album      string
	Track      string
	Plays      []int
	MaxPlays   int
	TrackPlays int
}

type trackHandler struct {
	db        *data.Database
	title     string
	templates *template.Template
}

func Track(db *data.Database, title string, templates *template.Template) http.Handler {
	return &trackHandler{db, title, templates}
}

func (h trackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")

	vars := route.Vars(r)

	artist, err := url.QueryUnescape(vars["artist"])
	if err != nil {
		http.NotFound(w, r)
		return
	}

	album, err := url.QueryUnescape(vars["album"])
	if err != nil {
		http.NotFound(w, r)
		return
	}

	track, err := url.QueryUnescape(vars["track"])
	if err != nil {
		http.NotFound(w, r)
		return
	}

	foundTrack, err := h.db.GetTrack(artist, album, track)
	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if foundTrack == nil {
		http.NotFound(w, r)
		return
	}

	ctx := trackCtx{}
	ctx.Title = h.title
	ctx.TotalPlays = h.db.TotalPlays()
	if nowPlaying, ok := h.db.GetNowPlaying(); ok {
		ctx.NowPlaying = nowPlaying
	}

	ctx.Artist = artist
	ctx.Album = album

	ctx.Track = foundTrack.Track
	ctx.TrackPlays = foundTrack.Count

	calcSegment := func(play data.PlayCount) int {
		return play.Year*12 + int(play.Month)
	}

	plays := h.db.TrackPlays(artist, album, track)
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

	h.templates.ExecuteTemplate(w, "track.gotmpl", ctx)
}
