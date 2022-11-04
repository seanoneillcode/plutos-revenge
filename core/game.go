package core

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"plutos-revenge/common"
	"time"
)

const menuGameState = "menu"
const playingGameState = "playing"
const gameOverGameState = "gameOver"
const levelOverGameState = "levelOver"
const gameWonGameState = "gameWon"

const numberOfLevels = 4

type Game struct {
	lastUpdateCalled time.Time
	player           *Player
	bullets          []*Bullet
	aliens           []*Alien
	blocks           []*Block
	alienGroup       *AlienGroup
	stars            []*Star
	effects          []*Animation
	earth            *Earth
	fader            *Fader
	mystery          *Mystery
	timer            float64
	state            string
	images           map[string]*ebiten.Image
	level            int
	score            int
	soundManager     *common.SoundManager
}

func NewGame() *Game {
	g := &Game{
		state: menuGameState,
		images: map[string]*ebiten.Image{
			"splash":       common.LoadImage("splash.png"),
			"explosion":    common.LoadImage("explosion.png"),
			"player-death": common.LoadImage("player-death.png"),
			"alien-death":  common.LoadImage("pluton-death.png"),
			"plus-one":     common.LoadImage("plus-one.png"),
		},
		stars:            []*Star{},
		fader:            NewFader(),
		earth:            NewEarth(),
		lastUpdateCalled: time.Now(),
		mystery:          NewMystery(common.ScreenWidth, 8, 1),
		soundManager:     common.NewManager(),
	}
	for index := 0; index < 100; index += 1 {
		g.stars = append(g.stars, NewStar())
	}
	return g
}

func (r *Game) Update() error {
	delta := float64(time.Now().Sub(r.lastUpdateCalled).Milliseconds()) / 1000
	r.lastUpdateCalled = time.Now()

	for _, s := range r.stars {
		s.Update(delta)
	}

	switch r.state {
	case menuGameState:
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			r.StartNewGame()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			return common.NormalEscapeError
		}
	case playingGameState:
		r.mystery.Update(delta, r)
		r.player.Update(delta, r)
		for _, b := range r.bullets {
			b.Update(delta, r)
		}
		r.alienGroup.Update(delta, r)
		for _, a := range r.aliens {
			a.Update(delta, r)
		}
		for _, b := range r.blocks {
			b.Update(delta, r)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			r.QuitToMenu()
		}
	case gameOverGameState:
		for _, b := range r.bullets {
			b.Update(delta, r)
		}
		r.mystery.Update(delta, r)
		r.alienGroup.Update(delta, r)
		for _, a := range r.aliens {
			a.Update(delta, r)
		}
		r.timer = r.timer - delta
		if r.timer < 0 {
			r.timer = 0
			r.QuitToMenu()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			r.QuitToMenu()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			r.QuitToMenu()
		}
	case levelOverGameState:
		r.timer = r.timer - delta
		if r.timer < 0 {
			r.timer = 0
			r.StartNewLevel()
		}
		r.mystery.Update(delta, r)
		r.player.Update(delta, r)
		for _, b := range r.bullets {
			b.Update(delta, r)
		}
	case gameWonGameState:
		r.player.Update(delta, r)
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			r.QuitToMenu()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			r.QuitToMenu()
		}
	}

	r.earth.Update(delta)
	for _, e := range r.effects {
		e.Update(delta)
		if e.done {
			r.RemoveAnimation(e)
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	return nil
}

func (r *Game) Draw(screen *ebiten.Image) {

	for _, s := range r.stars {
		s.Draw(screen)
	}
	switch r.state {
	case menuGameState:
		r.drawImage(screen, "splash", 40, 40)
		common.DrawText(screen, "start game", 60, 120)
	case playingGameState:
		r.earth.Draw(screen)
		r.player.Draw(screen)
		for _, b := range r.bullets {
			b.Draw(screen)
		}
		r.mystery.Draw(screen)
		for _, a := range r.aliens {
			a.Draw(screen)
		}
		for _, b := range r.blocks {
			b.Draw(screen)
		}
		for _, e := range r.effects {
			e.Draw(screen)
		}
		common.DrawText(screen, fmt.Sprintf("score %d", r.score), 4, 4)
		common.DrawText(screen, fmt.Sprintf("lives %d", r.player.lives), 130, 4)
	case gameOverGameState:
		r.earth.Draw(screen)
		r.player.Draw(screen)
		for _, b := range r.bullets {
			b.Draw(screen)
		}
		r.mystery.Draw(screen)
		for _, a := range r.aliens {
			a.Draw(screen)
		}
		for _, b := range r.blocks {
			b.Draw(screen)
		}
		for _, e := range r.effects {
			e.Draw(screen)
		}
		common.DrawText(screen, "game over", 60, 90)
	case levelOverGameState:
		r.earth.Draw(screen)
		r.player.Draw(screen)
		r.mystery.Draw(screen)
		for _, b := range r.bullets {
			b.Draw(screen)
		}
		for _, b := range r.blocks {
			b.Draw(screen)
		}
		for _, e := range r.effects {
			e.Draw(screen)
		}
		common.DrawText(screen, "   wave\ndestroyed", 60, 90)
	case gameWonGameState:
		r.earth.Draw(screen)
		r.player.Draw(screen)
		for _, e := range r.effects {
			e.Draw(screen)
		}
		common.DrawText(screen, " you win!", 60, 70)
		common.DrawText(screen, "the earth\nis saved!", 60, 120)
	}

}

func (r *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return common.ScreenWidth * common.Scale, common.ScreenHeight * common.Scale
}

func (r *Game) drawImage(screen *ebiten.Image, img string, x float64, y float64) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	op.GeoM.Scale(common.Scale, common.Scale)
	screen.DrawImage(r.images[img], op)
}
