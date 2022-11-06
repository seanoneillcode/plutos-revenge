package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"math/rand"
	"plutos-revenge/common"
)

const mysteryTimerAmount = 10 // mysteryTimerAmount

type Mystery struct {
	x         float64
	y         float64
	animation *Animation
	dir       int
	size      float64
	timer     float64
	active    bool
	speed     float64
	frame     int
}

func NewMystery(x float64, y float64, dir int) *Mystery {
	m := &Mystery{
		x:   x,
		y:   y,
		dir: dir,
		animation: &Animation{
			image:           common.LoadImage("mystery.png"),
			numFrames:       2,
			size:            12,
			frameTimeAmount: 0.5,
			isLoop:          true,
		},
		size:  alienSize,
		speed: 40,
	}
	m.reset()
	return m
}

func (r *Mystery) Update(delta float64, game *Game) {

	if r.active {
		r.x = r.x + (r.speed * delta * float64(r.dir))
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
			game.PlaySound("mystery-entrance")
			r.active = true
			r.x = common.ScreenWidth
			r.dir = -1
			if rand.Float64() > 0.5 {
				r.x = 0 - r.size
				r.dir = 1
			}
		}
	}
	r.animation.Update(delta)
}

func (r *Mystery) GetHit(game *Game) {
	game.player.lives += 1
	r.reset()
	game.PlaySound("pickup")
}

func (r *Mystery) reset() {
	r.active = false
	r.x = common.ScreenWidth
	r.timer = rand.Float64()*mysteryTimerAmount + mysteryTimerAmount
}

func (r *Mystery) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(r.x, r.y)
	op.GeoM.Scale(common.Scale, common.Scale)
	screen.DrawImage(r.animation.GetCurrentFrame(), op)
}
