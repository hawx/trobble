package data

import "time"

type Artist struct {
	Artist string
	Count  int
}

func (d *Database) TopArtists(limit int) []Artist {
	return d.TopArtistsAfter(limit, -100*365*24*time.Hour)
}

func (d *Database) TopArtistsAfter(limit int, after time.Duration) (artists []Artist) {
	rows, err := d.db.Query("SELECT Artist, COUNT(*) AS C FROM scrobbles WHERE Timestamp > ? GROUP BY Artist ORDER BY C DESC LIMIT ?",
		time.Now().Add(after).Unix(),
		limit)

	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var artist Artist
		if err = rows.Scan(&artist.Artist, &artist.Count); err != nil {
			panic(err)
		}
		artists = append(artists, artist)
	}

	if err = rows.Err(); err != nil {
		panic(err)
	}
	return
}
