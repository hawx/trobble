package handlers

import (
	"html/template"
	"net/http"
	"net/url"

	"hawx.me/code/route"
	"hawx.me/code/trobble/data"
)

type trackCtx struct {
	Title      string
	TotalPlays int
	NowPlaying *data.Playing

	Artist      string
	AlbumArtist string
	Album       string
	Track       string
	Plays       []int
	MaxPlays    int
	TrackPlays  int
}

type trackHandler struct {
	db        *data.Database
	title     string
	templates *template.Template
}

func Track(db *data.Database, title string, templates *template.Template) route.Handler {
	return &trackHandler{db, title, templates}
}

func (h trackHandler) ServeErrorHTTP(w http.ResponseWriter, r *http.Request) error {
	w.Header().Add("Content-Type", "text/html")

	vars := route.Vars(r)

	albumArtist, err := url.QueryUnescape(vars["albumArtist"])
	if err != nil {
		return ErrNotFound
	}

	album, err := url.QueryUnescape(vars["album"])
	if err != nil {
		return ErrNotFound
	}

	track, err := url.QueryUnescape(vars["track"])
	if err != nil {
		return ErrNotFound
	}

	foundTrack, err := h.db.GetTrack(albumArtist, album, track)
	if err != nil {
		return err
	}

	if foundTrack == nil {
		return ErrNotFound
	}

	ctx := trackCtx{}
	ctx.Title = h.title
	ctx.TotalPlays = h.db.TotalPlays()
	if nowPlaying, ok := h.db.GetNowPlaying(); ok {
		ctx.NowPlaying = nowPlaying
	}

	ctx.Artist = foundTrack.Artist
	ctx.Album = foundTrack.Album
	ctx.AlbumArtist = foundTrack.AlbumArtist
	ctx.Track = foundTrack.Track
	ctx.TrackPlays = foundTrack.Count

	calcSegment := func(play data.PlayCount) int {
		return play.Year*12 + int(play.Month)
	}

	plays := h.db.TrackPlays(albumArtist, album, track)
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

	return h.templates.ExecuteTemplate(w, "track.gotmpl", ctx)
}
