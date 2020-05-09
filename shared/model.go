package shared

import (
	"log"
)

type MoveType rune

type Logger struct {
	*log.Logger
}

const (
	MOVELEFT    MoveType = 's'
	MOVERIGHT   MoveType = 'd'
	ROTATELEFT  MoveType = 'a'
	ROTATERIGHT MoveType = 'f'
	DROP        MoveType = 'e'
	DOWN        MoveType = 'x'

	BOARDSIZEX = 20
	BOARDSIZEY = 17
)

type Board struct {
	Height int
	Width  int
}
