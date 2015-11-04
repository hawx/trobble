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

const helpMessage = `Usage: trobble [--db] [--port|--socket]

  Catches messages from last.fm scrobblers and stores them in a database
  instead. It does not forward requests to last.fm.

    --username VAL     # Username to use
    --api-key VAL      # API Key used by connecting clients
    --secret VAL       # Secret used by connecting clients

    --title TITLE      # Title of page (default: 'trobble')
    --url URL          # Url to host (default: 'http://localhost:8080/')

    --db PATH          # Path to sqlite3 db (default: 'trobble.db')
    --port PORT        # Port to serve on (default: '8080')
    --socket SOCK      # Socket to serve on
    --help             # Display this message
`

var (
	username = flag.String("username", "", "")
	apiKey   = flag.String("api-key", "", "")
	secret   = flag.String("secret", "", "")

	title = flag.String("title", "trobble", "")
	url   = flag.String("url", "http://localhost:8080/", "")

	dbPath = flag.String("db", "trobble.db", "")
	port   = flag.String("port", "8080", "")
	socket = flag.String("socket", "", "")
	help   = flag.Bool("help", false, "")
)

func main() {
	flag.Parse()

	if *help {
		fmt.Println(helpMessage)
		return
	}

	db, err := data.Open(*dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	route.Handle("/", handlers.Index(db, *title))
	route.Handle("/feed", handlers.Feed(db, *title, *url))
	route.Handle("/played", handlers.Played(db, *title))
	route.Handle("/artist/:name", handlers.Artist(db, *title))
	route.HandleFunc("/styles.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/css")
		fmt.Fprint(w, views.Styles)
	})

	auth := handlers.NewAuth(*username, *apiKey, *secret)
	route.Handle("/scrobble/*any", handlers.Scrobble(auth, db))

	serve.Serve(*port, *socket, route.Default)
}
