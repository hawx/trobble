package data

import "time"

type ArtistRankings struct {
	Overall, Year, Month, Week []Artist
}

func (d *Database) TopArtists(limit int) (rankings ArtistRankings, err error) {
	overall, err := d.topArtistsAfter(limit, -Forever)
	if err != nil {
		return
	}

	year, err := d.topArtistsAfter(limit, -Year)
	if err != nil {
		return
	}

	month, err := d.topArtistsAfter(limit, -Month)
	if err != nil {
		return
	}

	week, err := d.topArtistsAfter(limit, -Week)
	if err != nil {
		return
	}

	return ArtistRankings{
		Overall: overall,
		Year:    year,
		Month:   month,
		Week:    week,
	}, nil
}

type Artist struct {
	Artist string
	Count  int
}

func (d *Database) topArtistsAfter(limit int, after time.Duration) (artists []Artist, err error) {
	rows, err := d.db.Query("SELECT Artist, COUNT(*) AS C FROM scrobbles WHERE Timestamp > ? GROUP BY Artist ORDER BY C DESC LIMIT ?",
		time.Now().Add(after).Unix(),
		limit)

	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var artist Artist
		if err = rows.Scan(&artist.Artist, &artist.Count); err != nil {
			return
		}
		artists = append(artists, artist)
	}

	err = rows.Err()
	return
}
