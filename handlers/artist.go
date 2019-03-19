package handlers

import (
	"html/template"
	"net/http"
	"net/url"

	"hawx.me/code/route"
	"hawx.me/code/trobble/data"
)

const playsLength = 48

type artistCtx struct {
	Title      string
	TotalPlays int
	NowPlaying *data.Playing

	Name     string
	Plays    []int
	MaxPlays int
	Tracks   []data.Track
	Albums   []data.Album
}

type artistHandler struct {
	db        *data.Database
	title     string
	templates *template.Template
}

func Artist(db *data.Database, title string, templates *template.Template) http.Handler {
	return &artistHandler{db, title, templates}
}

func (h artistHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")

	name, err := url.QueryUnescape(route.Vars(r)["artist"])
	if err != nil {
		http.NotFound(w, r)
		return
	}

	ctx := artistCtx{}
	ctx.Title = h.title
	ctx.TotalPlays = h.db.TotalPlays()
	if nowPlaying, ok := h.db.GetNowPlaying(); ok {
		ctx.NowPlaying = nowPlaying
	}

	ctx.Name = name
	ctx.Tracks = h.db.ArtistTopTracks(name, 50)
	ctx.Albums = h.db.ArtistTopAlbums(name, 4)

	if len(ctx.Tracks) == 0 {
		http.NotFound(w, r)
		return
	}

	calcSegment := func(play data.PlayCount) int {
		return play.Year*12 + int(play.Month)
	}

	plays := h.db.ArtistPlays(name)
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

	h.templates.ExecuteTemplate(w, "artist.gotmpl", ctx)
}
