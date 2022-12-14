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
	effects          []*Effect
	earth            *Earth
	mystery          *Mystery
	timer            float64
	state            string
	images           map[string]*ebiten.Image
	level            int
	score            int
	soundManager     *common.SoundManager
}

func NewGame() *Game {
	r := &Game{
		state: menuGameState,
		images: map[string]*ebiten.Image{
			// load all the images once, up front
			"splash":       common.LoadImage("splash.png"),
			"explosion":    common.LoadImage("explosion.png"),
			"player-death": common.LoadImage("player-death.png"),
			"alien-death":  common.LoadImage("alien-death.png"),
			"plus-one":     common.LoadImage("plus-one.png"),
			"poison":       common.LoadImage("heart.png"),
			"tough":        common.LoadImage("tough.png"),
			"normal":       common.LoadImage("pluton.png"),
			"bomb":         common.LoadImage("bomb.png"),
			"hurt-tough":   common.LoadImage("tough-hurt.png"),
			"gas":          common.LoadImage("gas.png"),
			"block":        common.LoadImage("block.png"),
			"earth":        common.LoadImage("earth.png"),
			"mystery":      common.LoadImage("mystery.png"),
			"player":       common.LoadImage("player.png"),
		},
		lastUpdateCalled: time.Now(),
		soundManager:     common.NewManager(),
	}

	r.soundManager.LoadSound("alien-hurt", "res/fx/alien-hurt.ogg")
	r.soundManager.LoadSound("alien-shoot", "res/fx/alien-shoot.ogg")
	r.soundManager.LoadSound("blast", "res/fx/blast.ogg")
	r.soundManager.LoadSound("block", "res/fx/block.ogg")
	r.soundManager.LoadSound("level-start", "res/fx/level-start.ogg")
	r.soundManager.LoadSound("mystery-entrance", "res/fx/mystery-entrance.ogg")
	r.soundManager.LoadSound("pickup", "res/fx/pickup.ogg")
	r.soundManager.LoadSound("player-death", "res/fx/player-death-wash.ogg")
	r.soundManager.LoadSound("player-shoot", "res/fx/player-shoot.ogg")
	r.soundManager.LoadSound("cancel", "res/fx/cancel.ogg")
	r.soundManager.LoadSound("select", "res/fx/select.ogg")
	r.soundManager.LoadSound("bomb", "res/fx/bomb.ogg")
	r.soundManager.LoadSound("gas", "res/fx/gas.ogg")

	r.stars = []*Star{}
	for index := 0; index < 100; index += 1 {
		r.stars = append(r.stars, NewStar())
	}
	r.mystery = NewMystery(common.ScreenWidth, 20, 1, r)
	r.earth = NewEarth(r)
	return r
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
			r.PlaySound("select")
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
			r.PlaySound("cancel")
			r.QuitToMenu()
		}
		r.earth.Update(delta)
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
			r.PlaySound("cancel")
			r.QuitToMenu()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			r.PlaySound("select")
			r.QuitToMenu()
		}
		r.earth.Update(delta)
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
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			r.PlaySound("cancel")
			r.QuitToMenu()
		}
		r.earth.Update(delta)
	case gameWonGameState:
		r.player.Update(delta, r)
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			r.PlaySound("cancel")
			r.QuitToMenu()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			r.PlaySound("select")
			r.QuitToMenu()
		}
		r.earth.Update(delta)
	}

	for _, e := range r.effects {
		e.Update(delta)
		if e.animation.isDone {
			r.RemoveEffect(e)
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
		r.drawImage(screen, r.images["splash"], 40, 40)
		common.DrawText(screen, "controls", 60, 120)
		common.DrawText(screen, "press space to shoot", 26, 140)
		common.DrawText(screen, "wasd or arrow keys to move", 12, 150)
		common.DrawText(screen, "press space to play!", 26, 180)

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

func (r *Game) drawImage(screen *ebiten.Image, image *ebiten.Image, x float64, y float64) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	op.GeoM.Scale(common.Scale, common.Scale)
	screen.DrawImage(image, op)
}
