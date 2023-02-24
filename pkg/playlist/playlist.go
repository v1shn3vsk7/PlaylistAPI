package playlist

import (
	"container/list"
	"errors"
	"github.com/v1shn3vsk7/PlaylistAPI/pkg/song"
	"log"
	"sync"
	"time"
)

type Playlist struct {
	Songs     *list.List
	CurrSong  *song.Song
	NextSong  *list.Element
	PrevSong  *list.Element
	IsPlaying  bool
	mu         sync.Mutex
	timePlayed time.Duration
}

func Init() *Playlist {
	return &Playlist{
		Songs: list.New(),
	}
}

func (p *Playlist) Play() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.IsPlaying {
		return errors.New("already playing")
	}

	if p.Songs.Len() == 0 {
		return errors.New("playlist is empty")
	}

	p.IsPlaying = true
	if p.CurrSong == nil {
		front := p.Songs.Front()
		p.CurrSong = front.Value.(*song.Song)

		p.PrevSong = front.Prev()
		p.NextSong = front.Next()
	}

	if p.timePlayed == 0 {
		log.Printf("Playing song: %s by %s, Duration: %s\n",
			p.CurrSong.Name,
			p.CurrSong.Author,
			p.CurrSong.Duration.String())
	} else {
		log.Printf("Resuming playback of song: %s by %s, Duration: %s\n",
			p.CurrSong.Name,
			p.CurrSong.Author,
			p.CurrSong.Duration-p.timePlayed)
	}

	go func() {
		time.Sleep(p.CurrSong.Duration - p.timePlayed)

		p.mu.Lock()
		defer p.mu.Unlock()

		p.IsPlaying = false
		p.timePlayed = 0

		if p.NextSong != nil {
			p.CurrSong = p.NextSong.Value.(*song.Song)
			p.NextSong = p.NextSong.Next()
			go p.Play()
		} else {
			log.Println("End of playlist")
		}
	}()

	return nil
}

func (p *Playlist) Pause() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.IsPlaying {
		return errors.New("already paused")
	}

	p.IsPlaying = false

	log.Printf("Paused playback of song: %s by %s, Duration: %s\n",
		p.CurrSong.Name,
		p.CurrSong.Author,
		p.CurrSong.Duration-p.timePlayed)

	remainingTime := p.CurrSong.Duration - p.timePlayed
	p.timePlayed = 0

	go func() {
		time.Sleep(remainingTime)

		p.mu.Lock()
		defer p.mu.Unlock()

		if p.IsPlaying {
			return
		}

		p.IsPlaying = true
		p.timePlayed = 0
		go p.Play()
	}()

	return nil
}

func (p *Playlist) AddSongs(songs ...*song.Song) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, s := range songs {
		p.Songs.PushBack(s)
	}

	if p.CurrSong == nil {
		p.CurrSong = songs[0]
	}
}

func (p *Playlist) Next() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.Songs.Len() == 0 {
		return errors.New("playlist is empty")
	}

	if p.CurrSong == nil {
		p.CurrSong = p.Songs.Front().Value.(*song.Song)
	}

	if p.NextSong == nil {
		p.NextSong = p.Songs.Front()
	}

	if err := p.Pause(); err != nil {
		return err
	}

	p.CurrSong = p.NextSong.Value.(*song.Song)

	go p.Play()

	return nil
}

func (p *Playlist) Prev() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.Songs.Len() == 0 {
		return errors.New("playlist is empty")
	}

	if p.CurrSong == nil {
		p.CurrSong = p.Songs.Back().Value.(*song.Song)
		return nil
	}

	if p.PrevSong == nil {
		p.CurrSong = p.Songs.Back().Value.(*song.Song)
	}

	p.CurrSong = p.PrevSong.Value.(*song.Song)

	go p.Play()

	return nil
}