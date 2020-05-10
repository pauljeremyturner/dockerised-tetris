package client

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/nsf/termbox-go"
	"github.com/pauljeremyturner/dockerised-tetris/shared"
)

type GameState struct {
	Pixels    []Pixel
	NextPiece []Pixel
	GameOver  bool
	Lines     int
	Pieces    int
	Duration  int64
}

type Pixel struct {
	X       int
	Y       int
	Color termbox.Attribute
}

func (r *GameState) String() string {
	var s = "Pixels:"
	for _, pix := range r.Pixels {
		s = s + pix.String()
	}
	s = s + "\nNext Piece:"
	for _, pix := range r.NextPiece {
		s = s + pix.String()
	}
	s = s + "\nPiece Count:" + string(r.Pieces)
	s = s + "\nLine Count Count:" + string(r.Lines)
	return s
}

func (r *Pixel) String() string {
	return fmt.Sprintf("pixel, (%d, %d) color %d; ", r.X, r.Y, r.Color)
}

type ClientSession struct {
	Uuid               uuid.UUID
	PlayerName         string
	MoveChannel        chan shared.MoveType
	BoardUpdateChannel chan GameState
}

func (r ClientSession) String() string {
	return fmt.Sprintf("ClientSession uuid: %s, PlayerName: %s", r.Uuid.String(), r.PlayerName)
}
