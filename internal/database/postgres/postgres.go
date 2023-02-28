package postgres

import (
	"database/sql"
	"github.com/v1shn3vsk7/PlaylistAPI/pkg/song"
)

func NewDb(url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func GetSongs(db *sql.DB) []song.Song {

}
