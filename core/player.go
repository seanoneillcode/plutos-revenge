package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"plutos-revenge/common"
)

const playingState = "playing"
const deadState = "dead"

type Player struct {
	x     float64
	y     float64
	image *ebiten.Image
	frame int
	size  int
	state string
	speed float64
}

func NewPlayer(imageFileName string) *Player {
	return &Player{
		image: common.LoadImage(imageFileName),
		state: playingState,
		y:     200,
		speed: 60,
		size:  12,
	}
}

func (r *Player) Update(delta float64) {
	switch r.state {
	case playingState:
		inputX := 0
		if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
			inputX = -1
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
			inputX = 1
		}
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			// fire bullet
		}
		r.x = r.x + (float64(inputX) * delta * r.speed)
		if r.x < 0 {
			r.x = 0
		}
		if r.x > common.ScreenWidth-float64(r.size) {
			r.x = common.ScreenWidth - float64(r.size)
		}
	case deadState:
		// timer
		// animation
	}
}

func (r *Player) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(r.x, r.y)
	op.GeoM.Scale(common.Scale, common.Scale)

	screen.DrawImage(r.image.SubImage(image.Rect(r.frame*r.size, 0, (r.frame+1)*r.size, r.size)).(*ebiten.Image), op)
}
