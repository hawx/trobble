package handlers

import (
	"net/http"
	"time"

	"github.com/hawx/trobble/data"
	"github.com/hawx/trobble/views"
)

type playedCtx struct {
	TotalPlays     int
	RecentlyPlayed []data.Scrobble

	TopArtists struct {
		Overall, Year, Month, Week []data.Artist
	}

	TopTracks struct {
		Overall, Year, Month, Week []data.Track
	}
}

type playedHandler struct {
	db *data.Database
}

func Played(db *data.Database) http.Handler {
	return &playedHandler{db}
}

// Simplified constants
const (
	Week  = 7 * 24 * time.Hour
	Month = 4 * Week
	Year  = 52 * Week
)

func (h playedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")

	ctx := playedCtx{}
	ctx.TotalPlays = h.db.TotalPlays()
	ctx.RecentlyPlayed = h.db.RecentlyPlayed()

	ctx.TopArtists.Overall = h.db.TopArtists(10)
	ctx.TopArtists.Year = h.db.TopArtistsAfter(10, -Year)
	ctx.TopArtists.Month = h.db.TopArtistsAfter(10, -Month)
	ctx.TopArtists.Week = h.db.TopArtistsAfter(10, -Week)

	ctx.TopTracks.Overall = h.db.TopTracks(10)
	ctx.TopTracks.Year = h.db.TopTracksAfter(10, -Year)
	ctx.TopTracks.Month = h.db.TopTracksAfter(10, -Month)
	ctx.TopTracks.Week = h.db.TopTracksAfter(10, -Week)

	views.Played.Execute(w, ctx)
}
