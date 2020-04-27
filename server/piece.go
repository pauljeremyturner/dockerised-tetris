package server

import (
	"math/rand"

	"github.com/pauljeremyturner/dockerised-tetris/shared"
)

type LineState struct {
	pixels []shared.Pixel
}

type Line interface {
	Pixels() []shared.Pixel
}

func (r *LineState) Pixels() []shared.Pixel {
	return r.pixels
}

func (r *LineState) MoveDown() {
	for _, p := range r.pixels {
		p.MoveDown()
	}
}

func (r *PieceState) RotateClockwiseIfPossible() bool {

	copyPixels := make([]shared.Pixel, len(r.pixels), cap(r.pixels))

	pieceCopy := PieceState{pixels: copyPixels}

	for _, pc := range r.pixels {
		copyPixels = append(copyPixels, shared.Pixel{X: pc.X, Y: pc.Y, Color: pc.Color})
	}

	pieceCopy.RotateClockwise()

	for _, p := range pieceCopy.pixels {
		if p.Y < 0 {
			return false
		}
	}

	r.pixels = pieceCopy.pixels

}

func (r *LineState) Fits(p PieceState) bool {
	var highestPiecePixels []shared.Pixel
	var highestYposition = 0
	for _, p := range p.pixels {
		if highestYposition < p.Y {
			highestPiecePixels = nil
			highestYposition = p.Y
			highestPiecePixels = append(highestPiecePixels, p)
		} else if highestYposition == p.Y {
			highestPiecePixels = append(highestPiecePixels, p)
		}
	}

	fits := true

	for _, lp := range r.pixels {
		for _, pp := range p.pixels {
			if lp.X == pp.X {
				fits = false
				break
			}
		}
	}
	return fits

}

type Piece interface {
	RotateClockwise()
	RotateAnticlockwise()
	MoveLeft()
	MoveRight()
	MoveDown()
	Pixels() []shared.Pixel
}

func newPiece(ps []shared.Pixel) Piece {
	return &PieceState{pixels: ps}
}

type PieceState struct {
	pixels []shared.Pixel
}

func (r *PieceState) Pixels() []shared.Pixel {
	return r.pixels
}

func RandomPiece() Piece {
	switch rand.Intn(8) {
	case 0:
		return NewI()
	case 1:
		return NewL()
	case 2:
		return NewO()
	case 3:
		return NewP()
	case 4:
		return NewS()
	case 5:
		return NewT()
	case 6:
		return NewZ()
	}
	panic("Expecting a piece to return")
}

func (r *PieceState) RotateClockwise() {
	centre := r.pixels[0]
	for i, p := range r.pixels {
		if i == 0 {
			continue
		}
		p.RotateClockwise(centre)
	}
}

func (r *PieceState) RotateAnticlockwise() {
	centre := r.pixels[0]
	for i, p := range r.pixels {
		if i == 0 {
			continue
		}
		p.RotateAntiClockwise(centre)
	}
}

func (r *PieceState) MoveDown() {
	for _, p := range r.pixels {
		p.MoveDown()
	}
}
func (r *PieceState) MoveLeft() {
	for _, p := range r.pixels {
		p.MoveLeft()
	}
}
func (r *PieceState) MoveRight() {
	for _, p := range r.pixels {
		p.MoveRight()
	}
}

func NewI() Piece {
	points := []shared.Pixel{
		{X: 0, Y: 1, Color: 0},
		{X: 0, Y: 0, Color: 0},
		{X: 0, Y: 2, Color: 0},
		{X: 0, Y: 3, Color: 0},
	}
	return newPiece(points)
}

func NewT() Piece {
	points := []shared.Pixel{
		{X: 1, Y: 1, Color: 0},
		{X: 1, Y: 0, Color: 0},
		{X: 0, Y: 1, Color: 0},
		{X: 2, Y: 1, Color: 0},
	}
	return newPiece(points)
}

func NewL() Piece {
	points := []shared.Pixel{
		{X: 0, Y: 1, Color: 0},
		{X: 0, Y: 0, Color: 0},
		{X: 1, Y: 0, Color: 0},
		{X: 0, Y: 2, Color: 0},
	}
	return newPiece(points)
}

func NewP() Piece {
	points := []shared.Pixel{
		{X: 0, Y: 1, Color: 0},
		{X: 0, Y: 0, Color: 0},
		{X: 0, Y: 2, Color: 0},
		{X: 1, Y: 2, Color: 0},
	}
	return newPiece(points)
}

func NewS() Piece {
	points := []shared.Pixel{
		{X: 1, Y: 0, Color: 0},
		{X: 0, Y: 0, Color: 0},
		{X: 1, Y: 1, Color: 0},
		{X: 2, Y: 1, Color: 0},
	}
	return newPiece(points)
}

func NewZ() Piece {
	points := []shared.Pixel{
		{X: 1, Y: 1, Color: 0},
		{X: 0, Y: 1, Color: 0},
		{X: 1, Y: 0, Color: 0},
		{X: 2, Y: 0, Color: 0},
	}
	return newPiece(points)
}

func NewO() Piece {
	points := []shared.Pixel{
		{X: 1, Y: 0, Color: 0},
		{X: 0, Y: 0, Color: 0},
		{X: 0, Y: 1, Color: 0},
		{X: 1, Y: 1, Color: 0},
	}
	return newPiece(points)
}
