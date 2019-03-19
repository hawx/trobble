package handlers

import (
	"html/template"
	"log"
	"net/http"
	"time"

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

func Index(db *data.Database, title string, templates *template.Template) http.Handler {
	return &indexHandler{db, title, templates}
}

// Simplified constants
const (
	Week  = 7 * 24 * time.Hour
	Month = 4 * Week
	Year  = 52 * Week
)

func (h indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")

	ctx := indexCtx{}
	ctx.Title = h.title
	ctx.TotalPlays = h.db.TotalPlays()
	if nowPlaying, ok := h.db.GetNowPlaying(); ok {
		ctx.NowPlaying = nowPlaying
	}

	ctx.RecentlyPlayed = h.db.RecentlyPlayed()

	ctx.TopArtists.Overall = h.db.TopArtists(10)
	ctx.TopArtists.Year = h.db.TopArtistsAfter(10, -Year)
	ctx.TopArtists.Month = h.db.TopArtistsAfter(10, -Month)
	ctx.TopArtists.Week = h.db.TopArtistsAfter(10, -Week)

	ctx.TopTracks.Overall = h.db.TopTracks(10)
	ctx.TopTracks.Year = h.db.TopTracksAfter(10, -Year)
	ctx.TopTracks.Month = h.db.TopTracksAfter(10, -Month)
	ctx.TopTracks.Week = h.db.TopTracksAfter(10, -Week)

	if err := h.templates.ExecuteTemplate(w, "index.gotmpl", ctx); err != nil {
		log.Println(err)
	}
}
