package server

import (
	"github.com/google/uuid"
	"github.com/pauljeremyturner/dockerised-tetris/shared"
	"gotest.tools/v3/assert"
	"testing"
)

func TestPieceAtEndWhenOffBoardY(t *testing.T) {

	p := Piece{
		pixels: []shared.Pixel{shared.Pixel{X: 0, Y: 10, Color: 0}},
	}
	u, _ := uuid.NewRandom()
	ss := ServerSession{
		player:         Player{u, "test"},
		moveQueue:      nil,
		gameQueue:      nil,
		gameOverSignal: nil,
		activePiece:    p,
		lines:          Lines{},
		nextPiece:      Piece{},
		gameOver:       false,
		score:          0,
		board:          shared.Board{10, 10},
	}

	assert.Assert(t, ss.isPieceAtEnd())
}

func TestPieceAtEndWhenBlockedByLineY(t *testing.T) {

	p := Piece{
		pixels: []shared.Pixel{shared.Pixel{X: 4, Y: 9, Color: 0}},
	}
	u, _ := uuid.NewRandom()

	lines := make(map[int][]shared.Pixel)
	lines[10] = []shared.Pixel{shared.Pixel{X: 4, Y: 10, Color: 0}}

	ss := ServerSession{
		player:         Player{u, "test"},
		moveQueue:      nil,
		gameQueue:      nil,
		gameOverSignal: nil,
		activePiece:    p,
		lines:          Lines{},
		nextPiece:      Piece{},
		gameOver:       false,
		score:          0,
		board:          shared.Board{10, 10},
	}

	assert.Assert(t, ss.isPieceAtEnd())
}

func TestMovePossibleWhenNotBlockedByLineY(t *testing.T) {

	p := Piece{
		pixels: []shared.Pixel{shared.Pixel{X: 4, Y: 8, Color: 0}},
	}
	u, _ := uuid.NewRandom()

	lines := make(map[int][]shared.Pixel)
	lines[10] = []shared.Pixel{shared.Pixel{X: 4, Y: 10, Color: 0}}

	ss := ServerSession{
		player:         Player{u, "test"},
		moveQueue:      nil,
		gameQueue:      nil,
		gameOverSignal: nil,
		activePiece:    p,
		lines:          Lines{},
		nextPiece:      Piece{},
		gameOver:       false,
		score:          0,
		board:          shared.Board{10, 10},
	}

	assert.Assert(t, !ss.isPieceAtEnd())
}
