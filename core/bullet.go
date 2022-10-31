package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"plutos-revenge/common"
)

const normalState = "normal"
const hitState = "hit"

const bulletSize = 12

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
	kind   string
}

func NewBullet(x float64, y float64, dir int, kind string) *Bullet {
	return &Bullet{
		x:     x,
		y:     y,
		dir:   dir,
		image: "normal",
		images: map[string]*ebiten.Image{
			"normal": common.LoadImage("bullet.png"),
		},
		state: normalState,
		size:  bulletSize,
		speed: 60,
		kind:  kind,
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
		switch r.kind {
		case "alien":
			// check if hit player
			if game.player.state == playingState {
				if common.Overlap(game.player.x, game.player.y, float64(game.player.size), r.x, r.y, float64(r.size)) {
					r.state = hitState
					game.player.GetHit(game)
					game.AddEffect(r.x, r.y, "explosion")
					r.GetHit()
				}
			}
		case "player":
			// check if hit alien
			for index, a := range game.aliens {
				if a.state == normalAlienState {
					if common.Overlap(a.x, a.y, float64(a.size), r.x, r.y, float64(r.size)) {
						r.state = hitState
						game.aliens[index].GetHit()
						game.AddEffect(game.aliens[index].x, game.aliens[index].y, "explosion")
						r.GetHit()
						game.AddEffect(r.x, r.y, "explosion")
						game.ScorePoint()
					}
				}
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
}

func (r *Bullet) Draw(screen *ebiten.Image) {
	if r.state == normalState {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(r.x, r.y)
		op.GeoM.Scale(common.Scale, common.Scale)
		screen.DrawImage(r.images[r.image].SubImage(image.Rect(r.frame*r.size, 0, (r.frame+1)*r.size, r.size)).(*ebiten.Image), op)
	}
}
