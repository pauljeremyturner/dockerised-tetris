package server

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pauljeremyturner/dockerised-tetris/shared"
)


type tetris struct {
	activeGames *sync.Map
	board       shared.Board
}

func NewTetris() *tetris {

	return &tetris{
		activeGames: &sync.Map{},
		board: shared.Board{
			Height: shared.BOARDSIZEY,
			Width:  shared.BOARDSIZEX,
		},
	}
}

func (r *tetris) StartNewGame(player Player) *serverSession {

	session := &serverSession{
		player:         player,
		moveQueue:      make(chan shared.MoveType, 10),
		gameQueue:      make(chan GameState, 10),
		gameOverSignal: nil,
		activePiece:    RandomPiece(),
		lines:          Lines{lineMap: make(map[int][]Pixel)},
		nextPiece:      RandomPiece(),
		gameOver:       false,
		board:          r.board,
		startSeconds:   time.Now().Unix(),
	}

	GetFileLogger().Println("New Game", player.uuid.String())
	//r.activeGames.Store(player.uuid.String(), *session)

	go redraw(session)
	go tick(session)

	return session

}

func (r *tetris) EnqueueMove(u uuid.UUID, move shared.MoveType) {

	GetFileLogger().Printf("Queue MoveType: %s to tetris engine player:%s", u.String(), string(move))

	ss, ok := r.activeGames.Load(u)
	if ok {
		serverSession := ss.(*serverSession)
		serverSession.moveQueue <- move
	} else {


		GetFileLogger().Printf("Client Error: uuid for active game not found, uuid: %s", u)
	}
}

func tick(ss *serverSession) {

	for time := range time.Tick(1 * time.Second) {
		if ss.gameOver {
			break
		}
		GetFileLogger().Printf("Tick game for Player: %s at time: %s", ss.player.uuid, time.String())

		ss.moveQueue <- shared.DOWN
	}
}

func gameState(ss *serverSession) GameState {

	allPixels := make([]Pixel, 0)
	allPixels = append(allPixels, ss.activePiece.pixels...)
	for _, lps := range ss.lines.lineMap {
		allPixels = append(allPixels, lps...)
	}
	return GameState{
		Score:     ss.score,
		Pixels:    allPixels,
		NextPiece: ss.nextPiece.pixels,
		GameOver:  ss.gameOver,
		Duration:  time.Now().Unix() - ss.startSeconds,
	}

}

func nextPiece(ss *serverSession) {

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

	ss.lines.Compact(ss.board)

	ss.activePiece = ss.nextPiece
	ss.nextPiece = RandomPiece()
}

func redraw(ss *serverSession) {

	for r := range ss.moveQueue {

		GetFileLogger().Printf("Enqueue game update game for Player: %s", ss.player.uuid)

		GetFileLogger().Println("compare command with commands", r, shared.MOVERIGHT)
		GetFileLogger().Println("compare command with commands", r == shared.MOVERIGHT)

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
		for k, _ := range ss.lines.lineMap {
			if k < topLine {
				topLine = k
			}
		}

		ss.gameOver = topLine <= 2

		ss.gameQueue <- gameState(ss)

		if ss.gameOver {
			close(ss.gameQueue)
		}

	}
}
