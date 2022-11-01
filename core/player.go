package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
	"math"
	"plutos-revenge/common"
)

const playingState = "playing"
const dyingState = "dying"

const dyingTimeAmount = 0.8
const shootTimerAmount = 0.5
const moveFrameAmount = 0.04
const playerYNormal = 200.0

type Player struct {
	x           float64
	y           float64
	targetY     float64
	image       string
	frame       int
	frames      int
	size        int
	state       string
	speed       float64
	timer       float64
	shootTimer  float64
	images      map[string]*ebiten.Image
	lives       int
	animTimer   float64
	targetFrame int
	moveYSpeed  float64
}

func NewPlayer() *Player {
	p := &Player{
		image:      "player",
		state:      playingState,
		y:          common.ScreenHeight,
		x:          common.ScreenWidth / 2,
		targetY:    playerYNormal,
		moveYSpeed: 40,
		speed:      80,
		size:       12,
		lives:      1,
		shootTimer: -1,
		images: map[string]*ebiten.Image{
			"player": common.LoadImage("player.png"),
		},
		animTimer: -1,
		frame:     3,
	}
	return p
}

func (r *Player) Update(delta float64, game *Game) {
	switch r.state {
	case playingState:

		inputX := 0
		r.targetFrame = 3
		if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
			inputX = -1
			r.targetFrame = 1
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
			inputX = 1
			r.targetFrame = 5
		}
		if r.animTimer >= 0 {
			r.animTimer = r.animTimer - delta
		}
		if r.frame < r.targetFrame {
			if r.animTimer < 0 {
				r.frame = r.frame + 1
				r.animTimer = moveFrameAmount
			}
		}
		if r.frame > r.targetFrame {
			if r.animTimer < 0 {
				r.frame = r.frame - 1
				r.animTimer = moveFrameAmount
			}
		}
		if r.shootTimer > 0 {
			r.shootTimer = r.shootTimer - delta
		}
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			if r.shootTimer < 0 {
				game.AddBullet(r.x, r.y-float64(bulletSize+1), -1, "player")
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
			game.GameOver()
		}
	}
	if math.Abs(r.y-r.targetY) < (2 * r.moveYSpeed * delta) {
		r.y = r.targetY
	}
	if r.y < r.targetY {
		r.y = r.y + r.moveYSpeed*delta
	}
	if r.y > r.targetY {
		r.y = r.y - r.moveYSpeed*delta
	}
}

func (r *Player) Draw(screen *ebiten.Image) {
	if r.state == dyingState {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(r.x, r.y)
	op.GeoM.Scale(common.Scale, common.Scale)
	screen.DrawImage(r.images[r.image].SubImage(image.Rect(r.frame*r.size, 0, (r.frame+1)*r.size, r.size)).(*ebiten.Image), op)
}

func (r *Player) GetHit(game *Game) {
	r.lives = r.lives - 1
	if r.lives < 0 {
		r.state = dyingState
		r.timer = dyingTimeAmount
		game.AddEffect(r.x, r.y, "player-death")
	}
}

func (r *Player) Target(target float64) {
	r.targetY = target
}
