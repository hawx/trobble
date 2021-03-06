package data

import (
	"time"

	_ "github.com/mattn/go-sqlite3"

	"database/sql"
)

// Simplified constants
const (
	Week    = 7 * 24 * time.Hour
	Month   = 4 * Week
	Year    = 52 * Week
	Forever = 100 * Year
)

type Database struct {
	db *sql.DB
}

func Open(path string) (*Database, error) {
	sqlite, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	db := &Database{sqlite}

	return db, db.setup()
}

func (d *Database) setup() error {
	_, err := d.db.Exec(`CREATE TABLE IF NOT EXISTS scrobbles (
    Artist      TEXT,
    AlbumArtist TEXT,
    Album       TEXT,
    Track       TEXT,
    Timestamp   INTEGER PRIMARY KEY
  );

  CREATE TABLE IF NOT EXISTS nowplaying (
    Id          INTEGER PRIMARY KEY,
    Artist      TEXT,
    AlbumArtist TEXT,
    Album       TEXT,
    Track       TEXT,
    Timestamp   INTEGER
  );

  INSERT OR IGNORE INTO nowplaying VALUES(1, "", "", "", "", 0);

  UPDATE scrobbles SET AlbumArtist = Artist WHERE AlbumArtist = '';
`)

	return err
}

func (d *Database) Close() error {
	return d.db.Close()
}
