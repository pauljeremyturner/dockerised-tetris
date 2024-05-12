package shared

import (
	"fmt"
)

type MoveType rune

const (
	MOVELEFT    MoveType = 's'
	MOVERIGHT   MoveType = 'd'
	ROTATELEFT  MoveType = 'a'
	ROTATERIGHT MoveType = 'f'
	DROP        MoveType = 'e'
	DOWN        MoveType = 'x'
	UNKNOWN     MoveType = ' '

	BOARDSIZEX = 16
	BOARDSIZEY = 18
)

type Board struct {
	Height int
	Width  int
}

var DefaultBoard Board = Board{
	Height: BOARDSIZEY,
	Width:  BOARDSIZEX,
}

func (r Board) String() string {
	return fmt.Sprintf("Tetris board. Width: %d,  Height: %d", r.Width, r.Height)
}

func (r MoveType) String() string {
	switch r {
	case MOVELEFT:
		return "Shift Left"
	case MOVERIGHT:
		return "Shift Right"
	case ROTATELEFT:
		return "Rotate Left"
	case ROTATERIGHT:
		return "Rotate Right"
	case DROP:
		return "Drop"
	case DOWN:
		return "Down"
	default:
		return "Unknown"
	}
}
