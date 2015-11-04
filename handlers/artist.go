package handlers

import (
	"net/http"
	"net/url"

	"hawx.me/code/route"
	"hawx.me/code/trobble/data"
	"hawx.me/code/trobble/views"
)

type artistCtx struct {
	Title      string
	TotalPlays int
	NowPlaying *data.Track

	Name     string
	Plays    []int
	MaxPlays int
	Tracks   []data.Track
}

type artistHandler struct {
	db    *data.Database
	title string
}

func Artist(db *data.Database, title string) http.Handler {
	return &artistHandler{db, title}
}

func (h artistHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")

	name, err := url.QueryUnescape(route.Vars(r)["name"])
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

	if len(ctx.Tracks) == 0 {
		http.NotFound(w, r)
		return
	}

	plays := h.db.ArtistPlays(name)
	firstPlay, lastPlay := plays[0], plays[len(plays)-1]

	calcSegment := func(play data.PlayCount) int {
		return play.Year*12 + int(play.Month)
	}

	startPlay := calcSegment(firstPlay)
	segments := calcSegment(lastPlay) - startPlay + 1
	ctx.Plays = make([]int, segments)

	for _, play := range plays {
		ctx.Plays[calcSegment(play)-startPlay] = play.Count
		if play.Count > ctx.MaxPlays {
			ctx.MaxPlays = play.Count
		}
	}

	views.Artist.Execute(w, ctx)
}
