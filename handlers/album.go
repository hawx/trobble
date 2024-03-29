package handlers

import (
	"html/template"
	"net/http"
	"net/url"

	"hawx.me/code/route"
	"hawx.me/code/trobble/data"
)

type albumCtx struct {
	Title       string
	TotalPlays  int
	NowPlaying  *data.Playing
	ShowArtists bool

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

func Album(db *data.Database, title string, templates *template.Template) route.Handler {
	return &albumHandler{db, title, templates}
}

func (h albumHandler) ServeErrorHTTP(w http.ResponseWriter, r *http.Request) error {
	w.Header().Add("Content-Type", "text/html")

	albumArtist, err := url.QueryUnescape(route.Vars(r)["albumArtist"])
	if err != nil {
		return ErrNotFound
	}

	album, err := url.QueryUnescape(route.Vars(r)["album"])
	if err != nil {
		return ErrNotFound
	}

	ctx := albumCtx{}
	ctx.Title = h.title
	ctx.TotalPlays = h.db.TotalPlays()
	if nowPlaying, ok := h.db.GetNowPlaying(); ok {
		ctx.NowPlaying = nowPlaying
	}

	ctx.Artist = albumArtist
	ctx.Album = album

	tracks, err := h.db.AlbumTracks(albumArtist, album)
	if err != nil {
		return err
	}
	ctx.Tracks = tracks

	if len(ctx.Tracks) == 0 {
		return ErrNotFound
	}

	for _, track := range tracks {
		if track.Artist != track.AlbumArtist {
			ctx.ShowArtists = true
			break
		}
	}

	calcSegment := func(play data.PlayCount) int {
		return play.Year*12 + int(play.Month)
	}

	plays := h.db.AlbumPlays(albumArtist, album)
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

	return h.templates.ExecuteTemplate(w, "album.gotmpl", ctx)
}
