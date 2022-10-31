package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"plutos-revenge/common"
)

const normalAlienState = "normal"
const hitAlienState = "hit"

const alienSize = 12

type Alien struct {
	x      float64
	y      float64
	dir    int
	image  string
	images map[string]*ebiten.Image
	frame  int
	size   int
	state  string
	speed  float64
	timer  float64
}

func NewAlien(x float64, y float64, dir int) *Alien {
	return &Alien{
		x:     x,
		y:     y,
		dir:   dir,
		image: "normal",
		images: map[string]*ebiten.Image{
			"normal": common.LoadImage("pluton.png"),
		},
		state: normalAlienState,
		size:  alienSize,
		speed: 20,
	}
}

func (r *Alien) Update(delta float64, game *Game) {
	switch r.state {
	case normalState:
		// gone from the bottom of the screen
		if r.y > common.ScreenHeight+float64(r.size) {
			game.RemoveAlien(r)
		}
		// check if hit player
		if game.player.state == playingState {
			if common.Overlap(game.player.x, game.player.y, float64(game.player.size), r.x, r.y, float64(r.size)) {
				r.state = hitState
				game.player.GetHit(game)
				game.AddEffect(r.x, r.y, "explosion")
				r.GetHit()
			}
		}

	case hitState:
		r.timer = r.timer - delta
		if r.timer < 0 {
			game.RemoveAlien(r)
		}
	}
}

func (r *Alien) GetHit() {
	r.state = hitAlienState
	r.timer = 1.0
}

func (r *Alien) Draw(screen *ebiten.Image) {
	if r.state == normalAlienState {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(r.x, r.y)
		op.GeoM.Scale(common.Scale, common.Scale)
		screen.DrawImage(r.images[r.image].SubImage(image.Rect(r.frame*r.size, 0, (r.frame+1)*r.size, r.size)).(*ebiten.Image), op)
	}
}
