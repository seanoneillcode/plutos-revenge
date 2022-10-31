package core

import "plutos-revenge/common"

type AlienGroup struct {
	x       float64
	y       float64
	dir     int
	targetY float64
	speed   float64
}

func NewAlienGroup(game *Game, numberOfAliens int) *AlienGroup {
	group := &AlienGroup{
		x:     alienSize,
		y:     alienSize,
		dir:   1,
		speed: 20,
	}
	// add the aliens
	x := alienSize + 6
	y := 6
	for index := 0; index < numberOfAliens; index += 1 {
		alien := NewAlien(float64(x), float64(y), 1)
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
	if r.y < r.targetY {
		// move down y
		r.y = r.y + (delta * r.speed)
		for index := 0; index < len(game.aliens); index += 1 {
			alien := game.aliens[index]
			if alien.state == normalAlienState {
				alien.y = alien.y + (delta * r.speed)
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
				alien.x = alien.x + (float64(r.dir) * delta * r.speed)
			}
		}
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
