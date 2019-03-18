package data

import (
	"time"
)

type Scrobble struct {
	Artist, AlbumArtist, Album, Track string
	Timestamp                         int64
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
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var scrobble Scrobble
		if err = rows.Scan(&scrobble.Artist, &scrobble.Album, &scrobble.AlbumArtist, &scrobble.Track, &scrobble.Timestamp); err != nil {
			panic(err)
		}
		scrobbles = append(scrobbles, scrobble)
	}

	if err = rows.Err(); err != nil {
		panic(err)
	}
	return
}

func (d *Database) Played(from time.Time) (scrobbles []Scrobble) {
	rows, err := d.db.Query("SELECT Artist, Album, AlbumArtist, Track, Timestamp FROM scrobbles WHERE Timestamp < ? ORDER BY Timestamp DESC LIMIT 100",
		from)

	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var scrobble Scrobble
		if err = rows.Scan(&scrobble.Artist, &scrobble.Album, &scrobble.AlbumArtist, &scrobble.Track, &scrobble.Timestamp); err != nil {
			panic(err)
		}
		scrobbles = append(scrobbles, scrobble)
	}

	if err = rows.Err(); err != nil {
		panic(err)
	}
	return
}

func (d *Database) Scrobble(at string) (scrobble Scrobble, ok bool) {
	row := d.db.QueryRow("SELECT Artist, Album, AlbumArtist, Track, Timestamp FROM scrobbles WHERE Timestamp = ?",
		at)

	if err := row.Scan(&scrobble.Artist, &scrobble.Album, &scrobble.AlbumArtist, &scrobble.Track, &scrobble.Timestamp); err != nil {
		return
	}

	return scrobble, true
}
