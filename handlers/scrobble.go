package handlers

import (
	"crypto/md5"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"hawx.me/code/trobble/data"
)

type errorResponse struct {
	XMLName xml.Name `xml:"lfm"`
	Status  string   `xml:"status,attr"`
	Error   struct {
		Code int    `xml:"code,attr"`
		Text string `xml:",chardata"`
	} `xml:"error"`
}

func filter(pred func(*http.Request) bool, sub http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if pred(r) {
			sub.ServeHTTP(w, r)
		} else {
			resp := new(errorResponse)
			resp.Status = "failed"
			resp.Error.Code = 4
			resp.Error.Text = "Authentication Failed"
			xml.NewEncoder(w).Encode(resp)
		}
	})
}

type Auth struct {
	username, apiKey, secret, sessionId string
}

func NewAuth(username, apiKey, secret string) Auth {
	return Auth{username, apiKey, secret, strings.Replace(uuid.New().String(), "-", "", -1)}
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

type sessionResponse struct {
	XMLName xml.Name `xml:"lfm"`
	Status  string   `xml:"status,attr"`
	Session struct {
		Name       string `xml:"name"`
		Key        string `xml:"key"`
		Subscriber int    `xml:"subscriber"`
	} `xml:"session"`
}

func (handler *scrobbleHandler) getMobileSession(w http.ResponseWriter, r *http.Request) {
	resp := new(sessionResponse)
	resp.Status = "ok"
	resp.Session.Name = r.FormValue("username")
	resp.Session.Key = handler.auth.sessionId
	resp.Session.Subscriber = 0

	xml.NewEncoder(w).Encode(resp)
}

type playingResponse struct {
	XMLName    xml.Name `xml:"lfm"`
	Status     string   `xml:"status,attr"`
	NowPlaying struct {
		Track          string `xml:"track"`
		Artist         string `xml:"artist"`
		Album          string `xml:"album"`
		AlbumArtist    string `xml:"albumArtist"`
		IgnoredMessage string `xml:"ignoredMessage"`
	} `xml:"nowplaying"`
}

func (handler *scrobbleHandler) updateNowPlaying(w http.ResponseWriter, r *http.Request) {
	playing := data.Playing{
		Artist:      r.FormValue("artist"),
		Album:       r.FormValue("album"),
		AlbumArtist: r.FormValue("albumArtist"),
		Track:       r.FormValue("track"),
	}

	log.Println("now playing:", playing)
	if err := handler.db.NowPlaying(playing); err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	resp := new(playingResponse)
	resp.Status = "ok"
	resp.NowPlaying.Track = playing.Track
	resp.NowPlaying.Artist = playing.Artist
	resp.NowPlaying.Album = playing.Album
	resp.NowPlaying.AlbumArtist = playing.AlbumArtist

	xml.NewEncoder(w).Encode(resp)
}

func formValue(r *http.Request, keys ...string) string {
	for _, key := range keys {
		if val := r.FormValue(key); val != "" {
			return val
		}
	}
	return ""
}

type scrobbleResponse struct {
	XMLName   xml.Name `xml:"lfm"`
	Status    string   `xml:"status,attr"`
	Scrobbles struct {
		Accepted int `xml:"accepted,attr"`
		Ignored  int `xml:"ignored,attr"`
		Scrobble struct {
			Track          string `xml:"track"`
			Artist         string `xml:"artist"`
			Album          string `xml:"album"`
			AlbumArtist    string `xml:"albumArtist"`
			Timestamp      int64  `xml:"timestamp"`
			IgnoredMessage string `xml:"ignoredMessage"`
		} `xml:"scrobble"`
	} `xml:"scrobbles"`
}

func (handler *scrobbleHandler) scrobble(w http.ResponseWriter, r *http.Request) {
	scrobble := data.Scrobble{
		Artist:      formValue(r, "artist", "artist[0]"),
		Album:       formValue(r, "album", "album[0]"),
		AlbumArtist: formValue(r, "albumArtist", "albumArtist[0]"),
		Track:       formValue(r, "track", "track[0]"),
		Timestamp:   mustParseInt(formValue(r, "timestamp", "timestamp[0]")),
	}

	log.Println("scrobbled:", scrobble)
	if err := handler.db.Add(scrobble); err != nil {
		log.Println(err) // maybe return error response???
	}

	resp := new(scrobbleResponse)
	resp.Status = "ok"
	resp.Scrobbles.Accepted = 1
	resp.Scrobbles.Ignored = 0
	resp.Scrobbles.Scrobble.Track = scrobble.Track
	resp.Scrobbles.Scrobble.Artist = scrobble.Artist
	resp.Scrobbles.Scrobble.Album = scrobble.Album
	resp.Scrobbles.Scrobble.AlbumArtist = scrobble.AlbumArtist
	resp.Scrobbles.Scrobble.Timestamp = scrobble.Timestamp

	xml.NewEncoder(w).Encode(resp)
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
