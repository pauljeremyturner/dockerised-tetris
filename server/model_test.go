package server

import (
	"github.com/google/uuid"
	"github.com/pauljeremyturner/dockerised-tetris/shared"
	"gotest.tools/v3/assert"
	"testing"
)

var centre = Pixel{0, 0, 0}

func TestShouldRotateClockWise(t *testing.T) {

	got := Pixel{2, 2, 1}
	got.rotateClockwise(centre)

	want := Pixel{-2, 2, 0}

	assert.Equal(t, got.X, want.X)
	assert.Equal(t, got.Y, want.Y)
}

func TestShouldRotateAnticlockwise(t *testing.T) {

	got := Pixel{2, 2, 1}
	got.rotateAntiClockwise(centre)

	want := Pixel{2, -2, 0}

	assert.Equal(t, got.X, want.X)
	assert.Equal(t, got.Y, want.Y)
}

func TestShouldMoveLeft(t *testing.T) {

	got := Pixel{2, 2, 1}
	got.moveLeft()

	want := Pixel{1, 2, 0}

	assert.Equal(t, got.X, want.X)
	assert.Equal(t, got.Y, want.Y)
}

func TestShouldMoveRight(t *testing.T) {

	got := Pixel{2, 2, 1}
	got.moveRight()

	want := Pixel{3, 2, 0}

	assert.Equal(t, got.X, want.X)
	assert.Equal(t, got.Y, want.Y)
}

func TestShouldMoveDown(t *testing.T) {

	got := Pixel{2, 2, 1}
	got.moveDown()

	want := Pixel{2, 3, 0}

	assert.Equal(t, got.X, want.X)
	assert.Equal(t, got.Y, want.Y)
}

func TestPieceAtEndWhenOffBoardY(t *testing.T) {

	p := Piece{
		pixels: []Pixel{Pixel{X: 0, Y: 10, Color: 0}},
	}
	u, _ := uuid.NewRandom()
	ss := serverSession{
		player:         Player{u, "test"},
		moveQueue:      nil,
		gameQueue:      nil,
		activePiece:    p,
		lines:          Lines{},
		nextPiece:      Piece{},
		gameOver:       false,
		board:          shared.Board{10, 10},
	}

	assert.Assert(t, !ss.MoveActivePieceDownIfPossible())
}

func TestPieceAtEndWhenBlockedByLineY(t *testing.T) {

	p := Piece{
		pixels: []Pixel{Pixel{X: 4, Y: 9, Color: 0}},
	}
	u, _ := uuid.NewRandom()

	lines := make(map[int][]Pixel)
	lines[10] = []Pixel{Pixel{X: 4, Y: 10, Color: 0}}

	ss := serverSession{
		player:         Player{u, "test"},
		moveQueue:      nil,
		gameQueue:      nil,
		activePiece:    p,
		lines:          Lines{},
		nextPiece:      Piece{},
		gameOver:       false,
		board:          shared.Board{10, 10},
	}

	assert.Assert(t, !ss.MoveActivePieceDownIfPossible())
}

func TestMovePossibleWhenNotBlockedByLineY(t *testing.T) {

	p := Piece{
		pixels: []Pixel{Pixel{X: 4, Y: 8, Color: 0}},
	}
	u, _ := uuid.NewRandom()

	lines := make(map[int][]Pixel)
	lines[10] = []Pixel{Pixel{X: 4, Y: 10, Color: 0}}

	ss := serverSession{
		player:         Player{u, "test"},
		moveQueue:      nil,
		gameQueue:      nil,
		activePiece:    p,
		lines:          Lines{},
		nextPiece:      Piece{},
		gameOver:       false,
		board:          shared.Board{10, 10},
	}

	assert.Assert(t, ss.MoveActivePieceDownIfPossible())
}
