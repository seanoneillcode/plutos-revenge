package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"math/rand"
	"plutos-revenge/common"
)

type Star struct {
	x     float64
	y     float64
	color color.RGBA
	speed float64
}

func NewStar() *Star {
	s := &Star{
		x: rand.Float64() * float64(common.ScreenWidth),
		y: rand.Float64() * float64(common.ScreenHeight),
		color: color.RGBA{
			R: uint8(rand.Intn(255)),
			G: uint8(rand.Intn(255)),
			B: uint8(rand.Intn(255)),
			A: uint8(255),
		},
		speed: 10 + (rand.Float64() * 10),
	}

	return s
}

func (r *Star) Update(delta float64) {
	r.y = r.y + (-1 * delta * r.speed)
	// gone from top of screen
	if r.y < 0 {
		r.Reset()
	}
}

func (r *Star) Reset() {
	r.x = rand.Float64() * float64(common.ScreenWidth)
	r.y = common.ScreenHeight + 2
	r.color = color.RGBA{
		R: uint8(rand.Intn(255)),
		G: uint8(rand.Intn(255)),
		B: uint8(rand.Intn(255)),
		A: uint8(rand.Intn(255)),
	}
	r.speed = 10 + (rand.Float64() * 10)
}

func (r *Star) Draw(screen *ebiten.Image) {
	screen.Set(int(r.x*common.Scale), int(r.y*common.Scale), r.color)
}
