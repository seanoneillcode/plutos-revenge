package core

import (
	"fmt"
	"math"
	"math/rand"
	"plutos-revenge/common"
)

const alienYTargetNormal = alienSize
const maxAlienSpeed = 50

type AlienGroup struct {
	x                 float64
	y                 float64
	dir               int
	targetY           float64
	speed             float64
	timer             float64
	nextTimerAmount   float64
	shootingRate      float64
	originalNumAliens int
	speedIncrease     float64
}

func NewAlienGroup(game *Game, numberOfAliens int) *AlienGroup {
	group := &AlienGroup{
		x: alienSize,
		// put the group just above the screen
		y:                 float64((numberOfAliens / 5) * alienSize * -2),
		dir:               1,
		speed:             10,
		shootingRate:      math.Max(0.08, 0.4-(float64(numberOfAliens)*0.01)),
		originalNumAliens: numberOfAliens,
		speedIncrease:     1.0 / float64(numberOfAliens),
	}
	// add the aliens
	x := alienSize + 6.0
	y := group.y
	for index := 0; index < numberOfAliens; index += 1 {
		alien := NewAlien(x, y, 1)
		game.AddAlien(alien)
		x = x + (alienSize * 2)
		if (index+1)%5 == 0 {
			y = y + (alienSize * 2)
			if (index+1)%2 == 0 {
				x = alienSize + 6
			} else {
				x = alienSize
			}
		}
	}
	group.targetY = alienYTargetNormal
	return group
}

func (r *AlienGroup) Update(delta float64, game *Game) {
	if r.targetY > common.ScreenHeight {
		game.GameOver()
		return
	}
	numAliens := len(game.aliens)
	if numAliens == 0 {
		game.EndLevel()
		return
	}
	actualSpeed := r.speed
	if r.targetY == alienYTargetNormal {
		actualSpeed = actualSpeed * 3
	}
	var numAlive = 0
	if r.y < r.targetY {
		// move down y
		r.y = r.y + (delta * actualSpeed)
		for index := 0; index < numAliens; index += 1 {
			alien := game.aliens[index]
			if alien.state == normalAlienState {
				alien.y = alien.y + (delta * actualSpeed)
				numAlive = numAlive + 1
			}
		}
	} else {
		// move across x
		r.x = r.x + (float64(r.dir) * delta * r.speed)
		if r.dir == 1 {
			hitEdge := false
			for _, a := range game.aliens {
				if a.x+float64(a.size) > common.ScreenWidth {
					hitEdge = true
				}
			}
			if hitEdge {
				r.dir = -1
				r.moveDown(game.aliens)
			}
		} else {
			hitEdge := false
			for _, a := range game.aliens {
				if a.x < 0 {
					hitEdge = true
				}
			}
			if hitEdge {
				r.dir = 1
				r.moveDown(game.aliens)
			}
		}
		for index := 0; index < len(game.aliens); index += 1 {
			alien := game.aliens[index]
			if alien.state == normalAlienState {
				numAlive = numAlive + 1
				alien.x = alien.x + (float64(r.dir) * delta * r.speed)
			}
		}
	}
	// shooting
	r.timer = r.timer + delta
	if r.timer > r.nextTimerAmount && numAlive > 0 {
		r.timer = 0
		r.nextTimerAmount = r.shootingRate + (rand.Float64())
		randIndex := rand.Intn(numAliens)
		shootingAlien := game.aliens[randIndex]
		if shootingAlien.state == hitState {
			for _, a := range game.aliens {
				if a.state == normalAlienState {
					shootingAlien = a
					continue
				}
			}
		}
		// 50% to try shoot player directly
		if rand.Float64() > 0.7 {
			for _, a := range game.aliens {
				if a.x > game.player.x-4 && a.x < game.player.x+float64(game.player.size+4) {
					shootingAlien = a
					break
				}
			}
		}
		game.AddBullet(shootingAlien.x, shootingAlien.y+alienSize, 1, "alien")
	}
}

func (r *AlienGroup) moveDown(aliens []*Alien) {
	r.targetY = r.y + (alienSize)
	for _, a := range aliens {
		if a.state == normalAlienState {
			a.dir = r.dir
		}
	}
}

func (r *AlienGroup) SpeedUp() {
	r.speed = r.speed + (r.speedIncrease * 50)
	fmt.Printf("speed: %v\n", r.speed)
}
