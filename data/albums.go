package data

type Album struct {
	Artist string
	Album  string
	Count  int
}

func (d *Database) ArtistTopAlbums(name string, limit int) (albums []Album) {
	rows, err := d.db.Query("SELECT Artist, Album, COUNT(*) AS C FROM scrobbles WHERE Artist = ? GROUP BY Artist, Album ORDER BY C DESC LIMIT ?",
		name,
		limit)

	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var album Album
		if err = rows.Scan(&album.Artist, &album.Album, &album.Count); err != nil {
			panic(err)
		}
		albums = append(albums, album)
	}

	if err = rows.Err(); err != nil {
		panic(err)
	}
	return
}
