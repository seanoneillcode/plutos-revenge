package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"plutos-revenge/common"
)

type Effect struct {
	animation *Animation
	x         float64
	y         float64
}

func (r *Effect) Update(delta float64) {
	r.animation.Update(delta)
}

func (r *Effect) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(r.x, r.y)
	op.GeoM.Scale(common.Scale, common.Scale)
	screen.DrawImage(r.animation.GetCurrentFrame(), op)
}
