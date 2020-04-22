package server

import (
	"github.com/google/uuid"
	"github.com/pauljeremyturner/dockerised-tetris/shared"
)

type GameState struct {
	Score     int
	Pixels    []Pixel
	NextPiece []Pixel
}

type ServerSession struct {
	player         Player
	moveQueue      chan shared.MoveType
	gameQueue      chan GameState
	gameOverSignal chan bool
	activePiece    Piece
	lines          Lines
	nextPiece      Piece
	gameOver       bool
	score          int
	board          shared.Board
}

type Player struct {
	uuid       uuid.UUID
	playerName string
}

type Pixel struct {
	X     int
	Y     int
	Color int
}

func (r *Pixel) SameLocationAs(p Pixel) bool {
	return (r.X == p.X) && (r.Y == p.Y)
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
	r.Y = r.Y + 1
}

func (r *Pixel) MoveLeft() {
	r.X = r.X - 1
}

func (r *Pixel) MoveRight() {
	r.X = r.X + 1
}

func (r *ServerSession) MoveActivePieceDownIfPossible() bool {

	ap := &r.activePiece

	pieceCopy := ap.Clone()
	pieceCopy.MoveDown()

	if isMovePossible(r, pieceCopy) {
		r.activePiece = *pieceCopy
		return true
	} else {
		return false
	}
}

func (r *ServerSession) MoveActivePieceRightIfPossible() bool {

	ap := &r.activePiece

	pieceCopy := ap.Clone()
	pieceCopy.MoveRight()

	if isMovePossible(r, pieceCopy) {
		r.activePiece = *pieceCopy
		return true
	} else {
		return false
	}
}

func (r *ServerSession) MoveActivePieceLeftIfPossible() bool {

	ap := &r.activePiece

	pieceCopy := ap.Clone()
	pieceCopy.MoveLeft()

	if isMovePossible(r, pieceCopy) {
		r.activePiece = *pieceCopy
		return true
	} else {
		return false
	}
}

func (r *ServerSession) RotateActivePieceClockwiseIfPossible() bool {

	ap := &r.activePiece

	pieceCopy := ap.Clone()
	pieceCopy.RotateClockwise()

	if isMovePossible(r, pieceCopy) {
		r.activePiece = *pieceCopy
		return true
	} else {
		return false
	}
}

func (r *ServerSession) RotateActivePieceAnticlockwiseIfPossible() bool {

	ap := &r.activePiece

	pieceCopy := ap.Clone()
	pieceCopy.RotateAnticlockwise()

	if isMovePossible(r, pieceCopy) {
		r.activePiece = *pieceCopy
		return true
	} else {
		return false
	}
}

func isMovePossible(s *ServerSession, placement *Piece) bool {

	for _, pp := range placement.pixels {
		if pp.X < 0 || pp.X >= s.board.Width || pp.Y >= s.board.Height {
			return false
		}
	}

	for _, p := range placement.pixels {

		if lps, ok := s.lines.lineMap[p.Y]; ok {
			for _, lp := range lps {
				if p.SameLocationAs(lp) {
					return false
				}
			}
		}
	}
	return true
}
