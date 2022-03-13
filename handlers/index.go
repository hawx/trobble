package handlers

import (
	"html/template"
	"net/http"

	"hawx.me/code/route"
	"hawx.me/code/trobble/data"
)

type indexCtx struct {
	Title      string
	TotalPlays int
	NowPlaying *data.Playing

	RecentlyPlayed []data.Scrobble

	TopArtists struct {
		Overall, Year, Month, Week []data.Artist
	}

	TopTracks struct {
		Overall, Year, Month, Week []data.Track
	}
}

type indexHandler struct {
	db        *data.Database
	title     string
	templates *template.Template
}

func Index(db *data.Database, title string, templates *template.Template) route.Handler {
	return &indexHandler{db, title, templates}
}

func (h indexHandler) ServeErrorHTTP(w http.ResponseWriter, r *http.Request) error {
	w.Header().Add("Content-Type", "text/html")

	ctx := indexCtx{}
	ctx.Title = h.title
	ctx.TotalPlays = h.db.TotalPlays()
	if nowPlaying, ok := h.db.GetNowPlaying(); ok {
		ctx.NowPlaying = nowPlaying
	}

	recent, err := h.db.RecentlyPlayed()
	if err != nil {
		return err
	}
	ctx.RecentlyPlayed = recent

	topArtists, err := h.db.TopArtists(10)
	if err != nil {
		return err
	}

	ctx.TopArtists.Overall = topArtists.Overall
	ctx.TopArtists.Year = topArtists.Year
	ctx.TopArtists.Month = topArtists.Month
	ctx.TopArtists.Week = topArtists.Week

	topTracks, err := h.db.TopTracks(10)
	if err != nil {
		return err
	}

	ctx.TopTracks.Overall = topTracks.Overall
	ctx.TopTracks.Year = topTracks.Year
	ctx.TopTracks.Month = topTracks.Month
	ctx.TopTracks.Week = topTracks.Week

	return h.templates.ExecuteTemplate(w, "index.gotmpl", ctx)
}
