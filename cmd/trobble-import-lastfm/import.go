package main

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"

	"github.com/hawx/trobble/data"

	"flag"
	"fmt"
	"log"
)

const helpMessage = `Usage: trobble-import [--db] FILE

  Imports data from last.fm to a trobble database.

    --db <path>     # Path to trobble db (default: 'trobble.db')
    --help          # Display this message
`

var (
	dbPath = flag.String("db", "trobble.db", "")
	help   = flag.Bool("help", false, "")
)

type lastfmScrobble struct {
	Album            lastfmAlbum
	UncorrectedTrack lastfmTrack
	Track            lastfmTrack
	Timestamp        lastfmTimestamp
}

type lastfmAlbum struct {
	Artist     lastfmArtist
	Mbid, Name string
}

type lastfmTrack struct {
	Artist     lastfmArtist
	Mbid, Name string
}

type lastfmArtist struct {
	Mbid, Name string
}

type lastfmTimestamp struct {
	UnixTimestamp uint64
	Iso           string
}

// ISO time, unixtime, track name, track mbid, artist name, artist mbid, uncorrected track name, uncorrected track mbid, uncorrected artist name, uncorrected artist mbid, album name, album mbid, album artist name, album artist mbid, application
//
// [2008-12-10T16:19:27 1228925967 The Number Song d6efd638-2cef-4da1-85c8-7bcc20bc9746 DJ Shadow efa2c11a-1a35-4b60-bc1b-66d37de88511 The Number Song d6efd638-2cef-4da1-85c8-7bcc20bc9746 DJ Shadow efa2c11a-1a35-4b60-bc1b-66d37de88511 Endtroducing...  DJ Shadow efa2c11a-1a35-4b60-bc1b-66d37de88511 ]
type record struct {
	IsoTime               string
	Unixtime              string
	TrackName             string
	TrackMbid             string
	ArtistName            string
	ArtistMbid            string
	UncorrectedTrackName  string
	UncorrectedTrackMbid  string
	UncorrectedArtistName string
	UncorrectedArtistMbid string
	AlbumName             string
	AlbumMbid             string
	AlbumArtistName       string
	AlbumArtistMbid       string
	Application           string
}

func newRecord(vals []string) *record {
	return &record{
		IsoTime:               vals[0],
		Unixtime:              vals[1],
		TrackName:             vals[2],
		TrackMbid:             vals[3],
		ArtistName:            vals[4],
		ArtistMbid:            vals[5],
		UncorrectedTrackName:  vals[6],
		UncorrectedTrackMbid:  vals[7],
		UncorrectedArtistName: vals[8],
		UncorrectedArtistMbid: vals[9],
		AlbumName:             vals[10],
		AlbumMbid:             vals[11],
		AlbumArtistName:       vals[12],
		AlbumArtistMbid:       vals[13],
		Application:           vals[14],
	}
}

type recordReader struct {
	tsv *csv.Reader
}

func newReader(r io.Reader) (*recordReader, error) {
	tsv := csv.NewReader(r)
	tsv.Comma = '\t'
	tsv.Comment = '#'
	tsv.LazyQuotes = true
	tsv.TrailingComma = true
	tsv.TrimLeadingSpace = false

	_, err := tsv.Read() // skip header
	if err != nil {
		return nil, err
	}

	return &recordReader{tsv}, nil
}

func (r *recordReader) Read() (*record, error) {
	rcrd, err := r.tsv.Read()
	if err != nil {
		return nil, err
	}

	return newRecord(rcrd), nil
}

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

	file, _ := os.Open(flag.Args()[0])
	defer file.Close()

	r, err := newReader(file)
	if err != nil {
		log.Fatal(err)
	}

	scrobbles := []data.Scrobble{}
	timestamps := map[int64]data.Scrobble{}
	for {
		rcrd, err := r.Read()
		if err == io.EOF {
			log.Println("Reached last record")
			break
		}

		t, err := strconv.ParseInt(rcrd.Unixtime, 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		scrobble := data.Scrobble{
			Artist:      rcrd.ArtistName,
			AlbumArtist: rcrd.AlbumArtistName,
			Album:       rcrd.AlbumName,
			Track:       rcrd.TrackName,
			Timestamp:   t,
		}

		if _, ok := timestamps[t]; ok {
			continue
		}

		scrobbles = append(scrobbles, scrobble)
		timestamps[t] = scrobble
	}

	if err = db.AddMultiple(scrobbles); err != nil {
		log.Fatal(err)
	}
}
