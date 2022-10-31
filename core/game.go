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

type Game struct {
	lastUpdateCalled time.Time
	player           *Player
	bullets          []*Bullet
	timer            float64
	state            string
	images           map[string]*ebiten.Image
}

func NewGame() *Game {
	return &Game{
		//player:  NewPlayer(),
		//bullets: []*Bullet{},
		state: menuGameState,
		images: map[string]*ebiten.Image{
			"game-over-text": common.LoadImage("game-over-text.png"),
			"play-text":      common.LoadImage("play-text.png"),
		},
	}
}

func (r *Game) Update() error {
	delta := float64(time.Now().Sub(r.lastUpdateCalled).Milliseconds()) / 1000
	r.lastUpdateCalled = time.Now()

	switch r.state {
	case menuGameState:
		if ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsKeyPressed(ebiten.KeyEnter) {
			r.StartNewGame()
		}
	case playingGameState:
		r.player.Update(delta, r)
		for _, b := range r.bullets {
			b.Update(delta, r)
		}

		r.timer = r.timer + delta
		if r.timer > 1.5 {
			r.timer = 0
			r.AddBullet(90, 10, 1)
		}
	case gameOverGameState:
		if ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsKeyPressed(ebiten.KeyEnter) {
			return common.NormalEscapeError
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return common.NormalEscapeError
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	return nil
}

func (r *Game) Draw(screen *ebiten.Image) {

	switch r.state {
	case menuGameState:
		r.drawImage(screen, "play-text", 40, 40)
	case playingGameState:
		r.player.Draw(screen)
		for _, b := range r.bullets {
			b.Draw(screen)
		}
	case gameOverGameState:
		r.player.Draw(screen)
		for _, b := range r.bullets {
			b.Draw(screen)
		}
		r.drawImage(screen, "game-over-text", 40, 40)
	}

}

func (r *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return common.ScreenWidth * common.Scale, common.ScreenHeight * common.Scale
}

func (r *Game) StartNewGame() {
	r.bullets = []*Bullet{}
	r.player = NewPlayer()
	r.state = playingGameState
}

func (r *Game) PlayerDeath() {
	r.state = gameOverGameState
}

func (r *Game) AddBullet(x float64, y float64, dir int) {
	r.bullets = append(r.bullets, NewBullet(x, y, dir))
}

func (r *Game) RemoveBullet(bullet *Bullet) {
	var newBullets []*Bullet
	for _, b := range r.bullets {
		if b != bullet {
			newBullets = append(newBullets, b)
		} else {
			fmt.Println("removing bullet")
			fmt.Printf("num bullets %d\n", len(r.bullets))
		}
	}
	r.bullets = newBullets
}

func (r *Game) drawImage(screen *ebiten.Image, img string, x float64, y float64) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	op.GeoM.Scale(common.Scale, common.Scale)
	screen.DrawImage(r.images[img], op)
}
