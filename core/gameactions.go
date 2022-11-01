package core

import (
	"fmt"
	"plutos-revenge/common"
)

func (r *Game) StartNewGame() {
	r.effects = []*Animation{}
	fmt.Println("starting new game")
	r.bullets = []*Bullet{}
	r.aliens = []*Alien{}
	r.alienGroup = NewAlienGroup(r, 10)
	r.player = NewPlayer()
	r.state = playingGameState
	r.earth.Target(common.ScreenHeight - 24)
}

func (r *Game) StartNewLevel() {
	r.effects = []*Animation{}
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
	r.earth.Target(common.ScreenHeight)
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

func (r *Game) AddAlien(alien *Alien) {
	r.aliens = append(r.aliens, alien)
}

func (r *Game) AddEffect(x float64, y float64, kind string) {
	switch kind {
	case "explosion":
		r.effects = append(r.effects, &Animation{
			numFrames:       8,
			frameTimeAmount: 0.06,
			image:           r.images["explosion"],
			size:            12,
			x:               x,
			y:               y,
		})
	}
}
