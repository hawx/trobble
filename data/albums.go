package data

type Album struct {
	Artist string
	Album  string
	Count  int
}

func (d *Database) ArtistTopAlbums(name string, limit int) (albums []Album, err error) {
	rows, err := d.db.Query("SELECT Artist, Album, COUNT(*) AS C FROM scrobbles WHERE Artist = ? GROUP BY Artist, Album ORDER BY C DESC LIMIT ?",
		name,
		limit)

	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var album Album
		if err = rows.Scan(&album.Artist, &album.Album, &album.Count); err != nil {
			return
		}
		albums = append(albums, album)
	}

	err = rows.Err()
	return
}
