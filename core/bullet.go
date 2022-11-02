package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"plutos-revenge/common"
)

const normalState = "normal"
const hitState = "hit"

const bulletSize = 4

type Bullet struct {
	x         float64
	y         float64
	dir       int
	image     string
	images    map[string]*ebiten.Image
	imageSize int
	frame     int
	size      int
	offsetX   float64
	offsetY   float64
	state     string
	speed     float64
	timer     float64
	kind      string
}

func NewBullet(x float64, y float64, dir int, kind string) *Bullet {
	b := &Bullet{
		x:     x,
		y:     y,
		dir:   dir,
		image: "bullet",
		images: map[string]*ebiten.Image{
			"bullet": common.LoadImage("bullet.png"),
			"lazer":  common.LoadImage("lazer.png"),
		},
		state:     normalState,
		size:      bulletSize,
		imageSize: 8,
		offsetX:   -2,
		offsetY:   -4,
		speed:     120,
		kind:      kind,
	}
	if kind == "alien" {
		b.image = "lazer"
	}
	return b
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
		for _, b := range game.bullets {
			if b.state == normalState && b != r {
				if common.Overlap(b.x, b.y, float64(b.size-1), r.x, r.y, float64(r.size)) {
					r.state = hitState
					b.GetHit()
					game.AddEffect(b.x, b.y, "explosion")
					r.GetHit()
					game.AddEffect(r.x, r.y, "explosion")
				}
			}
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
			for _, a := range game.blocks {
				if a.lives > 0 {
					if common.Overlap(a.x, a.y, float64(a.size), r.x, r.y, float64(r.size)) {
						a.GetHit(game)
						r.state = hitState
						r.GetHit()
						game.AddEffect(r.x, r.y, "explosion")
					}
				}
			}
		case "player":
			// check if hit alien
			for _, a := range game.aliens {
				if a.state == normalAlienState {
					if common.Overlap(a.x, a.y, float64(a.size), r.x, r.y, float64(r.size)) {
						r.state = hitState
						a.GetHit()
						game.AddEffect(a.x, a.y, "alien-death")
						r.GetHit()
						game.AddEffect(r.x, r.y, "explosion")
						game.ScorePoint()
					}
				}
			}
			for _, a := range game.blocks {
				if a.lives > 0 {
					if common.Overlap(a.x, a.y, float64(a.size-1), r.x, r.y, float64(r.size-1)) {
						r.state = hitState
						r.GetHit()
						game.AddEffect(r.x, r.y, "explosion")
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
		op.GeoM.Translate(r.x+r.offsetX, r.y+r.offsetY)
		op.GeoM.Scale(common.Scale, common.Scale)
		screen.DrawImage(r.images[r.image].SubImage(image.Rect(r.frame*r.imageSize, 0, (r.frame+1)*r.imageSize, r.imageSize)).(*ebiten.Image), op)
	}
}
