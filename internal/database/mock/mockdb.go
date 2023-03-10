package mockDb

import (
	"database/sql"
	"errors"
	"github.com/v1shn3vsk7/PlaylistAPI/pkg/song"
	"time"
)

type MockDb struct {
	*sql.DB
}

func (db *MockDb) GetSongs() ([]*song.Song, error) {
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

func (db *MockDb) AddSong(song *song.Song) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	if _, err := tx.Exec("INSERT INTO songs (name, artist, duration) VALUES ($1, $2, $3)",
		song.Name, song.Artist, int64(song.Duration.Seconds())); err != nil {
		return err
	}

	return nil
}

func (db *MockDb) FindSong(song *song.Song) (int, error) {
	var id int

	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	if _ = tx.QueryRow("SELECT id FROM songs WHERE name = $1 AND artist = $2",
		song.Name, song.Artist).Scan(&id); id != 0 {
		return 0, errors.New("Not found")
	}

	return id, nil
}

func (db *MockDb) EditSong(song *song.Song, id int) {
	tx, err := db.Begin()
	if err != nil {
		return
	}
	_ = tx.QueryRow("UPDATE songs SET name = $1, artist = $2, duration = $3 WHERE id = $4",
		song.Name, song.Artist, song.Duration.Seconds(), id)
}

func (db *MockDb) DeleteSong(song *song.Song) error {
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
