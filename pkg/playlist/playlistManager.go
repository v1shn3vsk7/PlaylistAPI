package playlist

import "github.com/v1shn3vsk7/PlaylistAPI/pkg/song"

type PlaylistManager interface {
	New()                     *Playlist
	Play()                    error
	Pause()                   error
	AddSong(song *song.Song)
	Next()                    error
	Prev()                    error
}
