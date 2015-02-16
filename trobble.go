package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/hawx/serve"
	"github.com/hawx/trobble/data"
	"github.com/hawx/trobble/handlers"
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

	http.Handle("/", handlers.Played(db))
	http.Handle("/scrobble", handlers.Scrobble(db))
	serve.Serve(*port, *socket, http.DefaultServeMux)
}
