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

  Catches messages from last.fm scrobblers and stores them in a database
  instead. It does not forward requests to last.fm.

    --username <val>   # Username to use
    --api-key <val>    # API Key used by connecting clients
    --secret <val>     # Secret used by connecting clients

    --db <path>        # Path to sqlite3 db (default: 'trobble.db')
    --port <port>      # Port to serve on (default: '8080')
    --socket <sock>    # Socket to serve on
    --help             # Display this message
`

var (
	username = flag.String("username", "", "")
	apiKey   = flag.String("api-key", "", "")
	secret   = flag.String("secret", "", "")

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

	auth := handlers.NewAuth(*username, *apiKey, *secret)
	http.Handle("/", handlers.Played(db))
	http.Handle("/scrobble/", handlers.Scrobble(auth, db))
	serve.Serve(*port, *socket, http.DefaultServeMux)
}
