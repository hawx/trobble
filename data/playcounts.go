package data

import "time"

type PlayCount struct {
	Year  int
	Month time.Month
	Count int
}

func (d *Database) TotalPlays() (count int, err error) {
	err = d.db.QueryRow("SELECT COUNT(*) FROM scrobbles").Scan(&count)
	return
}

func (d *Database) ArtistPlays(name string) (plays []PlayCount, err error) {
	rows, err := d.db.Query("SELECT COUNT(Timestamp), strftime('%Y', Timestamp, 'unixepoch'), strftime('%m', Timestamp, 'unixepoch') "+
		"FROM scrobbles "+
		"WHERE Artist = ? "+
		"GROUP BY strftime('%Y-%m', Timestamp, 'unixepoch')",
		name)

	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var count, year, month int
		if err = rows.Scan(&count, &year, &month); err != nil {
			return
		}

		plays = append(plays, PlayCount{
			Year:  year,
			Month: time.Month(month),
			Count: count,
		})
	}

	err = rows.Err()
	return
}
