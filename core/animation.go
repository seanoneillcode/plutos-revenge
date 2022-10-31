package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"plutos-revenge/common"
)

type Animation struct {
	frame           int
	numFrames       int
	timer           float64
	frameTimeAmount float64
	loop            bool
	image           *ebiten.Image
	size            int
	x               float64
	y               float64
}

func (r *Animation) Update(delta float64) {
	r.timer = r.timer + delta
	if r.timer > r.frameTimeAmount {
		r.timer = r.timer - r.frameTimeAmount
		r.frame = r.frame + 1
		if r.frame == r.numFrames {
			if r.loop {
				r.frame = 0
			} else {
				r.frame = r.numFrames - 1
			}
		}
	}
}

func (r *Animation) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(r.x, r.y)
	op.GeoM.Scale(common.Scale, common.Scale)
	screen.DrawImage(r.image.SubImage(image.Rect(r.frame*r.size, 0, (r.frame+1)*r.size, r.size)).(*ebiten.Image), op)
}
