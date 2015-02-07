package main

import (
	_ "code.google.com/p/go-sqlite/go1/sqlite3"
	"github.com/hawx/serve"

	"strings"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

const helpMessage = `Usage: trobble [--db] [--port|--socket]

  Catches messages from last.fm scrobblers (ymmv) and stores them in
  a database instead.

    --db <path>        # Path to sqlite3 db (default: 'trobble.db')
    --port <port>      # Port to serve on (default: '8080')
    --socket <sock>    # Socket to serve on
    --help             # Display this message
`

var (
	dbPath = flag.String("db", "trobble.db", "")
	port   = flag.String("port", "8080", "")
	socket = flag.String("socket", "", "")
	help   = flag.Bool("help", false, "")
)

func mustParseInt(val string) int64 {
	n, _ := strconv.ParseInt(val, 10, 64)
	return n
}

type Scrobble struct {
	Artist, AlbumArtist, Album, Track string
	Duration, Timestamp               int64
}

func (scrobble Scrobble) Respond(w io.Writer) {
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

type Playing struct {
	Artist, AlbumArtist, Album, Track string
	Duration                          int64
}

func (playing Playing) Respond(w io.Writer) {
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

type Database struct {
	db *sql.DB
}

func OpenDatabase(path string) (*Database, error) {
	sqlite, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	db := &Database{sqlite}

	return db, db.setup()
}

func (d *Database) setup() error {
	_, err := d.db.Exec(`CREATE TABLE IF NOT EXISTS scrobbles (
    Artist      TEXT,
    AlbumArtist TEXT,
    Album       TEXT,
    Track       TEXT,
    Duration    INTEGER,
    Timestamp   INTEGER PRIMARY KEY
  )`)

	return err
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) Add(scrobble Scrobble) error {
	_, err := d.db.Exec("INSERT INTO scrobbles VALUES(?, ?, ?, ?, ?, ?)",
		scrobble.Artist,
		scrobble.AlbumArtist,
		scrobble.Album,
		scrobble.Track,
		scrobble.Duration,
		scrobble.Timestamp)

	return err
}

type scrobbleHandler struct {
	db *Database
}

func ScrobbleHandler(db *Database) http.Handler {
	return &scrobbleHandler{db}
}

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

		playing := Playing{
			Artist:      r.FormValue("artist"),
			Album:       r.FormValue("album"),
			AlbumArtist: r.FormValue("albumArtist"),
			Track:       r.FormValue("track"),
			Duration:    mustParseInt(r.FormValue("duration")),
		}

		log.Println("now playing:", playing)
		playing.Respond(w)

	case "track.scrobble":
		scrobble := Scrobble{
			Artist:      r.FormValue("artist"),
			Album:       r.FormValue("album"),
			AlbumArtist: r.FormValue("albumArtist"),
			Track:       r.FormValue("track"),
			Duration:    mustParseInt(r.FormValue("duration")),
			Timestamp:   mustParseInt(r.FormValue("timestamp")),
		}

		log.Println("scrobbled:", scrobble)
		if err := handler.db.Add(scrobble); err != nil {
			log.Println(err)
		}
		scrobble.Respond(w)
	}
}

func main() {
	flag.Parse()

	if *help {
		fmt.Println(helpMessage)
		return
	}

	db, err := OpenDatabase(*dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.Handle("/", ScrobbleHandler(db))
	serve.Serve(*port, *socket, http.DefaultServeMux)
}
