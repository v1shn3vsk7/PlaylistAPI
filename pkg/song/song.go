package song

import "time"

type Song struct {
	Name     string
	Artist   string
	Duration time.Duration
}

func NewSong(name, author string, duration time.Duration) *Song {
	return &Song {
		Name:     name,
		Artist:   author,
		Duration: duration,
	}
}
