package handlers

import (
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"hawx.me/code/assert"
	"hawx.me/code/trobble/data"
)

type fakeScrobbleDB struct{}

func (db *fakeScrobbleDB) NowPlaying(playing data.Playing) error {
	return nil
}

func (db *fakeScrobbleDB) Add(scrobble data.Scrobble) error {
	return nil
}

func TestScrobbleGetMobileSession(t *testing.T) {
	assert := assert.New(t)
	auth := NewAuth("testuser", "testkey", "testsecret")
	db := &fakeScrobbleDB{}

	s := httptest.NewServer(Scrobble(auth, db))
	defer s.Close()

	resp, err := http.PostForm(s.URL, url.Values{
		"method":   {"auth.getMobileSession"},
		"api_key":  {"testkey"},
		"username": {"testuser"},
		"api_sig":  {"ccedc18325ea9a841a5a1d50027f1941"},
	})
	assert.Nil(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	var v sessionResponse
	assert.Nil(xml.NewDecoder(resp.Body).Decode(&v))

	assert.Equal("ok", v.Status)
	assert.Equal("testuser", v.Session.Name)
	assert.Equal(0, v.Session.Subscriber)
}

func TestScrobbleGetMobileSessionWithUnknownUsername(t *testing.T) {
	assert := assert.New(t)
	auth := NewAuth("testuser", "testkey", "testsecret")
	db := &fakeScrobbleDB{}

	s := httptest.NewServer(Scrobble(auth, db))
	defer s.Close()

	resp, err := http.PostForm(s.URL, url.Values{
		"method":   {"auth.getMobileSession"},
		"api_key":  {"testkey"},
		"username": {"who"},
		"api_sig":  {"9896ae52b0a2586ccf9d52bae55bc709"},
	})
	assert.Nil(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	var v errorResponse
	assert.Nil(xml.NewDecoder(resp.Body).Decode(&v))

	assert.Equal("failed", v.Status)
	assert.Equal(4, v.Error.Code)
	assert.Equal("Authentication Failed", v.Error.Text)
}

func TestScrobbleUpdateNowPlaying(t *testing.T) {
	assert := assert.New(t)
	auth := NewAuth("testuser", "testkey", "testsecret")
	auth.sessionId = "39ba0d7ed7834b4e8498a0c463cba6ed"
	db := &fakeScrobbleDB{}

	s := httptest.NewServer(Scrobble(auth, db))
	defer s.Close()

	resp, err := http.PostForm(s.URL, url.Values{
		"method":      {"track.updateNowPlaying"},
		"api_key":     {"testkey"},
		"sk":          {auth.sessionId},
		"artist":      {"who that"},
		"album":       {"what that"},
		"albumArtist": {"group that"},
		"track":       {"sound that"},
		"api_sig":     {"928ee05c0603e8e4916054ab10aaf04f"},
	})
	assert.Nil(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	var v playingResponse
	assert.Nil(xml.NewDecoder(resp.Body).Decode(&v))

	assert.Equal("who that", v.NowPlaying.Artist)
	assert.Equal("what that", v.NowPlaying.Album)
	assert.Equal("group that", v.NowPlaying.AlbumArtist)
	assert.Equal("sound that", v.NowPlaying.Track)
}
