package server

import (
	"fmt"
	"math/rand"
)

type Lines struct {
	lineMap map[int][]Pixel
}

func (r *Lines) LinesBlockPieceMoveDown(p *Piece) bool {
	for _, p := range p.pixels {

		if lps, ok := r.lineMap[p.Y]; ok {
			for _, lp := range lps {
				if p.SameLocationAs(lp) {
					return true
				}
			}
		}
	}
	return false
}

func (r *Piece) lowestPixels() map[int]Pixel {

	lowestPixelMap := make(map[int]Pixel)
	for _, pxl := range r.pixels {
		if p, ok := lowestPixelMap[pxl.X]; ok {
			if p.Y < pxl.Y {
				lowestPixelMap[pxl.X] = pxl
			}
		} else {
			lowestPixelMap[pxl.X] = pxl
		}
	}
	return lowestPixelMap
}

func (r *Piece) Clone() *Piece {
	copyPixels := make([]Pixel, len(r.pixels))

	pieceCopy := &Piece{pixels: copyPixels}

	copy(copyPixels, r.pixels)

	return pieceCopy
}

func newPiece(ps []Pixel) Piece {
	return Piece{pixels: ps}
}

type Piece struct {
	// pixels[0] is used as the centre piece for calculating rotations
	pixels []Pixel
}

func RandomPiece() Piece {
	switch rand.Intn(7) {
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

func (r *Piece) RotateClockwise() {
	centre := r.pixels[0]
	for i := range r.pixels {
		r.pixels[i].RotateClockwise(centre)
	}
}

func (r *Piece) RotateAnticlockwise() {
	centre := r.pixels[0]
	for i := range r.pixels {
		r.pixels[i].RotateAntiClockwise(centre)
	}
}

func (r *Piece) MoveDown() {

	GetFileLogger().Println("Move Down")

	for i := range r.pixels {
		r.pixels[i].MoveDown()
	}
}
func (r *Piece) MoveLeft() {

	GetFileLogger().Println("Move Left")

	for i := range r.pixels {
		r.pixels[i].MoveLeft()
	}
}
func (r *Piece) MoveRight() {

	GetFileLogger().Println("Move Right")

	for i := range r.pixels {
		r.pixels[i].MoveRight()
	}
}

func NewI() Piece {
	/*
	   XXXXXX
	*/
	points := []Pixel{
		{X: 0, Y: 1, Color: 1},
		{X: 0, Y: 0, Color: 1},
		{X: 0, Y: 2, Color: 1},
		{X: 0, Y: 3, Color: 1},
	}
	return newPiece(points)
}

func NewT() Piece {
	/*
		     X
			XXX
	*/
	pixels := []Pixel{
		{X: 1, Y: 1, Color: 2},
		{X: 1, Y: 0, Color: 2},
		{X: 0, Y: 1, Color: 2},
		{X: 2, Y: 1, Color: 2},
	}
	return newPiece(pixels)
}

func NewL() Piece {
	pixels := []Pixel{
		{X: 0, Y: 1, Color: 3},
		{X: 0, Y: 0, Color: 3},
		{X: 1, Y: 0, Color: 3},
		{X: 0, Y: 2, Color: 3},
	}
	return newPiece(pixels)
}

func NewP() Piece {
	points := []Pixel{
		{X: 0, Y: 1, Color: 4},
		{X: 0, Y: 0, Color: 4},
		{X: 0, Y: 2, Color: 4},
		{X: 1, Y: 2, Color: 4},
	}
	return newPiece(points)
}

func NewS() Piece {
	pixels := []Pixel{
		{X: 1, Y: 0, Color: 5},
		{X: 0, Y: 0, Color: 5},
		{X: 1, Y: 1, Color: 5},
		{X: 2, Y: 1, Color: 5},
	}
	return newPiece(pixels)
}

func NewZ() Piece {
	pixels := []Pixel{
		{X: 1, Y: 1, Color: 6},
		{X: 0, Y: 1, Color: 6},
		{X: 1, Y: 0, Color: 6},
		{X: 2, Y: 0, Color: 6},
	}
	return newPiece(pixels)
}

func NewO() Piece {
	pixels := []Pixel{
		{X: 1, Y: 0, Color: 7},
		{X: 0, Y: 0, Color: 7},
		{X: 0, Y: 1, Color: 7},
		{X: 1, Y: 1, Color: 7},
	}
	return newPiece(pixels)
}

func (r *Piece) String() string {
	var s = ""
	for _, pix := range r.pixels {
		s = s + fmt.Sprintf("pixel, (%d, %d) color %d; ", pix.X, pix.Y, pix.Color)
	}
	return s
}
