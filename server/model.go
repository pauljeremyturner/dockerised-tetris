package server

import (
	"github.com/google/uuid"
	"github.com/pauljeremyturner/dockerised-tetris/shared"
)

type GameState struct {
	Score     int
	Pixels    []shared.Pixel
	NextPiece []shared.Pixel
}

type ServerSession struct {
	player         Player
	moveQueue      chan shared.MoveType
	gameQueue      chan GameState
	gameOverSignal chan bool
	activePiece    Piece
	lines          []Line
	nextPiece      Piece
	gameOver       bool
	score          int
}

type Player struct {
	uuid       uuid.UUID
	playerName string
}
