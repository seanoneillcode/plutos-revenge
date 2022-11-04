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
	}

	ctx := audio.NewContext(sampleRate)

	m.sounds["alien-hurt"] = loadSound(ctx, "res/fx/alien-hurt.ogg")
	m.sounds["alien-shoot"] = loadSound(ctx, "res/fx/alien-shoot.ogg")
	m.sounds["blast"] = loadSound(ctx, "res/fx/blast.ogg")
	m.sounds["block"] = loadSound(ctx, "res/fx/block.ogg")
	m.sounds["level-start"] = loadSound(ctx, "res/fx/level-start.ogg")
	m.sounds["mystery-entrance"] = loadSound(ctx, "res/fx/mystery-entrance.ogg")
	m.sounds["pickup"] = loadSound(ctx, "res/fx/pickup.ogg")
	m.sounds["player-death"] = loadSound(ctx, "res/fx/player-death-wash.ogg")
	m.sounds["player-shoot"] = loadSound(ctx, "res/fx/player-shoot.ogg")

	m.sounds["blast"].SetVolume(0.5)

	return m
}

func loadSound(ctx *audio.Context, file string) *audio.Player {

	oggSound, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
	}
	s, err := vorbis.DecodeWithSampleRate(sampleRate, oggSound)

	player, err := ctx.NewPlayer(s)
	if err != nil {
		fmt.Fprint(os.Stderr, "failed to create player "+err.Error())
		return nil
	}
	return player
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
