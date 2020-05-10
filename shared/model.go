package shared

import (
	"fmt"
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

	BOARDSIZEX = 16
	BOARDSIZEY = 18
)

type Board struct {
	Height int
	Width  int
}

func (r Board) String() string {
	return fmt.Sprintf("Tetris board. Width: %d,  Height: %d", r.Width, r.Height)
}
