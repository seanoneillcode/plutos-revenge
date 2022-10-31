package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"plutos-revenge/common"
)

const playingState = "playing"
const dyingState = "dying"

const dyingTimeAmount = 1.5
const shootTimerAmount = 0.5

type Player struct {
	x          float64
	y          float64
	image      string
	frame      int
	size       int
	state      string
	speed      float64
	timer      float64
	shootTimer float64
	images     map[string]*ebiten.Image
	lives      int
}

func NewPlayer() *Player {
	return &Player{
		image:      "playing",
		state:      playingState,
		y:          200,
		speed:      60,
		size:       12,
		lives:      3,
		shootTimer: -1,
		images: map[string]*ebiten.Image{
			"playing":   common.LoadImage("player.png"),
			"explosion": common.LoadImage("explosion.png"),
		},
	}
}

func (r *Player) Update(delta float64, game *Game) {
	switch r.state {
	case playingState:
		inputX := 0
		if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
			inputX = -1
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
			inputX = 1
		}
		if r.shootTimer > 0 {
			r.shootTimer = r.shootTimer - delta
		}
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			if r.shootTimer < 0 {
				game.AddBullet(r.x, r.y-float64(r.size), -1)
				r.shootTimer = shootTimerAmount
			}
		}
		r.x = r.x + (float64(inputX) * delta * r.speed)
		if r.x < 0 {
			r.x = 0
		}
		if r.x > common.ScreenWidth-float64(r.size) {
			r.x = common.ScreenWidth - float64(r.size)
		}
	case dyingState:
		r.timer = r.timer - delta
		if r.timer < 0 {
			game.PlayerDeath()
		}
	}
}

func (r *Player) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(r.x, r.y)
	op.GeoM.Scale(common.Scale, common.Scale)
	screen.DrawImage(r.images[r.image].SubImage(image.Rect(r.frame*r.size, 0, (r.frame+1)*r.size, r.size)).(*ebiten.Image), op)
}

func (r *Player) GetHit() {
	r.lives = r.lives - 1
	if r.lives < 0 {
		r.state = dyingState
		r.timer = dyingTimeAmount
		r.image = "explosion"
		r.frame = 0
	}
}
