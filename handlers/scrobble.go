package handlers

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"code.google.com/p/go-uuid/uuid"
	"github.com/hawx/trobble/data"
)

func filter(pred func(*http.Request) bool, sub http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if pred(r) {
			sub.ServeHTTP(w, r)
		} else {
			fmt.Fprintf(w, `<lfm status="failed">
  <error code="4">Authentication Failed</error>
</lfm>`)
		}
	})
}

type Auth struct {
	username, apiKey, secret, sessionId string
}

func NewAuth(username, apiKey, secret string) Auth {
	return Auth{username, apiKey, secret, strings.Replace(uuid.New(), "-", "", -1)}
}

func (auth Auth) calcSignature(r *http.Request) string {
	keys := make([]string, 0, len(r.Form))
	for k := range r.Form {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	sigStr := ""
	for _, k := range keys {
		if k != "api_sig" {
			sigStr += k + r.FormValue(k)
		}
	}
	sigStr += auth.secret

	h := md5.New()
	io.WriteString(h, sigStr)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (auth Auth) checkApiDetails(r *http.Request) bool {
	return r.FormValue("api_key") == auth.apiKey && r.FormValue("api_sig") == auth.calcSignature(r)
}

func (auth Auth) checkUsername(r *http.Request) bool {
	return auth.checkApiDetails(r) && r.FormValue("username") == auth.username
}

func (auth Auth) checkSession(r *http.Request) bool {
	return auth.checkApiDetails(r) && r.FormValue("sk") == auth.sessionId
}

type scrobbleHandler struct {
	db         *data.Database
	auth       Auth
	responders map[string]http.Handler
}

func Scrobble(auth Auth, db *data.Database) http.Handler {
	handler := new(scrobbleHandler)
	handler.db = db
	handler.auth = auth
	handler.responders = map[string]http.Handler{
		"auth.getmobilesession":  filter(auth.checkUsername, http.HandlerFunc(handler.getMobileSession)),
		"track.updatenowplaying": filter(auth.checkSession, http.HandlerFunc(handler.updateNowPlaying)),
		"track.scrobble":         filter(auth.checkSession, http.HandlerFunc(handler.scrobble)),
	}
	return handler
}

func (handler *scrobbleHandler) getMobileSession(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<lfm status="ok">
  <session>
    <name>%s</name>
    <key>%s</key>
    <subscriber>0</subscriber>
  </session>
</lfm>`, r.FormValue("username"), handler.auth.sessionId)
}

func (handler *scrobbleHandler) updateNowPlaying(w http.ResponseWriter, r *http.Request) {
	playing := data.Playing{
		Artist:      r.FormValue("artist"),
		Album:       r.FormValue("album"),
		AlbumArtist: r.FormValue("albumArtist"),
		Track:       r.FormValue("track"),
	}

	log.Println("now playing:", playing)
	respondPlaying(playing, w)
}

func (handler *scrobbleHandler) scrobble(w http.ResponseWriter, r *http.Request) {
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

func (handler *scrobbleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if responder, ok := handler.responders[strings.ToLower(r.FormValue("method"))]; ok {
		responder.ServeHTTP(w, r)
	}
}

func mustParseInt(val string) int64 {
	n, _ := strconv.ParseInt(val, 10, 64)
	return n
}
