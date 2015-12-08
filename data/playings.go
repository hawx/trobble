package data

import "time"

type Playing struct {
	Artist, AlbumArtist, Album, Track string
}

func (d *Database) NowPlaying(playing Playing) error {
	_, err := d.db.Exec("UPDATE nowplaying SET Artist=?, AlbumArtist=?, Album=?, Track=?, Timestamp=? WHERE Id=1",
		playing.Artist,
		playing.AlbumArtist,
		playing.Album,
		playing.Track,
		time.Now().UTC().Unix())

	return err
}

func (d *Database) GetNowPlaying() (*Playing, bool, error) {
	row := d.db.QueryRow("SELECT Artist, AlbumArtist, Album, Track, Timestamp FROM nowplaying WHERE Id=1")

	var playing Playing
	var timestamp int64
	if err := row.Scan(&playing.Artist, &playing.AlbumArtist, &playing.Album, &playing.Track, &timestamp); err != nil {
		return nil, false, err
	}

	if time.Unix(timestamp, 0).After(time.Now().Add(-10 * time.Minute)) {
		return &playing, true, nil
	}

	return nil, false, nil
}
