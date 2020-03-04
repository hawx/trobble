package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"hawx.me/code/route"
	"hawx.me/code/serve"
	"hawx.me/code/trobble/data"
	"hawx.me/code/trobble/handlers"
	"hawx.me/code/trobble/views"
)

const helpMessage = `Usage: trobble [options]

  Catches messages from last.fm scrobblers and stores them in a database
  instead. It does not forward requests to last.fm.

    --username VAL     # Username to use
    --api-key VAL      # API Key used by connecting clients
    --secret VAL       # Secret used by connecting clients

    --title TITLE      # Title of page (default: 'trobble')
    --url URL          # URL running at (default: 'http://localhost:8080/')
    --web PATH         # Path to 'web' directory (default: 'web')
    --db PATH          # Path to sqlite3 db (default: 'trobble.db')

    --port PORT        # Port to serve on (default: '8080')
    --socket SOCK      # Socket to serve on
    --help             # Display this message
`

func main() {
	var (
		username = flag.String("username", "", "")
		apiKey   = flag.String("api-key", "", "")
		secret   = flag.String("secret", "", "")

		title   = flag.String("title", "trobble", "")
		url     = flag.String("url", "http://localhost:8080/", "")
		webPath = flag.String("web", "web", "")
		dbPath  = flag.String("db", "trobble.db", "")

		port   = flag.String("port", "8080", "")
		socket = flag.String("socket", "", "")
	)

	flag.Usage = func() {
		fmt.Println(helpMessage)
	}
	flag.Parse()

	templates, err := views.Parse(*webPath + "/template/*.gotmpl")
	if err != nil {
		log.Fatal(err)
	}

	db, err := data.Open(*dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	route.Handle("/", handlers.Index(db, *title, templates))
	route.Handle("/feed", handlers.Feed(db, *title, *url))
	route.Handle("/played", handlers.Played(db, *title, templates))
	route.Handle("/listen/:timestamp", handlers.Listen(db, *title, templates))
	route.Handle("/artist/:artist", handlers.Artist(db, *title, templates))
	route.Handle("/album/:albumArtist/:album", handlers.Album(db, *title, templates))
	route.Handle("/track/:albumArtist/:album/:track", handlers.Track(db, *title, templates))

	route.Handle("/public/*path", http.StripPrefix("/public", http.FileServer(http.Dir(*webPath+"/static"))))

	auth := handlers.NewAuth(*username, *apiKey, *secret)
	route.Handle("/scrobble/*any", handlers.Scrobble(auth, db))

	serve.Serve(*port, *socket, serve.Recover(route.Default))
}
