package data

import "time"

type Track struct {
	Artist string
	Track  string
	Count  int
}

func (d *Database) TopTracks(limit int) ([]Track, error) {
	return d.TopTracksAfter(limit, -100*365*24*time.Hour)
}

func (d *Database) TopTracksAfter(limit int, after time.Duration) (tracks []Track, err error) {
	rows, err := d.db.Query("SELECT Artist, Track, COUNT(*) AS C FROM scrobbles WHERE Timestamp > ? GROUP BY Artist, Track ORDER BY C DESC LIMIT ?",
		time.Now().Add(after).Unix(),
		limit)

	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var track Track
		if err = rows.Scan(&track.Artist, &track.Track, &track.Count); err != nil {
			return
		}
		tracks = append(tracks, track)
	}

	err = rows.Err()
	return
}

func (d *Database) ArtistTopTracks(name string, limit int) (tracks []Track, err error) {
	rows, err := d.db.Query("SELECT Artist, Track, COUNT(*) AS C FROM scrobbles WHERE Artist = ? GROUP BY Artist, Track ORDER BY C DESC LIMIT ?",
		name,
		limit)

	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var track Track
		if err = rows.Scan(&track.Artist, &track.Track, &track.Count); err != nil {
			return
		}
		tracks = append(tracks, track)
	}

	err = rows.Err()
	return
}
