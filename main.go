package main

import (
	"errors"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"plutos-revenge/common"
	"plutos-revenge/core"
)

func main() {
	game := core.NewGame()

	ebiten.SetWindowSize(common.ScreenWidth*common.Scale, common.ScreenHeight*common.Scale)
	ebiten.SetWindowTitle("Pluto's Revenge")
	err := ebiten.RunGame(game)
	if err != nil {
		if errors.Is(err, common.NormalEscapeError) {
			log.Println("exiting normally")
		} else {
			log.Fatal(err)
		}
	}
}
