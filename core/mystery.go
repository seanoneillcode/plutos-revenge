package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"math/rand"
	"plutos-revenge/common"
)

type Mystery struct {
	x      float64
	y      float64
	image  *ebiten.Image
	dir    int
	size   float64
	timer  float64
	active bool
}

func NewMystery(x float64, y float64, dir int) *Mystery {
	return &Mystery{
		x:     x,
		y:     y,
		dir:   dir,
		image: common.LoadImage("mystery.png"),
		size:  alienSize,
	}
}

func (r *Mystery) Update(delta float64) {

	if r.active {
		r.x = r.x + (30 * delta * float64(r.dir))
		// gone from the side of the screen
		if r.x > common.ScreenWidth {
			r.reset()
		}
		if r.x < 0-r.size {
			r.reset()
		}
	} else {
		r.timer = r.timer - delta
		if r.timer < 0 {
			r.active = true
			r.x = common.ScreenWidth
			r.dir = -1
			if rand.Float64() > 0.5 {
				r.x = 0 - r.size
				r.dir = 1
			}
		}
	}

}

func (r *Mystery) GetHit(game Game) {
	r.reset()
	game.AddEffect(r.x, r.y, "explosion")
	game.player.lives += 1
}

func (r *Mystery) reset() {
	r.active = false
	r.x = common.ScreenWidth
	r.timer = rand.Float64()*20 + 20
}

func (r *Mystery) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(r.x, r.y)
	op.GeoM.Scale(common.Scale, common.Scale)
	screen.DrawImage(r.image, op)
}
