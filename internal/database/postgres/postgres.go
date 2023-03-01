package postgres

import (
	"database/sql"
	_ "github.com/lib/pq"
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

func GetSongs(db *sql.DB) ([]song.Song, error) {
	rows, err := db.Query("SELECT name, artist, duration FROM songs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songs []song.Song

	for rows.Next() {
		var song song.Song
		if err := rows.Scan(&song.Name, &song.Artist,
			&song.Duration); err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return songs, nil
}

func AddSong(db *sql.DB, song *song.Song) error {
	if  _ ,err := db.Query("INSERT INTO songs (name, artist, duration) VALUES ($1, $2, $3)",
		song.Name, song.Artist, song.Duration); err != nil {
		return err
	}

	return nil
}