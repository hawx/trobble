package data

type Scrobble struct {
	Artist, AlbumArtist, Album, Track string
	Timestamp                         int64
}

type Playing struct {
	Artist, AlbumArtist, Album, Track string
}

type Artist struct {
	Artist string
	Count  int
}

type Track struct {
	Artist string
	Track  string
	Count  int
}
