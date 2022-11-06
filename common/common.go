package common

import (
	"bytes"
	"errors"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"io/ioutil"
	"log"

	_ "image/png" // evil, required for decoder to 'know' what a png is...
)

const (
	ScreenWidth  = 180
	ScreenHeight = 240
	Scale        = 4
)

var NormalEscapeError = errors.New("normal escape termination")

func LoadImage(imageFileName string) *ebiten.Image {
	return loadImage("res/" + imageFileName)
}

func loadImage(imageFileName string) *ebiten.Image {
	b, err := ioutil.ReadFile(imageFileName)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		log.Fatal(err)
	}
	return ebiten.NewImageFromImage(img)
}

func Overlap(x1, y1, s1, x2, y2, s2 float64) bool {
	if x2 > x1+s1 || x2+s2 < x1 {
		return false
	}
	if y2 > y1+s1 || y2+s2 < y1 {
		return false
	}
	return true
}
