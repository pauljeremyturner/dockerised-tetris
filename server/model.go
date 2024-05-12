package server

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pauljeremyturner/dockerised-tetris/protogen"
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

//go:generate stringer -type=PieceType
type PieceType int

const (
	PieceTypeUnknown PieceType = iota - 1
	PieceTypeI
	PieceTypeO
	PieceTypeT
	PieceTypeL
	PieceTypeP
	PieceTypeS
	PieceTypeZ
)

//go:generate stringer -type=Color
type Color int

const (
	ColorUnknown Color = iota - 1
	ColorMajenta
	ColorCyan
	ColorYellow
	ColorBlue
	ColorGreen
	ColorRed
	ColorWhite
	ColorBlack
)

func colorToProto(c Color) protogen.Square_ColorEnum {

	switch c {
	case ColorMajenta:
		return protogen.Square_MAGENTA
	case ColorCyan:
		return protogen.Square_CYAN
	case ColorYellow:
		return protogen.Square_YELLOW
	case ColorBlue:
		return protogen.Square_BLUE
	case ColorGreen:
		return protogen.Square_GREEN
	case ColorRed:
		return protogen.Square_RED
	case ColorWhite:
		return protogen.Square_WHITE
	default:
		return protogen.Square_BLACK
	}
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

func newPiece(ps []Pixel, pt PieceType) Piece {
	return Piece{
		pixels:    ps,
		pieceType: pt,
	}
}

type Piece struct {
	// pixels[0] is used as the centre piece for calculating rotations
	pixels    []Pixel
	pieceType PieceType
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
	default:
		return NewZ()
	}
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

	for i := range r.pixels {
		r.pixels[i].moveDown()
	}
}
func (r *Piece) moveLeft() {

	for i := range r.pixels {
		r.pixels[i].moveLeft()
	}
}
func (r *Piece) moveRight() {

	for i := range r.pixels {
		r.pixels[i].moveRight()
	}
}

func NewO() Piece {
	/**
	  XX
	  XX
	*/
	pixels := []Pixel{
		{X: 1, Y: 0, Color: ColorMajenta},
		{X: 0, Y: 0, Color: ColorMajenta},
		{X: 0, Y: 1, Color: ColorMajenta},
		{X: 1, Y: 1, Color: ColorMajenta},
	}
	return newPiece(pixels, PieceTypeO)
}

func NewI() Piece {
	/*
	   XXXXXX
	*/
	points := []Pixel{
		{X: 0, Y: 1, Color: ColorCyan},
		{X: 0, Y: 0, Color: ColorCyan},
		{X: 0, Y: 2, Color: ColorCyan},
		{X: 0, Y: 3, Color: ColorCyan},
	}
	return newPiece(points, PieceTypeI)
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
	return newPiece(pixels, PieceTypeT)
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
	return newPiece(pixels, PieceTypeL)
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
	return newPiece(points, PieceTypeP)
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
	return newPiece(pixels, PieceTypeS)
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
	return newPiece(pixels, PieceTypeZ)
}

func (r *Piece) String() string {
	var s = "Piece: " + r.pieceType.String()
	for _, pix := range r.pixels {
		s = s + pix.String()
	}
	return s
}

type Player struct {
	UUID uuid.UUID
	Name string
}

func (p Player) String() string {
	return fmt.Sprintf("player, Name: %s, UUID: %s", p.Name, p.UUID.String())
}

type Pixel struct {
	X     int
	Y     int
	Color Color
}

func (r *Pixel) String() string {
	return fmt.Sprintf("Pixel: (%d, %d) %s", r.X, r.Y, r.Color.String())
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

func (s *serverSession) MoveActivePieceDownIfPossible() bool {

	ap := s.activePiece

	pieceCopy := ap.Clone()
	pieceCopy.moveDown()

	if isMovePossible(s, pieceCopy) {
		s.activePiece = *pieceCopy
		return true
	} else {
		return false
	}
}

func (s *serverSession) MoveActivePieceRightIfPossible() bool {

	ap := s.activePiece

	pieceCopy := ap.Clone()
	pieceCopy.moveRight()

	if isMovePossible(s, pieceCopy) {
		s.activePiece = *pieceCopy
		return true
	} else {
		return false
	}
}

func (s *serverSession) MoveActivePieceLeftIfPossible() bool {

	ap := s.activePiece

	pieceCopy := ap.Clone()
	pieceCopy.moveLeft()

	if isMovePossible(s, pieceCopy) {
		s.activePiece = *pieceCopy
		return true
	} else {
		return false
	}
}

func (s *serverSession) RotateActivePieceClockwiseIfPossible() bool {

	ap := s.activePiece

	pieceCopy := ap.Clone()
	pieceCopy.rotateClockwise()

	if isMovePossible(s, pieceCopy) {
		s.activePiece = *pieceCopy
		return true
	} else {
		return false
	}
}

func (s *serverSession) RotateActivePieceAnticlockwiseIfPossible() bool {

	ap := s.activePiece

	pieceCopy := ap.Clone()
	pieceCopy.rotateAnticlockwise()

	if isMovePossible(s, pieceCopy) {
		s.activePiece = *pieceCopy
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
