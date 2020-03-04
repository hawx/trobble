package data

type Album struct {
	AlbumArtist string
	Album       string
	Count       int
}

func (d *Database) ArtistTopAlbums(name string, limit int) (albums []Album, err error) {
	rows, err := d.db.Query("SELECT AlbumArtist, Album, COUNT(*) AS C FROM scrobbles WHERE (Artist = ? OR AlbumArtist = ?) GROUP BY AlbumArtist, Album ORDER BY C DESC LIMIT ?",
		name,
		name,
		limit)

	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var album Album
		if err = rows.Scan(&album.AlbumArtist, &album.Album, &album.Count); err != nil {
			return
		}

		if album.Album == "" {
			album.Album = "—"
		}
		albums = append(albums, album)
	}

	err = rows.Err()
	return
}
