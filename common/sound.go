package common

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"os"
)

type SoundManager struct {
	ctx    *audio.Context
	sounds map[string]*audio.Player
}

const sampleRate = 44100

func NewManager() *SoundManager {
	m := &SoundManager{
		sounds: map[string]*audio.Player{},
		ctx:    audio.NewContext(sampleRate),
	}

	//m.sounds["blast"].SetVolume(0.5)

	return m
}

func (r *SoundManager) LoadSound(name string, file string) {
	oggSound, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
	}
	s, err := vorbis.DecodeWithSampleRate(sampleRate, oggSound)
	player, err := r.ctx.NewPlayer(s)
	if err != nil {
		fmt.Fprint(os.Stderr, "failed to create player "+err.Error())
		return
	}
	r.sounds[name] = player
}

func (r *SoundManager) PlaySound(name string) {
	p, ok := r.sounds[name]
	if !ok {
		fmt.Fprint(os.Stderr, "failed to play sound, not loaded: "+name)
		return
	}
	if p == nil {
		fmt.Fprint(os.Stderr, "failed to play sound, nil player: "+name)
		return
	}
	p.Rewind()
	p.Play()
}
