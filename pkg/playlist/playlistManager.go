package playlist

type PlaylistManager interface {
	Init()
	Play()    error
	Pause()   error
	AddSong()
	Next()    error
	Prev()    error
}
