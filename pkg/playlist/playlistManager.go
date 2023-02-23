package playlist

type PlaylistManager interface {
	Play()    error
	Pause()   error
	AddSong() error
	Next()    error
	Prev()    error
}
