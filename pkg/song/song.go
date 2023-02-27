package song

import "time"

type Song struct {
	Name     string
	Author   string
	Duration time.Duration
}

func NewSong(name, author string, duration time.Duration) *Song {
	return &Song {
		Name:     name,
		Author:   author,
		Duration: duration,
	}
}
