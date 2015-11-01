package handlers

import (
	"net/http"

	"hawx.me/code/trobble/data"
)

func Played(db *data.Database, title string) http.Handler {
	return http.NotFoundHandler()
}
