package playlist

import (
	"errors"
	"github.com/v1shn3vsk7/PlaylistAPI/pkg/song"
	"log"
	"sync"
	"time"
)

type Node struct {
	Song *song.Song
	Next *Node
	Prev *Node
}

type Playlist struct {
	mu         sync.Mutex
	currSong   *Node
	front      *Node
	back       *Node
	IsPlaying  bool
	timePlayed time.Duration //change to time.Time
	startTime  time.Time
}

func New() *Playlist {
	return &Playlist{}
}

func (p *Playlist) Play() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.IsPlaying {
		return errors.New("already playing")
	}
	if p.front == nil {
		return errors.New("playlist is empty")
	}

	if p.timePlayed == 0 {
		log.Printf("Playing song: %s by %s, duration: %s\n",
			p.currSong.Song.Name,
			p.currSong.Song.Artist,
			p.currSong.Song.Duration.String()) //Minutes()
	} else {
		log.Printf("Resuming playback of song: %s by %s, duration: %s\n",
			p.currSong.Song.Name,
			p.currSong.Song.Artist,
			p.currSong.Song.Duration - p.timePlayed)
	}

	p.IsPlaying = true
	p.startTime = time.Now()

	go func() {
		time.Sleep(p.currSong.Song.Duration - p.timePlayed)

		p.timePlayed = 0

		p.Next()
	}()

	return nil
}

func (p *Playlist) Pause() error {
	if !p.IsPlaying {
		return errors.New("already paused")
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	log.Printf("paused playback\n")

	p.timePlayed += time.Since(p.startTime)
	p.IsPlaying = false

	return nil
}

func (p *Playlist) AddSong(song *song.Song) {
	p.mu.Lock()
	defer p.mu.Unlock()

	el := &Node {Song: song}

	if p.front == nil {
		p.front, p.back, p.currSong = el, el, el
	} else {
		p.back.Next, el.Prev, p.back = el, p.back, el
	}
}

func (p *Playlist) Next() error {
	if p.front == nil {
		return errors.New("playlist is empty")
	}
	if err := p.Pause(); err != nil {
		return err
	}

	p.mu.Lock()

	p.timePlayed = 0

	if p.currSong.Next != nil {
		p.currSong = p.currSong.Next
	} else {
		p.currSong = p.front
	}

	p.mu.Unlock()

	if err := p.Play(); err != nil {
		return err
	}

	return nil
}

func (p *Playlist) Prev() error {
	if p.back == nil {
		return errors.New("playlist is empty")
	}
	if err := p.Pause(); err != nil {
		return err
	}

	p.mu.Lock()

	p.timePlayed = 0

	if p.currSong.Prev != nil {
		p.currSong = p.currSong.Prev
	} else {
		p.currSong = p.back
	}

	p.mu.Unlock()

	if err := p.Play(); err != nil {
		return err
	}

	return nil
}

func (p *Playlist) Delete(song *song.Song) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.front == p.back {
		p.front, p.back, p.currSong = nil, nil, nil
	}

	node := p.front

	for node != nil {
		if node.Song.Name == song.Name &&
			node.Song.Artist == song.Artist {
				if node.Next == nil {
					node.Next = p.front
				} else {
					node.Prev.Next = node.Next
				}
				if node.Prev == nil {
					node.Prev = p.back
				} else {
					node.Next.Prev = node.Prev
				}

				return
		}
		node = node.Next
	}
}

func (p *Playlist) GetCurrentSong() *song.Song {
	return p.currSong.Song
}