package main

import (
	"errors"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"log"
	"time"
)

const (
	ScreenWidth  = 240
	ScreenHeight = 180
	Scale        = 4
)

var NormalEscapeError = errors.New("normal escape termination")

func main() {
	game := &Game{}

	ebiten.SetWindowSize(ScreenWidth*Scale, ScreenHeight*Scale)
	ebiten.SetWindowTitle("Pluto's Revenge")
	err := ebiten.RunGame(game)
	if err != nil {
		if errors.Is(err, NormalEscapeError) {
			log.Println("exiting normally")
		} else {
			log.Fatal(err)
		}
	}
}

type Game struct {
	lastUpdateCalled time.Time
}

func (g *Game) Update() error {
	//_ := time.Now().Sub(g.lastUpdateCalled).Milliseconds()
	g.lastUpdateCalled = time.Now()

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return NormalEscapeError
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth * Scale, ScreenHeight * Scale
}
