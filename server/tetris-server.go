package server

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pauljeremyturner/dockerised-tetris/shared"
)

type tetrisState struct {
	activeGames sync.Map
	board       shared.Board
}

type Tetris interface {
	EnqueueMove(u uuid.UUID, move shared.MoveType)
	StartNewGame(player Player) *ServerSession
}

func NewTetris() Tetris {

	return &tetrisState{
		activeGames: sync.Map{},
		board: shared.Board{
			Height: shared.BOARDSIZEY,
			Width:  shared.BOARDSIZEX,
		},
	}
}

func (ts *tetrisState) StartNewGame(player Player) *ServerSession {

	session := ServerSession{
		player:         player,
		moveQueue:      make(chan shared.MoveType, 10),
		gameQueue:      make(chan GameState, 10),
		gameOverSignal: nil,
		activePiece:    RandomPiece(),
		lines:          Lines{lineMap: make(map[int][]Pixel)},
		nextPiece:      RandomPiece(),
		gameOver:       false,
		board:          ts.board,
	}

	ts.activeGames.Store(player.uuid, session)

	go redraw(&session)
	go tick(&session)

	return &session

}

func (ts *tetrisState) EnqueueMove(u uuid.UUID, move shared.MoveType) {

	GetFileLogger().Printf("Queue MoveType: %s to tetris engine player:%s", u.String(), string(move))

	ss, ok := ts.activeGames.Load(u)
	if ok {
		serverSession := ss.(ServerSession)
		serverSession.moveQueue <- move
	}
}

func tick(ss *ServerSession) {

	for time := range time.Tick(300 * time.Millisecond) {
		if ss.gameOver {
			break
		}
		GetFileLogger().Printf("Tick game for Player: %s at time: %s", ss.player.uuid, time.String())

		ss.moveQueue <- shared.DOWN

	}
}

func gameState(ss *ServerSession) GameState {

	allPixels := make([]Pixel, 0)

	allPixels = append(allPixels, ss.activePiece.pixels...)

	for _, pslice := range ss.lines.lineMap {
		allPixels = append(allPixels, pslice...)
	}

	return GameState{
		Score:     ss.score,
		Pixels:    allPixels,
		NextPiece: ss.nextPiece.pixels,
	}

}

func nextPiece(ss *ServerSession) {

	lineMap := ss.lines.lineMap
	for _, pp := range ss.activePiece.pixels {

		if line, ok := lineMap[pp.Y]; ok {
			line = append(line, pp)
			lineMap[pp.Y] = line
		} else {
			newLine := make([]Pixel, 0)
			newLine = append(newLine, pp)
			lineMap[pp.Y] = newLine
		}
	}

	ss.activePiece = ss.nextPiece
	ss.nextPiece = RandomPiece()
}

func redraw(ss *ServerSession) {

	for r := range ss.moveQueue {

		GetFileLogger().Printf("Enqueue game update game for Player: %s", ss.player.uuid)

		GetFileLogger().Println("compare command with commands", r, shared.ROTATERIGHT)
		GetFileLogger().Println("compare command with commands", r == shared.ROTATERIGHT)

		GetFileLogger().Printf("Active Piece Before Move piece: %s move: %s", ss.activePiece.String(), string(r))

		var pieceMoved bool
		switch r {
		case shared.MOVELEFT:
			pieceMoved = ss.MoveActivePieceLeftIfPossible()
		case shared.MOVERIGHT:
			pieceMoved = ss.MoveActivePieceRightIfPossible()
		case shared.ROTATELEFT:
			pieceMoved = ss.RotateActivePieceAnticlockwiseIfPossible()
		case shared.ROTATERIGHT:
			pieceMoved = ss.RotateActivePieceClockwiseIfPossible()
		case shared.DOWN:
			if ss.MoveActivePieceDownIfPossible() {
				pieceMoved = true
			} else {
				nextPiece(ss)
			}
		case shared.DROP:
			for b := true; b; b = ss.MoveActivePieceDownIfPossible() {
			}
			nextPiece(ss)
		}

		GetFileLogger().Printf("Active Piece After Move piece: %s, piece moved?: %t", ss.activePiece.String(), pieceMoved)

		var topLine = ss.board.Height
		for k := range ss.lines.lineMap {
			if k < topLine {
				topLine = k
			}
		}
		ss.gameOver = topLine <= 2

		ss.gameQueue <- gameState(ss)

	}
}
