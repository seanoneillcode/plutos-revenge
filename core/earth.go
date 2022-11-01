package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"math"
	"plutos-revenge/common"
)

type Earth struct {
	y          float64
	targetY    float64
	image      *ebiten.Image
	moveAmount float64
}

func NewEarth() *Earth {
	return &Earth{
		y:          common.ScreenHeight,
		targetY:    0,
		image:      common.LoadImage("earth.png"),
		moveAmount: 6,
	}
}

func (r *Earth) Update(delta float64) {
	if math.Abs(r.y-r.targetY) < (2 * r.moveAmount * delta) {
		r.y = r.targetY
	}
	if r.y < r.targetY {
		r.y = r.y + r.moveAmount*delta
	}
	if r.y > r.targetY {
		r.y = r.y - r.moveAmount*delta
	}
}

func (r *Earth) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, r.y)
	op.GeoM.Scale(common.Scale, common.Scale)
	screen.DrawImage(r.image, op)
}

func (r *Earth) Target(target float64) {
	r.targetY = target
}
