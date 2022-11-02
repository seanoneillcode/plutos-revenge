package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"math"
	"plutos-revenge/common"
)

type Fader struct {
	fadeColor   color.RGBA
	fadeTargetA uint8
	image       *ebiten.Image
}

func NewFader() *Fader {
	return &Fader{
		fadeColor: color.RGBA{
			R: 255,
			G: 255,
			B: 255,
			A: 255,
		},
		fadeTargetA: 0,
		image:       common.LoadImage("fade.png"),
	}
}

func (r *Fader) Update() error {
	var fadeSpeed uint8 = 8
	if math.Abs(float64(r.fadeColor.A-r.fadeTargetA)) < float64(fadeSpeed) {
		r.fadeColor.A = r.fadeTargetA
	}
	if r.fadeColor.A < r.fadeTargetA {
		r.fadeColor.A = r.fadeColor.A + fadeSpeed
	}
	if r.fadeColor.A > r.fadeTargetA {
		r.fadeColor.A = r.fadeColor.A - fadeSpeed
	}
	return nil
}

func (r *Fader) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(common.ScreenWidth*common.Scale, common.ScreenHeight*common.Scale)
	op.ColorM.ScaleWithColor(r.fadeColor)
	screen.DrawImage(r.image, op)
}
