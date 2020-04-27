package shared

import "log"

type MoveType rune

type Logger struct {
	*log.Logger
}

const (
	MOVELEFT    = 's'
	MOVERIGHT   = 'd'
	ROTATELEFT  = 'a'
	ROTATERIGHT = 'f'
	DROP        = 'e'
	DOWN        = 'x'
)

type Pixel struct {
	X     int
	Y     int
	Color int
}

func (r *Pixel) subtract(p Pixel) {
	r.X = r.X - p.X
	r.Y = r.Y - p.Y
}
func (r *Pixel) add(p Pixel) {
	r.X = r.X + p.X
	r.Y = r.Y + p.Y
}

func (r *Pixel) RotateClockwise(centre Pixel) {
	r.subtract(centre)
	newX := 0 - r.Y
	r.Y = r.X
	r.X = newX
	r.add(centre)
}

func (r *Pixel) RotateAntiClockwise(centre Pixel) {
	r.subtract(centre)
	newY := 0 - r.X
	r.X = r.Y
	r.Y = newY
	r.add(centre)
}

func (r *Pixel) MoveDown() {
	r.Y = r.Y - 1
}

func (r *Pixel) MoveLeft() {
	r.X = r.X - 1
}

func (r *Pixel) MoveRight() {
	r.X = r.X + 1
}
