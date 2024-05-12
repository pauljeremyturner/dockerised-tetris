package server

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pauljeremyturner/dockerised-tetris/shared"
	"go.uber.org/zap"
	"sync"
	"time"
)

var nextPieceFunc func() Piece

func init() {
	nextPieceFunc = RandomPiece
}

type tetris struct {
	activeGames *sync.Map
	sugar       *zap.SugaredLogger
}

type Tetris interface {
	StartNewGame(player Player) ServerSession
	GetGame(u uuid.UUID) ServerSession
	EnqueueMove(u uuid.UUID, move shared.MoveType) error
}

func NewTetris(sugar *zap.SugaredLogger) *tetris {
	return &tetris{
		activeGames: &sync.Map{},
		sugar:       sugar,
	}
}

func (r *tetris) GetGame(u uuid.UUID) (ServerSession, error) {
	ss, ok := r.activeGames.Load(u)
	if ok {
		serverSession := ss.(ServerSession)
		return serverSession, nil
	}

	return nil, NewErrorInvalidRequest("game not found", nil)
}

func (r *tetris) StartNewGame(player Player) ServerSession {

	s := NewServerSession(player, r)

	r.sugar.Infof("New Game, name :%s, UUID: %s", player.Name, player.UUID.String())
	r.activeGames.Store(player.UUID.String(), s)

	go s.Redraw()
	go s.Tick()

	return s

}

func (r *tetris) EnqueueMove(u uuid.UUID, move shared.MoveType) error {

	r.sugar.Infof("Move.  player UUID: %s, move: %s", u, move.String())

	ss, err := r.GetGame(u)
	if err != nil {
		return fmt.Errorf("could not find game, cause: %w", err)
	}
	ss.Move(move)

	return nil

}

type serverSession struct {
	tetris       *tetris
	player       Player
	moveQueue    chan shared.MoveType
	gameQueue    chan GameState
	activePiece  Piece
	lines        Lines
	nextPiece    Piece
	gameOver     bool
	lineCount    int
	pieceCount   int
	board        shared.Board
	startSeconds int64
	sugar        *zap.SugaredLogger
}

func NewServerSession(player Player, tetris *tetris) ServerSession {
	return &serverSession{
		activePiece:  RandomPiece(),
		board:        shared.Board{Width: shared.BOARDSIZEX, Height: shared.BOARDSIZEY},
		gameOver:     false,
		gameQueue:    make(chan GameState, 10),
		lines:        Lines{lineMap: make(map[int][]Pixel)},
		moveQueue:    make(chan shared.MoveType, 10),
		nextPiece:    nextPieceFunc(),
		player:       player,
		startSeconds: time.Now().Unix(),
		sugar:        tetris.sugar,
		tetris:       tetris,
	}
}

type CompletedMovePublisher interface {
	PublishCompletedMove(ctx context.Context, gs GameState) error
}

type ServerSession interface {
	GetGameState() GameState
	GetPlayer() Player
	IsGameOver() bool
	Move(move shared.MoveType)
	PublishCompletedMoves(ctx context.Context, p CompletedMovePublisher) error
	Redraw()
	Tick()
}

func (s *serverSession) GetPlayer() Player {
	return s.player
}

func (s *serverSession) PublishCompletedMoves(ctx context.Context, pub CompletedMovePublisher) error {
	for gs := range s.gameQueue {
		err := pub.PublishCompletedMove(ctx, gs)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *serverSession) IsGameOver() bool {
	return s.gameOver
}
func (s *serverSession) Move(move shared.MoveType) {
	s.moveQueue <- move
}
func (s *serverSession) Tick() {

	for time := range time.Tick(1 * time.Second) {
		if s.gameOver {
			break
		}
		s.tetris.sugar.Infof("Tick game for Player: %s at time: %s", s.player.UUID, time.String())

		s.moveQueue <- shared.DOWN
	}
}

func (s *serverSession) GetGameState() GameState {

	allPixels := make([]Pixel, 0)
	allPixels = append(allPixels, s.activePiece.pixels...)
	for _, lps := range s.lines.lineMap {
		allPixels = append(allPixels, lps...)
	}
	return GameState{
		LineCount:  s.lineCount,
		Pixels:     allPixels,
		NextPiece:  s.nextPiece.pixels,
		GameOver:   s.gameOver,
		PieceCount: s.pieceCount,
		Duration:   time.Now().Unix() - s.startSeconds,
	}

}

func (s *serverSession) ChooseNextPiece() {

	lineMap := s.lines.lineMap
	for _, pp := range s.activePiece.pixels {

		if line, ok := lineMap[pp.Y]; ok {
			line = append(line, pp)
			lineMap[pp.Y] = line
		} else {
			newLine := make([]Pixel, 0)
			newLine = append(newLine, pp)
			lineMap[pp.Y] = newLine
		}
	}

	s.pieceCount = s.pieceCount + 1
	compactedLines := s.lines.Compact(s.board)
	s.lineCount = s.lineCount + compactedLines
	s.activePiece = s.nextPiece
	s.nextPiece = RandomPiece()
}

func (s *serverSession) Redraw() {

	for r := range s.moveQueue {

		s.sugar.Debugf("Enqueue game update game for Player: %s", s.player.UUID)
		s.sugar.Debugf("Active Piece Before Move piece: %s move: %s", s.activePiece.String(), string(r))

		var pieceMoved bool
		switch r {
		case shared.MOVELEFT:
			pieceMoved = s.MoveActivePieceLeftIfPossible()
		case shared.MOVERIGHT:
			pieceMoved = s.MoveActivePieceRightIfPossible()
		case shared.ROTATELEFT:
			pieceMoved = s.RotateActivePieceAnticlockwiseIfPossible()
		case shared.ROTATERIGHT:
			pieceMoved = s.RotateActivePieceClockwiseIfPossible()
		case shared.DOWN:
			if s.MoveActivePieceDownIfPossible() {
				pieceMoved = true
			} else {
				s.ChooseNextPiece()
			}
		case shared.DROP:
			for b := true; b; b = s.MoveActivePieceDownIfPossible() {
			}
			s.ChooseNextPiece()
		}

		s.tetris.sugar.Debugf("Active Piece After Move piece: %s, piece moved?: %t", s.activePiece.String(), pieceMoved)

		var topLine = s.board.Height
		for k, _ := range s.lines.lineMap {
			if k < topLine {
				topLine = k
			}
		}

		s.gameOver = topLine <= 2

		s.gameQueue <- s.GetGameState()

		if s.gameOver {
			close(s.gameQueue)
		}

	}
}
