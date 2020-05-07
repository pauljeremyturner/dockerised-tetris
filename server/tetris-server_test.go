package server

import (
	"github.com/google/uuid"
	"github.com/pauljeremyturner/dockerised-tetris/shared"
	"gotest.tools/v3/assert"
	"testing"
	"time"
)

func TestShouldMoveDownPeriodically(t *testing.T) {
	u, _ := uuid.NewRandom()
	mq := make(chan shared.MoveType, 5)
	ss := &ServerSession{
		player:         Player{u, "player"},
		moveQueue:      mq,
		gameQueue:      nil,
		gameOverSignal: nil,
		activePiece:    NewI(),
		lines:          Lines{},
		nextPiece:      NewI(),
		gameOver:       false,
		score:          0,
	}

	go tick(ss)

	time.Sleep(3 * time.Second)

	ss.gameOver = true

	time.Sleep(1 * time.Second)
	assert.Assert(t, len(mq) >= 2)

	close(mq)

	for got := range mq {
		assert.Assert(t, got == shared.DOWN)
	}

}
