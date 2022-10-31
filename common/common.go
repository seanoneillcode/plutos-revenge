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
	b, err := ioutil.ReadFile("res/" + imageFileName)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		log.Fatal(err)
	}
	return ebiten.NewImageFromImage(img)
}
