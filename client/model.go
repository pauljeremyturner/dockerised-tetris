package client

import (
	"github.com/google/uuid"
	"github.com/pauljeremyturner/dockerised-tetris/shared"
)

type GameState struct {
	Pixels    []shared.Pixel
	NextPiece []shared.Pixel
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
