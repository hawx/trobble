package data

import "time"

type PlayCount struct {
	Year  int
	Month time.Month
	Count int
}

func (d *Database) TotalPlays() (count int) {
	err := d.db.QueryRow("SELECT COUNT(*) FROM scrobbles").Scan(&count)
	if err != nil {
		panic(err)
	}

	return
}

func (d *Database) ArtistPlays(name string) (plays []PlayCount) {
	rows, err := d.db.Query("SELECT COUNT(Timestamp), strftime('%Y', Timestamp, 'unixepoch'), strftime('%m', Timestamp, 'unixepoch') "+
		"FROM scrobbles "+
		"WHERE Artist = ? "+
		"GROUP BY strftime('%Y-%m', Timestamp, 'unixepoch')",
		name)

	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var count, year, month int
		if err = rows.Scan(&count, &year, &month); err != nil {
			panic(err)
		}

		plays = append(plays, PlayCount{
			Year:  year,
			Month: time.Month(month),
			Count: count,
		})
	}

	if err = rows.Err(); err != nil {
		panic(err)
	}
	return
}
