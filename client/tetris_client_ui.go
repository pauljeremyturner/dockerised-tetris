package client

import (
	"time"

	"github.com/nsf/termbox-go"
	"github.com/pauljeremyturner/dockerised-tetris/shared"
)

const (
	backgroundColor   = termbox.ColorBlack
	boardColor        = termbox.ColorBlack
	instructionsColor = termbox.ColorYellow

	originXBoard = 10
	originYBoard = 10

	originXNextPiece = 10
	originYNextPiece = 10

	BLOCK       = '▇'
	HORIZONTAL  = '═'
	VERTICAL    = '║'
	TOPLEFT     = '╔'
	TOPRIGHT    = '╗'
	BOTTOMLEFT  = '╚'
	BOTTOMRIGHT = '╝'
)

type TetrisClientUi interface {
	ListenToBoardUpdates()
	StartGame()
}
type ClientUiState struct {
	eventChannel  chan termbox.Event
	playerSession ClientSession
	appLog        *Logger
}

func NewTetrisClientUi(ps ClientSession) TetrisClientUi {
	return ClientUiState{
		eventChannel:  make(chan termbox.Event, 1),
		appLog:        GetFileLogger(),
		playerSession: ps,
	}
}

func (r ClientUiState) StartGame() {

	err := termbox.Init()

	if err != nil {
		appLog.Panicf("Initialise terminal failed: %s", err)
	}

	termbox.HideCursor()

	go r.readKey()
	go r.listenKeyPress()

	defer termbox.Close()

	drawBorder(0, 0, 40, 20)

	b := shared.Pixel{
		X:     5,
		Y:     5,
		Color: 5,
	}
	r.drawBoardPixel(b)

	termbox.Flush()
	time.Sleep(5 * time.Minute)
	termbox.SetInputMode(termbox.InputEsc)
}

func (r ClientUiState) ListenToBoardUpdates() {

	for gm := range r.playerSession.BoardUpdateChannel {

		r.appLog.Println("Board Update", gm)

		for _, p := range gm.Pixels {
			r.drawBoardPixel(p)
		}
		for _, p := range gm.NextPiece {
			r.drawNextPieceBlock(p)
		}
		termbox.Flush()
	}
}

func (r ClientUiState) String() string {
	return "tetris" //do more here!
}

func (r ClientUiState) readKey() {
	switch ev := termbox.PollEvent(); ev.Type {
	case termbox.EventKey:
		r.eventChannel <- ev
	}
}

func (r ClientUiState) listenKeyPress() {
	for e := range r.eventChannel {
		r.onKeyPress(e)
		r.readKey()
	}
}

func (r ClientUiState) onKeyPress(event termbox.Event) {

	moveType := shared.MoveType(event.Ch)

	switch moveType {
	case shared.ROTATELEFT:
		fallthrough
	case shared.ROTATERIGHT:
		fallthrough
	case shared.MOVELEFT:
		fallthrough
	case shared.MOVERIGHT:
		fallthrough
	case shared.DROP:
		fallthrough
	case shared.DOWN:
		r.appLog.Printf("UI -> enqueue move %s", string(moveType))
		r.playerSession.MoveChannel <- moveType
	default:
		r.appLog.Println("UI -> unknown comamnd, ignoring")
	}
}

func (r ClientUiState) drawBoardPixel(p shared.Pixel) {

	r.appLog.Println("draw board pixel ", p)

	termbox.SetCell(p.X/2, p.Y, BLOCK, termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(p.X/2+1, p.Y, BLOCK, termbox.ColorDefault, termbox.ColorDefault)
}

func (r ClientUiState) drawNextPieceBlock(pixel shared.Pixel) {
	termbox.SetCell(pixel.X/2, pixel.Y, BLOCK, termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(pixel.X/2+1, pixel.Y, BLOCK, termbox.ColorDefault, termbox.ColorDefault)
}

func drawBorder(leftEdge int, topEdge int, width int, height int) {
	termbox.SetCell(leftEdge, topEdge, TOPLEFT, termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(leftEdge+width, topEdge, TOPRIGHT, termbox.ColorDefault, termbox.ColorDefault)

	for x := leftEdge + 1; x < width; x++ {
		termbox.SetCell(x, topEdge, HORIZONTAL, termbox.ColorDefault, termbox.ColorDefault)
		termbox.SetCell(x, topEdge+height, HORIZONTAL, termbox.ColorDefault, termbox.ColorDefault)
	}
	for y := topEdge + 1; y < height; y++ {
		termbox.SetCell(leftEdge, y, VERTICAL, termbox.ColorDefault, termbox.ColorDefault)
		termbox.SetCell(leftEdge+width, y, VERTICAL, termbox.ColorDefault, termbox.ColorDefault)
	}
	termbox.SetCell(leftEdge, topEdge+height, BOTTOMLEFT, termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(leftEdge+width, topEdge+height, BOTTOMRIGHT, termbox.ColorDefault, termbox.ColorDefault)
}
