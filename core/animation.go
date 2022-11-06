package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

type Animation struct {
	image           *ebiten.Image
	numFrames       int
	size            int
	frameTimeAmount float64
	isLoop          bool
	// state
	frame  int
	timer  float64
	isDone bool
}

func (r *Animation) Update(delta float64) {
	if r.isDone {
		return
	}
	r.timer = r.timer + delta
	if r.timer > r.frameTimeAmount {
		r.timer = r.timer - r.frameTimeAmount
		r.frame = r.frame + 1
		if r.frame == r.numFrames {
			if r.isLoop {
				r.frame = 0
			} else {
				r.frame = r.numFrames - 1
				r.isDone = true
			}
		}
	}
}

func (r *Animation) Play() {
	r.timer = 0
	r.frame = 0
	r.isDone = false
}

func (r *Animation) GetCurrentFrame() *ebiten.Image {
	return r.image.SubImage(image.Rect(r.frame*r.size, 0, (r.frame+1)*r.size, r.size)).(*ebiten.Image)
}
