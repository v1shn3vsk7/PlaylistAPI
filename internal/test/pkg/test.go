package test

import (
	"github.com/v1shn3vsk7/PlaylistAPI/pkg/playlist"
	"github.com/v1shn3vsk7/PlaylistAPI/pkg/song"
	"testing"
	"time"
)

func TestPlaylist(t *testing.T) {
	p := playlist.New()

	s1 := &song.Song{Name: "Song 1", Artist: "Artist 1", Duration: 5 * time.Second}
	s2 := &song.Song{Name: "Song 2", Artist: "Artist 2", Duration: 10 * time.Second}
	s3 := &song.Song{Name: "Song 3", Artist: "Artist 3", Duration: 15 * time.Second}

	p.AddSong(s1)
	p.AddSong(s2)
	p.AddSong(s3)

	if err := p.Play(); err != nil {
		t.Errorf("unexpected error while playing playlist: %v", err)
	}

	if err := p.Pause(); err != nil {
		t.Errorf("unexpected error while pausing playlist: %v", err)
	}

	if err := p.Play(); err != nil {
		t.Errorf("unexpected error while resuming playlist: %v", err)
	}

	if err := p.Next(); err != nil {
		t.Errorf("unexpected error while skipping to next song: %v", err)
	}

	if err := p.Prev(); err != nil {
		t.Errorf("unexpected error while skipping to previous song: %v", err)
	}

	if err := p.Prev(); err != nil {
		t.Errorf("unexpected error while skipping to previous song: %v", err)
	}

	if err := p.Prev(); err == nil {
		t.Error("expected error while skipping past the first song")
	}
}