package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"plutos-revenge/common"
	"time"
)

type Game struct {
	lastUpdateCalled time.Time
	player           *Player
}

func NewGame() *Game {
	return &Game{
		player: NewPlayer("player.png"),
	}
}

func (r *Game) Update() error {
	delta := float64(time.Now().Sub(r.lastUpdateCalled).Milliseconds()) / 1000
	r.lastUpdateCalled = time.Now()

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return common.NormalEscapeError
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}
	r.player.Update(delta)

	return nil
}

func (r *Game) Draw(screen *ebiten.Image) {
	r.player.Draw(screen)
}

func (r *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return common.ScreenWidth * common.Scale, common.ScreenHeight * common.Scale
}
