package handlers

import (
	"html/template"
	"log"
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

	topTracks, err := h.db.ArtistTopTracks(name, 50)
	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	ctx.Tracks = topTracks

	albums, err := h.db.ArtistTopAlbums(name, 4)
	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	ctx.Albums = albums

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

	if err := h.templates.ExecuteTemplate(w, "artist.gotmpl", ctx); err != nil {
		log.Println(err)
	}
}
