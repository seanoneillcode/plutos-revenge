package core

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"plutos-revenge/common"
	"time"
)

type Game struct {
	lastUpdateCalled time.Time
	player           *Player
	bullets          []*Bullet
	timer            float64
}

func NewGame() *Game {
	return &Game{
		player:  NewPlayer(),
		bullets: []*Bullet{},
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
	r.player.Update(delta, r)
	for _, b := range r.bullets {
		b.Update(delta, r)
	}

	r.timer = r.timer + delta
	if r.timer > 1.5 {
		r.timer = 0
		r.AddBullet(90, 10, 1)
	}

	return nil
}

func (r *Game) Draw(screen *ebiten.Image) {
	r.player.Draw(screen)
	for _, b := range r.bullets {
		b.Draw(screen)
	}
}

func (r *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return common.ScreenWidth * common.Scale, common.ScreenHeight * common.Scale
}

func (r *Game) PlayerDeath() {
	// change state..
	// show game over text
	// change to high scores screen
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
