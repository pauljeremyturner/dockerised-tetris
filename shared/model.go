package shared

import (
	"log"
)

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

	BOARDSIZEX = 20
	BOARDSIZEY = 20
)

type Board struct {
	Height int
	Width  int
}
