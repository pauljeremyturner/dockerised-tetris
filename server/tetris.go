package server

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pauljeremyturner/dockerised-tetris/shared"
)

type tetrisState struct {
	newGameChannel chan Player
	activeGames    sync.Map
}

type Tetris interface {
	EnqueueMove(sg ServerSession, move shared.MoveType)
	NewGame(uuid uuid.UUID, playerName string)
}

func (ts *tetrisState) StartNewGame(player Player) ServerSession {

	session := ServerSession{
		player:         player,
		moveQueue:      nil,
		gameQueue:      nil,
		gameOverSignal: nil,
		activePiece:    nil,
		lines:          nil,
		nextPiece:      nil,
		gameOver:       false,
	}

	ts.activeGames.Store(player.uuid, session)

	go redraw(&session)
	go tick(&session)

	return session
}

func (ts *tetrisState) EnqueueMove(u uuid.UUID, move shared.MoveType) {

	ss, ok := ts.activeGames.Load(u)
	if ok {
		serverSession := ss.(ServerSession)
		serverSession.moveQueue <- move
	}
}

func tick(sg *ServerSession) {

	for time := range time.Tick(100 * time.Millisecond) {

		GetFileLogger().Printf("Game tick game: %s time: %s", sg.player.uuid, time)

		sg.moveQueue <- shared.DOWN

		if sg.gameOver {
			break
		}
	}
}

func gameState(ss *ServerSession) GameState {

	allPixels := make([]shared.Pixel, 10)

	allPixels = append(allPixels, ss.activePiece.Pixels()...)

	for _, l := range ss.lines {
		allPixels = append(allPixels, l.Pixels()...)
	}

	return GameState{
		Score:     ss.score,
		Pixels:    allPixels,
		NextPiece: ss.nextPiece.Pixels(),
	}

}

func redraw(sg *ServerSession) {

	for r := range sg.moveQueue {

		switch r {
		case shared.MOVELEFT:
			sg.activePiece.MoveLeft()
		case shared.MOVERIGHT:
			sg.activePiece.MoveRight()
		case shared.ROTATELEFT:
			sg.activePiece.RotateClockwise()
		case shared.ROTATERIGHT:
			sg.activePiece.RotateAnticlockwise()
		case shared.DOWN:
			sg.activePiece.MoveDown()
		}

		sg.gameQueue <- gameState(sg)

	}

}
