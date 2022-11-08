package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"math/rand"
	"plutos-revenge/common"
)

const alienSize = 12

type Alien struct {
	x         float64
	y         float64
	dir       int
	animation *Animation
	size      int
	speed     float64
	timer     float64
	kind      string
	lives     int
}

func NewAlien(x float64, y float64, dir int, kind string, image *ebiten.Image) *Alien {
	a := &Alien{
		x:     x,
		y:     y,
		dir:   dir,
		size:  alienSize,
		speed: 20,
		kind:  kind,
	}
	switch kind {
	case "normal":
		a.animation = &Animation{
			image:           image,
			numFrames:       4,
			size:            alienSize,
			frameTimeAmount: 0.25,
			isLoop:          true,
		}
	case "tough":
		a.animation = &Animation{
			image:           image,
			numFrames:       4,
			size:            alienSize,
			frameTimeAmount: 0.25,
			isLoop:          true,
		}
		a.lives = 1
	case "poison":
		a.animation = &Animation{
			image:           image,
			numFrames:       4,
			size:            alienSize,
			frameTimeAmount: 0.25,
			isLoop:          true,
		}
	case "bomb":
		a.animation = &Animation{
			image:           image,
			numFrames:       4,
			size:            alienSize,
			frameTimeAmount: 0.25,
			isLoop:          true,
		}
	}
	return a
}

func (r *Alien) Update(delta float64, game *Game) {
	// gone from the bottom of the screen
	if r.y > common.ScreenHeight+float64(r.size) {
		game.RemoveAlien(r)
	}
	if r.y > common.ScreenHeight-float64(r.size) {
		game.GameOver()
	}
	// check if hit player
	if game.player.state == playingState {
		if common.Overlap(game.player.x, game.player.y, float64(game.player.size), r.x, r.y, float64(r.size)) {
			game.player.GetHit(game)
			r.GetHit(game)
		}
	}
	r.animation.Update(delta)
}

func (r *Alien) GetHit(game *Game) {
	if r.kind == "tough" {
		if r.lives == 1 {
			r.lives = 0
			r.animation.image = game.images["hurt-tough"]
			game.PlaySound("block")
			return
		}
	}
	game.ScorePoint()
	game.RemoveAlien(r)
	game.PlaySound("alien-hurt")
	switch r.kind {
	case "poison":
		offset := 64.0
		halfOffset := offset / 2.0
		for index := 0; index < 4; index += 1 {
			game.AddEffect(r.x-halfOffset+(rand.Float64()*offset), r.y-halfOffset+(rand.Float64()*offset), "gas")
		}
		game.alienGroup.SpeedUp()
		game.alienGroup.SpeedUp()
		game.alienGroup.SpeedUp()
		game.PlaySound("gas")
	case "bomb":
		game.RemoveAdjacentAliens(r)
		game.PlaySound("bomb")
	default:
		game.AddEffect(r.x, r.y, "alien-death")
	}

}

func (r *Alien) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(r.x, r.y)
	op.GeoM.Scale(common.Scale, common.Scale)
	screen.DrawImage(r.animation.GetCurrentFrame(), op)
}
