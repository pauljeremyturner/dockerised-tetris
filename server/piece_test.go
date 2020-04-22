package server

import (
	"github.com/pauljeremyturner/dockerised-tetris/shared"
	"gotest.tools/v3/assert"
	"testing"
)

var centre = shared.Pixel{0, 0, 0}

func TestShouldRotateClockWise(t *testing.T) {

	got := shared.Pixel{2, 2, 1}
	got.RotateClockwise(centre)

	want := shared.Pixel{-2, 2, 0}

	assert.Equal(t, got.X, want.X)
	assert.Equal(t, got.Y, want.Y)
}

func TestShouldRotateAnticlockwise(t *testing.T) {

	got := shared.Pixel{2, 2, 1}
	got.RotateAntiClockwise(centre)

	want := shared.Pixel{2, -2, 0}

	assert.Equal(t, got.X, want.X)
	assert.Equal(t, got.Y, want.Y)
}

func TestShouldMoveLeft(t *testing.T) {

	got := shared.Pixel{2, 2, 1}
	got.MoveLeft()

	want := shared.Pixel{1, 2, 0}

	assert.Equal(t, got.X, want.X)
	assert.Equal(t, got.Y, want.Y)
}

func TestShouldMoveRight(t *testing.T) {

	got := shared.Pixel{2, 2, 1}
	got.MoveRight()

	want := shared.Pixel{3, 2, 0}

	assert.Equal(t, got.X, want.X)
	assert.Equal(t, got.Y, want.Y)
}

func TestShouldMoveDown(t *testing.T) {

	got := shared.Pixel{2, 2, 1}
	got.MoveDown()

	want := shared.Pixel{2, 3, 0}

	assert.Equal(t, got.X, want.X)
	assert.Equal(t, got.Y, want.Y)
}
