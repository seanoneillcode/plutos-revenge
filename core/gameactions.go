package core

import (
	"fmt"
	"plutos-revenge/common"
)

func (r *Game) StartNewGame() {
	r.effects = []*Effect{}
	fmt.Println("starting new game")
	r.bullets = []*Bullet{}
	r.aliens = []*Alien{}
	r.blocks = []*Block{
		NewBlock(30),
		NewBlock(90),
		NewBlock(150),
	}
	r.alienGroup = NewAlienGroup(r, 10)
	r.player = NewPlayer()
	r.state = playingGameState
	r.earth.Target(common.ScreenHeight - 24)
	r.player.Target(playerYNormal)
	r.score = 0
	r.level = 0
}

func (r *Game) StartNewLevel() {
	r.effects = []*Effect{}
	r.level = r.level + 1
	r.aliens = []*Alien{}
	r.alienGroup = NewAlienGroup(r, 5*(r.level+2))
	r.state = playingGameState
	r.blocks = []*Block{
		NewBlock(30),
		NewBlock(90),
		NewBlock(150),
	}
	r.PlaySound("level-start")
}

func (r *Game) GameOver() {
	if r.state == playingGameState {
		fmt.Println("game over")
		r.timer = 6 // seconds
		r.state = gameOverGameState
	}
	r.player.Target(common.ScreenHeight)
	r.alienGroup.targetY = common.ScreenHeight
	for _, b := range r.blocks {
		b.GetDestroyed(r)
	}
}

func (r *Game) QuitToMenu() {
	r.state = menuGameState
	r.player = nil
	r.bullets = nil
	r.aliens = nil
	r.alienGroup = nil
	r.earth.y = common.ScreenHeight
}

func (r *Game) EndLevel() {
	if r.level == numberOfLevels {
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
	r.player.Target(-common.ScreenHeight)
	r.earth.Target(common.ScreenHeight)
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

func (r *Game) RemoveEffect(effect *Effect) {
	var effects []*Effect
	for _, e := range r.effects {
		if e != effect {
			effects = append(effects, e)
		}
	}
	r.effects = effects
}

func (r *Game) RemoveAlien(alien *Alien) {
	var newAliens []*Alien
	for _, a := range r.aliens {
		if a != alien {
			newAliens = append(newAliens, a)
		}
	}
	r.aliens = newAliens
	r.alienGroup.SpeedUp()
}

func (r *Game) ScorePoint() {
	r.score = r.score + 1
}

func (r *Game) AddAlien(alien *Alien) {
	r.aliens = append(r.aliens, alien)
}

func (r *Game) PlaySound(name string) {
	r.soundManager.PlaySound(name)
}

func (r *Game) RemoveAdjacentAliens(source *Alien) {
	x := source.x - 16
	y := source.y - 16
	size := 44.0
	for _, a := range r.aliens {
		if common.Overlap(a.x, a.y, float64(a.size), x, y, size) {
			a.GetHit(r)
		}
	}
	r.AddEffect(source.x-12, source.y-12, "explosion")
	r.AddEffect(source.x-12, source.y+12, "explosion")
	r.AddEffect(source.x+12, source.y-12, "explosion")
	r.AddEffect(source.x+12, source.y+12, "explosion")
}

func (r *Game) AddEffect(x float64, y float64, kind string) {
	switch kind {
	case "explosion":
		r.effects = append(r.effects, &Effect{
			animation: &Animation{
				numFrames:       8,
				frameTimeAmount: 0.06,
				image:           r.images["explosion"],
				size:            12,
			},
			x: x - 4,
			y: y - 4,
		})
	case "player-death":
		r.effects = append(r.effects,
			&Effect{
				x: x - 6,
				y: y - 6,
				animation: &Animation{
					numFrames:       6,
					frameTimeAmount: 0.06,
					image:           r.images["player-death"],
					size:            24,
				},
			},
		)
	case "gas":
		r.effects = append(r.effects,
			&Effect{
				x: x - 6,
				y: y - 6,
				animation: &Animation{
					numFrames:       6,
					frameTimeAmount: 0.2,
					image:           r.images["gas"],
					size:            24,
				},
			},
		)
	case "alien-death":
		r.effects = append(r.effects,
			&Effect{
				x: x - 6,
				y: y - 6,
				animation: &Animation{
					numFrames:       8,
					frameTimeAmount: 0.04,
					image:           r.images["alien-death"],
					size:            24,
				},
			},
		)
	case "plus-one":
		r.effects = append(r.effects,
			&Effect{
				x: x,
				y: y,
				animation: &Animation{
					numFrames:       6,
					frameTimeAmount: 0.1,
					image:           r.images["plus-one"],
					size:            12,
				},
			},
		)
	}
}
