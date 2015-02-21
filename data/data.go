package data

import (
	"log"
	"time"

	_ "code.google.com/p/go-sqlite/go1/sqlite3"

	"database/sql"
)

type Database struct {
	db *sql.DB
}

func Open(path string) (*Database, error) {
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
    Timestamp   INTEGER PRIMARY KEY
  )`)

	return err
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) Add(scrobble Scrobble) error {
	_, err := d.db.Exec("INSERT INTO scrobbles VALUES(?, ?, ?, ?, ?)",
		scrobble.Artist,
		scrobble.AlbumArtist,
		scrobble.Album,
		scrobble.Track,
		scrobble.Timestamp)

	return err
}

func (d *Database) AddMultiple(scrobbles []Scrobble) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	for _, scrobble := range scrobbles {
		_, err := tx.Exec("INSERT OR IGNORE INTO scrobbles VALUES(?, ?, ?, ?, ?)",
			scrobble.Artist,
			scrobble.AlbumArtist,
			scrobble.Album,
			scrobble.Track,
			scrobble.Timestamp)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (d *Database) RecentlyPlayed() (scrobbles []Scrobble) {
	rows, err := d.db.Query("SELECT Artist, Album, AlbumArtist, Track, Timestamp FROM scrobbles ORDER BY Timestamp DESC LIMIT 10")
	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var scrobble Scrobble
		if err := rows.Scan(&scrobble.Artist, &scrobble.Album, &scrobble.AlbumArtist, &scrobble.Track, &scrobble.Timestamp); err != nil {
			log.Println(err)
			return
		}
		scrobbles = append(scrobbles, scrobble)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return
	}

	return
}

func (d *Database) TotalPlays() (count int) {
	err := d.db.QueryRow("SELECT COUNT(*) FROM scrobbles").Scan(&count)
	if err != nil {
		log.Println(err)
	}
	return
}

func (d *Database) TopArtists(limit int) (artists []Artist) {
	return d.TopArtistsAfter(limit, -100*365*24*time.Hour)
}

func (d *Database) TopArtistsAfter(limit int, after time.Duration) (artists []Artist) {
	rows, err := d.db.Query("SELECT Artist, COUNT(*) AS C FROM scrobbles WHERE Timestamp > ? GROUP BY Artist ORDER BY C DESC LIMIT ?",
		time.Now().Add(after).Unix(),
		limit)

	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var artist Artist
		if err := rows.Scan(&artist.Artist, &artist.Count); err != nil {
			log.Println(err)
			return
		}
		artists = append(artists, artist)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return
	}

	return
}

func (d *Database) TopTracks(limit int) (tracks []Track) {
	return d.TopTracksAfter(limit, -100*365*24*time.Hour)
}

func (d *Database) TopTracksAfter(limit int, after time.Duration) (tracks []Track) {
	rows, err := d.db.Query("SELECT Artist, Track, COUNT(*) AS C FROM scrobbles WHERE Timestamp > ? GROUP BY Artist, Track ORDER BY C DESC LIMIT ?",
		time.Now().Add(after).Unix(),
		limit)

	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var track Track
		if err := rows.Scan(&track.Artist, &track.Track, &track.Count); err != nil {
			log.Println(err)
			return
		}
		tracks = append(tracks, track)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return
	}

	return
}
