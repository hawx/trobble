package handlers

import (
	"net/http"
	"time"

	"hawx.me/code/trobble/data"
	"hawx.me/code/trobble/views"
)

type Date struct {
	Year  int
	Month time.Month
	Day   int
}

type DatedTracks struct {
	Date   Date
	Tracks []data.Scrobble
}

type playedCtx struct {
	Title      string
	TotalPlays int
	NowPlaying *data.Playing

	Tracks   []DatedTracks
	MoreTime string
}

type playedHandler struct {
	db    *data.Database
	title string
}

func Played(db *data.Database, title string) http.Handler {
	return &playedHandler{db, title}
}

func (h playedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")

	fromTime := time.Now()
	if fromStr := r.FormValue("from"); fromStr != "" {
		t, err := time.Parse(time.RFC3339, fromStr)
		if err == nil {
			fromTime = t
		}
	}

	ctx := playedCtx{}
	ctx.Title = h.title
	ctx.TotalPlays = h.db.TotalPlays()
	if nowPlaying, ok := h.db.GetNowPlaying(); ok {
		ctx.NowPlaying = nowPlaying
	}

	ctx.Tracks = []DatedTracks{}

	tracks := h.db.Played(fromTime)
	var ltracks []data.Scrobble

	year, month, day := time.Unix(tracks[0].Timestamp, 0).Date()
	ldate := Date{year, month, day}

	for _, track := range tracks {
		year, month, day := time.Unix(track.Timestamp, 0).Date()
		date := Date{year, month, day}

		if date != ldate {
			if len(ltracks) > 0 {
				ctx.Tracks = append(ctx.Tracks, DatedTracks{ldate, ltracks})
			}

			ldate = date
			ltracks = []data.Scrobble{}
		} else {
			ltracks = append(ltracks, track)
		}
	}

	if len(ltracks) > 0 {
		ctx.Tracks = append(ctx.Tracks, DatedTracks{ldate, ltracks})
	}

	ctx.MoreTime = time.Unix(tracks[len(tracks)-1].Timestamp, 0).Format(time.RFC3339)

	views.Played.Execute(w, ctx)
}
