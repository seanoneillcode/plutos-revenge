package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"plutos-revenge/common"
)

const normalState = "normal"
const hitState = "hit"

type Bullet struct {
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

func NewBullet(x float64, y float64, dir int) *Bullet {
	return &Bullet{
		x:     x,
		y:     y,
		dir:   dir,
		image: "normal",
		images: map[string]*ebiten.Image{
			"normal":    common.LoadImage("bullet.png"),
			"explosion": common.LoadImage("explosion.png"),
		},
		state: normalState,
		size:  12,
		speed: 50,
	}
}

func (r *Bullet) Update(delta float64, game *Game) {
	switch r.state {
	case normalState:
		r.y = r.y + (float64(r.dir) * delta * r.speed)
		// gone from top of screen
		if r.y < (0 - float64(r.size)) {
			game.RemoveBullet(r)
		}
		// gone from the bottom of the screen
		if r.y > common.ScreenHeight+float64(r.size) {
			game.RemoveBullet(r)
		}
		// check if hit player
		if game.player.state == playingState {
			if common.Overlap(game.player.x, game.player.y, float64(game.player.size), r.x, r.y, float64(r.size)) {
				r.state = hitState
				game.player.GetHit()
				r.GetHit()
			}
		}
	case hitState:
		r.timer = r.timer - delta
		if r.timer < 0 {
			game.RemoveBullet(r)
		}
	}
}

func (r *Bullet) GetHit() {
	r.state = hitState
	r.timer = 1.0
	r.image = "explosion"
	r.frame = 0
}

func (r *Bullet) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(r.x, r.y)
	op.GeoM.Scale(common.Scale, common.Scale)
	screen.DrawImage(r.images[r.image].SubImage(image.Rect(r.frame*r.size, 0, (r.frame+1)*r.size, r.size)).(*ebiten.Image), op)
}
