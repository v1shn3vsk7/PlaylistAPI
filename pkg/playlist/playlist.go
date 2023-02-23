package playlist

import (
	"container/list"
)

type Playlist struct {
	Songs *list.List
}

func (p *Playlist) Init() {
	p.Songs.Init()
}

func (p *Playlist) Play() error {
	return nil
}

func (p *Playlist) Pause() error {
	return nil
}

func (p *Playlist) AddSong() error {
	return nil
}

func (p *Playlist) Next() error {
	return nil
}

func (p *Playlist) Prev() error {
	return nil
}


