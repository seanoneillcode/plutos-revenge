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

type Game struct {
	lastUpdateCalled time.Time
	player           *Player
	bullets          []*Bullet
	aliens           []*Alien
	alienGroup       *AlienGroup
	timer            float64
	state            string
	images           map[string]*ebiten.Image
	level            int
	score            int
}

func NewGame() *Game {
	return &Game{
		state: menuGameState,
		images: map[string]*ebiten.Image{
			"splash": common.LoadImage("splash.png"),
		},
	}
}

func (r *Game) Update() error {
	delta := float64(time.Now().Sub(r.lastUpdateCalled).Milliseconds()) / 1000
	r.lastUpdateCalled = time.Now()

	switch r.state {
	case menuGameState:
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			r.StartNewGame()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			return common.NormalEscapeError
		}
	case playingGameState:
		r.player.Update(delta, r)
		for _, b := range r.bullets {
			b.Update(delta, r)
		}
		r.alienGroup.Update(delta, r)
		for _, a := range r.aliens {
			a.Update(delta, r)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			r.QuitToMenu()
		}
	case gameOverGameState:
		for _, b := range r.bullets {
			b.Update(delta, r)
		}
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

	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	return nil
}

func (r *Game) Draw(screen *ebiten.Image) {

	switch r.state {
	case menuGameState:
		r.drawImage(screen, "splash", 40, 40)
		common.DrawText(screen, "start game", 60, 120)
	case playingGameState:
		r.player.Draw(screen)
		for _, b := range r.bullets {
			b.Draw(screen)
		}
		for _, a := range r.aliens {
			a.Draw(screen)
		}
		common.DrawText(screen, fmt.Sprintf("score %d", r.score), 4, 4)
		common.DrawText(screen, fmt.Sprintf("lives %d", r.player.lives), 130, 4)
	case gameOverGameState:
		r.player.Draw(screen)
		for _, b := range r.bullets {
			b.Draw(screen)
		}
		for _, a := range r.aliens {
			a.Draw(screen)
		}
		common.DrawText(screen, "game over", 60, 90)
	case levelOverGameState:
		r.player.Draw(screen)
		for _, b := range r.bullets {
			b.Draw(screen)
		}
		common.DrawText(screen, "   wave\ndestroyed", 60, 90)
	case gameWonGameState:
		r.player.Draw(screen)
		common.DrawText(screen, " you win!", 60, 70)
		common.DrawText(screen, "the earth\nis saved!", 60, 120)
	}

}

func (r *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return common.ScreenWidth * common.Scale, common.ScreenHeight * common.Scale
}

func (r *Game) StartNewGame() {
	fmt.Println("starting new game")
	r.bullets = []*Bullet{}
	r.aliens = []*Alien{}
	r.alienGroup = NewAlienGroup(r, 10)
	r.player = NewPlayer()
	r.state = playingGameState
}

func (r *Game) StartNewLevel() {
	r.level = r.level + 1
	r.aliens = []*Alien{}
	r.alienGroup = NewAlienGroup(r, 5*(r.level+2))
	r.state = playingGameState
}

func (r *Game) GameOver() {
	if r.state == playingGameState {
		fmt.Println("game over")
		r.timer = 6 // seconds
		r.state = gameOverGameState
	}
}

func (r *Game) QuitToMenu() {
	r.state = menuGameState
	r.player = nil
	r.bullets = nil
	r.aliens = nil
	r.alienGroup = nil
}

func (r *Game) EndLevel() {
	if r.level == 1 {
		r.WinGame()
		return
	}
	if r.state == playingGameState {
		fmt.Println("ending level")
		r.state = levelOverGameState
		r.timer = 3.0
	}
}

func (r *Game) WinGame() {
	r.state = gameWonGameState
}

func (r *Game) AddBullet(x float64, y float64, dir int, kind string) {
	r.bullets = append(r.bullets, NewBullet(x, y, dir, kind))
}

func (r *Game) RemoveBullet(bullet *Bullet) {
	var newBullets []*Bullet
	for _, b := range r.bullets {
		if b != bullet {
			newBullets = append(newBullets, b)
		}
	}
	r.bullets = newBullets
}

func (r *Game) RemoveAlien(alien *Alien) {
	var newAliens []*Alien
	for _, a := range r.aliens {
		if a != alien {
			newAliens = append(newAliens, a)
		}
	}
	r.aliens = newAliens
}

func (r *Game) ScorePoint() {
	r.score = r.score + 1
}

func (r *Game) drawImage(screen *ebiten.Image, img string, x float64, y float64) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	op.GeoM.Scale(common.Scale, common.Scale)
	screen.DrawImage(r.images[img], op)
}

func (r *Game) AddAlien(alien *Alien) {
	r.aliens = append(r.aliens, alien)
}
