package database

import "github.com/v1shn3vsk7/PlaylistAPI/pkg/song"

type Database interface {
	GetSongs() ([]*song.Song, error)
	AddSong(song *song.Song) error
	FindSong(song *song.Song) (int, error)
	EditSong(song *song.Song, id int)
	DeleteSong(song *song.Song) error
}
