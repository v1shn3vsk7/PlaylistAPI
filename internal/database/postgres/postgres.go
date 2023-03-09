package postgres

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"github.com/v1shn3vsk7/PlaylistAPI/pkg/song"
	"time"
)

type Postgres struct {
	*sql.DB
}

func NewDb(url string) (*Postgres, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Postgres{db}, nil
}

func (db *Postgres) GetSongs() ([]*song.Song, error) {
	rows, err := db.Query("SELECT name, artist, duration FROM songs ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songs []*song.Song

	for rows.Next() {
		var song = &song.Song{}
		if err := rows.Scan(&song.Name, &song.Artist,
			&song.Duration); err != nil {
			return nil, err
		}
		song.Duration *= time.Second
		songs = append(songs, song)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return songs, nil
}

func (db *Postgres) AddSong(song *song.Song) error {
	if db == nil {
		return errors.New("db connections is nil")
	}
	if _, err := db.Query("INSERT INTO songs (name, artist, duration) VALUES ($1, $2, $3)",
		song.Name, song.Artist, song.Duration.Seconds()); err != nil {
		return err
	}

	return nil
}

func (db *Postgres) FindSong(song *song.Song) (int, error) {
	var id int

	if err := db.QueryRow("SELECT id FROM songs WHERE name = $1 AND artist = $2",
		song.Name, song.Artist).Scan(&id); err != nil {
		return 0, nil
	}

	return id, nil
}

func (db *Postgres) EditSong(song *song.Song, id int) {
	db.QueryRow("UPDATE songs SET name = $1, artist = $2, duration = $3 WHERE id = $4",
		song.Name, song.Artist, song.Duration.Seconds(), id)
}

func (db *Postgres) DeleteSong(song *song.Song) error {
	result, err := db.Exec("DELETE FROM songs WHERE name = $1 AND artist = $2", song.Name, song.Artist)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("song not found")
	}

	return nil
}
