package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"plutos-revenge/common"
)

type Block struct {
	x     float64
	y     float64
	image *ebiten.Image
	lives int
	size  int
}

func NewBlock(x float64, game *Game) *Block {
	return &Block{
		x:     x,
		y:     playerYNormal - 16,
		lives: 3,
		size:  12,
		image: game.images["block"],
	}
}

func (r *Block) Update(delta float64, game *Game) {
	if r.lives > 0 {
		for _, a := range game.aliens {
			if common.Overlap(a.x, a.y, float64(a.size), r.x, r.y, float64(r.size)) {
				r.GetDestroyed(game)
			}
		}
	}
}

func (r *Block) Draw(screen *ebiten.Image) {
	if r.lives > 0 {
		frame := r.lives - 1
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(r.x, r.y)
		op.GeoM.Scale(common.Scale, common.Scale)
		screen.DrawImage(r.image.SubImage(image.Rect(frame*r.size, 0, (frame+1)*r.size, r.size)).(*ebiten.Image), op)
	}
}

func (r *Block) GetHit(game *Game) {
	r.lives = r.lives - 1
	if r.lives == 0 {
		game.AddEffect(r.x, r.y, "explosion")
	}
	game.PlaySound("block")
}

func (r *Block) GetDestroyed(game *Game) {
	r.lives = 0
	game.AddEffect(r.x, r.y, "explosion")
	game.PlaySound("blast")
}
