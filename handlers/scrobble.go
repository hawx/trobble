package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/hawx/trobble/data"
)

type scrobbleHandler struct {
	db *data.Database
}

func Scrobble(db *data.Database) http.Handler {
	return &scrobbleHandler{db}
}

func respondScrobble(scrobble data.Scrobble, w io.Writer) {
	fmt.Fprintf(w, `<?xml version='1.0' encoding='utf-8'?>
  <lfm status="ok">
    <scrobbles accepted="1" ignored="0">
      <scrobble>
        <track corrected="0">%s</track>
        <artist corrected="0">%s</artist>
        <album corrected="0">%s</album>
        <albumArtist corrected="0">%s</albumArtist>
        <timestamp>%d</timestamp>
        <ignoredMessage code="0"></ignoredMessage>
      </scrobble>
    </scrobbles>
  </lfm>`, scrobble.Track, scrobble.Artist, scrobble.Album, scrobble.AlbumArtist, scrobble.Timestamp)
}

func respondPlaying(playing data.Playing, w io.Writer) {
	fmt.Fprintf(w, `<?xml version='1.0' encoding='utf-8'?>
  <lfm status="ok">
    <nowplaying>
      <track corrected="0">%s</track>
      <artist corrected="0">%s</artist>
      <album corrected="0">%s</album>
      <albumArtist corrected="0">%s</albumArtist>
      <ignoredMessage code="0"></ignoredMessage>
   </nowplaying>
  </lfm>`, playing.Track, playing.Artist, playing.Album, playing.AlbumArtist)
}

// TODO: Deal with auth properly so this can be hosted on the web.
func (handler *scrobbleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	switch strings.ToLower(r.FormValue("method")) {
	case "auth.getmobilesession":
		fmt.Fprintf(w, `<lfm status="ok">
  <session>
    <name>%s</name>
    <key>d580d57f32848f5dcf574d1ce18d78b2</key>
    <subscriber>0</subscriber>
  </session>
</lfm>`, r.FormValue("username"))

	case "track.updatenowplaying":

		playing := data.Playing{
			Artist:      r.FormValue("artist"),
			Album:       r.FormValue("album"),
			AlbumArtist: r.FormValue("albumArtist"),
			Track:       r.FormValue("track"),
		}

		log.Println("now playing:", playing)
		respondPlaying(playing, w)

	case "track.scrobble":
		scrobble := data.Scrobble{
			Artist:      r.FormValue("artist"),
			Album:       r.FormValue("album"),
			AlbumArtist: r.FormValue("albumArtist"),
			Track:       r.FormValue("track"),
			Timestamp:   mustParseInt(r.FormValue("timestamp")),
		}

		log.Println("scrobbled:", scrobble)
		if err := handler.db.Add(scrobble); err != nil {
			log.Println(err)
		}
		respondScrobble(scrobble, w)
	}
}

func mustParseInt(val string) int64 {
	n, _ := strconv.ParseInt(val, 10, 64)
	return n
}
