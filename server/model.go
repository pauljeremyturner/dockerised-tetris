package server

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pauljeremyturner/dockerised-tetris/shared"
	"math/rand"
	"sort"
)

type GameState struct {
	LineCount  int
	PieceCount int
	Pixels     []Pixel
	NextPiece  []Pixel
	GameOver   bool
	Duration   int64
}

type Lines struct {
	lineMap map[int][]Pixel
}

func (r *Lines) Compact(board shared.Board) int {

	keys := make([]int, len(r.lineMap))
	i := 0
	for k := range r.lineMap {
		keys[i] = k
		i++
	}
	sort.Ints(keys)

	keep := make([]int, 0)
	for _, y := range keys {
		line := r.lineMap[y]
		if len(line) != board.Width {
			keep = append(keep, y)
		}
	}

	y := board.Height - len(keep)
	newLineMap := make(map[int][]Pixel)
	for _, i := range keep {
		newLineMap[y] = r.lineMap[i]
		for _, p := range newLineMap[y] {
			p.Y = y
		}
		y++
	}

	r.lineMap = newLineMap

	return len(keys) - len(keep)
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

func (r *Piece) rotateClockwise() {
	centre := r.pixels[0]
	for i := range r.pixels {
		r.pixels[i].rotateClockwise(centre)
	}
}

func (r *Piece) rotateAnticlockwise() {
	centre := r.pixels[0]
	for i := range r.pixels {
		r.pixels[i].rotateAntiClockwise(centre)
	}
}

func (r *Piece) moveDown() {

	GetFileLogger().Println("Move Down")

	for i := range r.pixels {
		r.pixels[i].moveDown()
	}
}
func (r *Piece) moveLeft() {

	GetFileLogger().Println("Move Left")

	for i := range r.pixels {
		r.pixels[i].moveLeft()
	}
}
func (r *Piece) moveRight() {

	GetFileLogger().Println("Move Right")

	for i := range r.pixels {
		r.pixels[i].moveRight()
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
	/**
	  X
	  X
	  XX
	*/
	pixels := []Pixel{
		{X: 0, Y: 1, Color: 3},
		{X: 0, Y: 0, Color: 3},
		{X: 1, Y: 0, Color: 3},
		{X: 0, Y: 2, Color: 3},
	}
	return newPiece(pixels)
}

func NewP() Piece {
	/**
	  XX
	  X
	  X
	*/
	points := []Pixel{
		{X: 0, Y: 1, Color: 4},
		{X: 0, Y: 0, Color: 4},
		{X: 0, Y: 2, Color: 4},
		{X: 1, Y: 2, Color: 4},
	}
	return newPiece(points)
}

func NewS() Piece {
	/**
	   XX
	  XX
	*/
	pixels := []Pixel{
		{X: 1, Y: 0, Color: 5},
		{X: 0, Y: 0, Color: 5},
		{X: 1, Y: 1, Color: 5},
		{X: 2, Y: 1, Color: 5},
	}
	return newPiece(pixels)
}

func NewZ() Piece {
	/**
	  XX
	   XX
	*/
	pixels := []Pixel{
		{X: 1, Y: 1, Color: 6},
		{X: 0, Y: 1, Color: 6},
		{X: 1, Y: 0, Color: 6},
		{X: 2, Y: 0, Color: 6},
	}
	return newPiece(pixels)
}

func NewO() Piece {
	/**
	  XX
	  XX
	*/
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

type serverSession struct {
	player       Player
	moveQueue    chan shared.MoveType
	gameQueue    chan GameState
	activePiece  Piece
	lines        Lines
	nextPiece    Piece
	gameOver     bool
	lineCount    int
	pieceCount   int
	board        shared.Board
	startSeconds int64
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

func (r *Pixel) rotateClockwise(centre Pixel) {
	r.subtract(centre)
	newX := 0 - r.Y
	r.Y = r.X
	r.X = newX
	r.add(centre)
}

func (r *Pixel) rotateAntiClockwise(centre Pixel) {
	r.subtract(centre)
	newY := 0 - r.X
	r.X = r.Y
	r.Y = newY
	r.add(centre)
}

func (r *Pixel) moveDown() {
	r.Y = r.Y + 1
}

func (r *Pixel) moveLeft() {
	r.X = r.X - 1
}

func (r *Pixel) moveRight() {
	r.X = r.X + 1
}

func (r *serverSession) MoveActivePieceDownIfPossible() bool {

	ap := r.activePiece

	pieceCopy := ap.Clone()
	pieceCopy.moveDown()

	if isMovePossible(r, pieceCopy) {
		r.activePiece = *pieceCopy
		return true
	} else {
		return false
	}
}

func (r *serverSession) MoveActivePieceRightIfPossible() bool {

	ap := r.activePiece

	pieceCopy := ap.Clone()
	pieceCopy.moveRight()

	if isMovePossible(r, pieceCopy) {
		r.activePiece = *pieceCopy
		return true
	} else {
		return false
	}
}

func (r *serverSession) MoveActivePieceLeftIfPossible() bool {

	ap := r.activePiece

	pieceCopy := ap.Clone()
	pieceCopy.moveLeft()

	if isMovePossible(r, pieceCopy) {
		r.activePiece = *pieceCopy
		return true
	} else {
		return false
	}
}

func (r *serverSession) RotateActivePieceClockwiseIfPossible() bool {

	ap := r.activePiece

	pieceCopy := ap.Clone()
	pieceCopy.rotateClockwise()

	if isMovePossible(r, pieceCopy) {
		r.activePiece = *pieceCopy
		return true
	} else {
		return false
	}
}

func (r *serverSession) RotateActivePieceAnticlockwiseIfPossible() bool {

	ap := r.activePiece

	pieceCopy := ap.Clone()
	pieceCopy.rotateAnticlockwise()

	if isMovePossible(r, pieceCopy) {
		r.activePiece = *pieceCopy
		return true
	} else {
		return false
	}
}

func isMovePossible(s *serverSession, placement *Piece) bool {

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
