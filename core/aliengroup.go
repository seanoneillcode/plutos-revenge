package core

import (
	"math/rand"
	"plutos-revenge/common"
)

const alienYTargetNormal = alienSize

type AlienGroup struct {
	x               float64
	y               float64
	dir             int
	targetY         float64
	speed           float64
	timer           float64
	nextTimerAmount float64
}

func NewAlienGroup(game *Game, numberOfAliens int) *AlienGroup {
	group := &AlienGroup{
		x: alienSize,
		// put the group just out of the way above the screen
		y:     float64((numberOfAliens / 5) * alienSize * -2),
		dir:   1,
		speed: 10,
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
	if len(game.aliens) == 0 {
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
		for index := 0; index < len(game.aliens); index += 1 {
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
			if r.x > (53 + alienSize) {
				r.dir = -1
				r.moveDown(game.aliens)
			}
		} else {
			if r.x < 0 {
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
		r.nextTimerAmount = 0.5 + (rand.Float64() * 2)
		randIndex := rand.Intn(len(game.aliens))
		randAlien := game.aliens[randIndex]
		game.AddBullet(randAlien.x, randAlien.y+alienSize, 1, "alien")
	}
}

func (r *AlienGroup) moveDown(aliens []*Alien) {
	r.targetY = r.y + (alienSize * 2)
	r.speed = r.speed + 5
	for _, a := range aliens {
		if a.state == normalAlienState {
			a.dir = r.dir
		}
	}
}
