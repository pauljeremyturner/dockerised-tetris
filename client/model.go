package client

import (
	"github.com/google/uuid"
	"github.com/nsf/termbox-go"
	"github.com/pauljeremyturner/dockerised-tetris/shared"
)

type GameState struct {
	Pixels    []Pixel
	NextPiece []Pixel
	GameOver  bool
	Score     int
	Duration  int64
}

type ClientSession struct {
	Uuid               uuid.UUID
	PlayerName         string
	MoveChannel        chan shared.MoveType
	BoardUpdateChannel chan GameState
}

type Pixel struct {
	X     int
	Y     int
	Color termbox.Attribute
}
