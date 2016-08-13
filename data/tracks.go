package data

import (
	"database/sql"
	"time"
)

type Track struct {
	Artist string
	Album  string
	Track  string
	Count  int
}

func (d *Database) TopTracks(limit int) []Track {
	return d.TopTracksAfter(limit, -100*365*24*time.Hour)
}

func (d *Database) TopTracksAfter(limit int, after time.Duration) (tracks []Track) {
	rows, err := d.db.Query("SELECT Artist, Album, Track, COUNT(*) AS C FROM scrobbles WHERE Timestamp > ? GROUP BY Artist, Track ORDER BY C DESC LIMIT ?",
		time.Now().Add(after).Unix(),
		limit)

	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var track Track
		if err = rows.Scan(&track.Artist, &track.Album, &track.Track, &track.Count); err != nil {
			panic(err)
		}
		tracks = append(tracks, track)
	}

	if err = rows.Err(); err != nil {
		panic(err)
	}
	return
}

func (d *Database) ArtistTopTracks(name string, limit int) (tracks []Track) {
	rows, err := d.db.Query("SELECT Artist, Album, Track, COUNT(*) AS C FROM scrobbles WHERE Artist = ? GROUP BY Artist, Track ORDER BY C DESC LIMIT ?",
		name,
		limit)

	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var track Track
		if err = rows.Scan(&track.Artist, &track.Album, &track.Track, &track.Count); err != nil {
			panic(err)
		}
		tracks = append(tracks, track)
	}

	if err = rows.Err(); err != nil {
		panic(err)
	}
	return
}

func (d *Database) AlbumTracks(artist, album string) (tracks []Track) {
	rows, err := d.db.Query("SELECT Artist, Album, Track, COUNT(*) AS C FROM scrobbles WHERE Artist = ? AND Album = ? GROUP BY Artist, Track ORDER BY C DESC",
		artist,
		album)

	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var track Track
		if err = rows.Scan(&track.Artist, &track.Album, &track.Track, &track.Count); err != nil {
			panic(err)
		}
		tracks = append(tracks, track)
	}

	if err = rows.Err(); err != nil {
		panic(err)
	}
	return
}

func (d *Database) GetTrack(artist, album, track string) *Track {
	var result Track

	err := d.db.QueryRow("SELECT Artist, Album, Track, COUNT(*) AS C FROM scrobbles WHERE Artist = ? AND Album = ? AND Track = ? GROUP BY Artist, Track ORDER BY C DESC",
		artist,
		album,
		track).Scan(&result.Artist, &result.Album, &result.Track, &result.Count)

	switch {
	case err == sql.ErrNoRows:
		return nil
	case err != nil:
		panic(err)
	}

	return &result
}
