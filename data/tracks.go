package data

import (
	"database/sql"
	"time"
)

type TrackRankings struct {
	Overall, Year, Month, Week []Track
}

func (d *Database) TopTracks(limit int) (rankings TrackRankings, err error) {
	overall, err := d.topTracksAfter(limit, -Forever)
	if err != nil {
		return
	}

	year, err := d.topTracksAfter(limit, -Year)
	if err != nil {
		return
	}

	month, err := d.topTracksAfter(limit, -Month)
	if err != nil {
		return
	}

	week, err := d.topTracksAfter(limit, -Week)
	if err != nil {
		return
	}

	return TrackRankings{
		Overall: overall,
		Year:    year,
		Month:   month,
		Week:    week,
	}, nil
}

type Track struct {
	Artist      string
	AlbumArtist string
	Album       string
	Track       string
	Count       int
}

func (d *Database) topTracksAfter(limit int, after time.Duration) (tracks []Track, err error) {
	rows, err := d.db.Query("SELECT Artist, AlbumArtist, Album, Track, COUNT(*) AS C FROM scrobbles WHERE Timestamp > ? GROUP BY Artist, Track ORDER BY C DESC LIMIT ?",
		time.Now().Add(after).Unix(),
		limit)

	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var track Track
		if err = rows.Scan(&track.Artist, &track.AlbumArtist, &track.Album, &track.Track, &track.Count); err != nil {
			return
		}
		tracks = append(tracks, track)
	}

	err = rows.Err()
	return
}

func (d *Database) ArtistTopTracks(name string, limit int) (tracks []Track, err error) {
	rows, err := d.db.Query("SELECT Artist, AlbumArtist, Album, Track, COUNT(*) AS C FROM scrobbles WHERE (Artist = ? OR AlbumArtist = ?) GROUP BY Artist, Track ORDER BY C DESC LIMIT ?",
		name,
		name,
		limit)

	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var track Track
		if err = rows.Scan(&track.Artist, &track.AlbumArtist, &track.Album, &track.Track, &track.Count); err != nil {
			return
		}
		tracks = append(tracks, track)
	}

	err = rows.Err()
	return
}

func (d *Database) AlbumTracks(albumArtist, album string) (tracks []Track, err error) {
	rows, err := d.db.Query("SELECT Artist, AlbumArtist, Album, Track, COUNT(*) AS C FROM scrobbles WHERE AlbumArtist = ? AND Album = ? GROUP BY Artist, Track ORDER BY C DESC",
		albumArtist,
		album)

	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var track Track
		if err = rows.Scan(&track.Artist, &track.AlbumArtist, &track.Album, &track.Track, &track.Count); err != nil {
			return
		}
		tracks = append(tracks, track)
	}

	err = rows.Err()
	return
}

func (d *Database) GetTrack(albumArtist, album, track string) (*Track, error) {
	var result Track

	err := d.db.QueryRow(`
    SELECT Artist, Album, AlbumArtist, Track, COUNT(*) AS C
    FROM scrobbles
    WHERE AlbumArtist = ? AND Album = ? AND Track = ?
    GROUP BY AlbumArtist, Track
    ORDER BY C DESC`,
		albumArtist,
		album,
		track).Scan(&result.Artist, &result.Album, &result.AlbumArtist, &result.Track, &result.Count)

	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	}

	return &result, nil
}
