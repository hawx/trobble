package main

import (
	"flag"
	"strconv"

	"github.com/hawx/trobble/data"
	"github.com/shkh/lastfm-go/lastfm"
)

const helpMessage = `Usage: trobble-import-lastfm [--db]

  Imports data from a last.fm profile.

    --db <path>     # Path to trobble db (default: 'trobble.db')
    --help          # Display this message
`

var (
	apiKey   = flag.String("api-key", "", "")
	secret   = flag.String("secret", "", "")
	username = flag.String("username", "", "")
	from     = flag.Int64("from", "", "")

	dbPath = flag.String("db", "trobble.db", "")
	help   = flag.Bool("help", false, "")
)

func main() {
	flag.Parse()

	api := lastfm.New(*apiKey, *secret)

	var resp UserGetRecentTracks
	page := 1

	for page := 1; page < resp.TotalPages; page++ {
		resp = api.User.GetRecentTracks(lastfm.P{
			"user": *username,
			"from": *from,
			"page": page,
		})

		for _, track := range resp.Tracks {
			n, _ := strconv.ParseInt(track.Date.Uts, 10, 64)

			scrobble := data.Scrobble{
				Artist:    track.Artist.Name,
				Album:     track.Album.Name,
				Track:     track.Name,
				Timestamp: n,
			}

		}
	}
}
