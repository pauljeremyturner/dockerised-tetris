package client

import (
	"github.com/nsf/termbox-go"
	"log"
	"time"
)

const backgroundColor = termbox.ColorBlack
const boardColor = termbox.ColorBlack
const instructionsColor = termbox.ColorYellow

const originXBoard = 10
const originYBoard = 10

const originXNextPiece = 10
const originYNextPiece = 10

const BLOCK = '▇'
const HORIZONTAL = '═'
const VERTICAL = '║'
const TOPLEFT = '╔'
const TOPRIGHT = '╗'
const BOTTOMLEFT = '╚'
const BOTTOMRIGHT = '╝'

//todo wrap in const {}
const MOVELEFT = 's'
const MOVERIGHT = 'd'
const ROTATELEFT = 'a'
const ROTATERIGHT = 'f'
const DROP = 'e'
const DOWN = 'x'

type TetrisClientUi interface {
	Update(gs GameState)
	NewGame()
}
type KeyListener func(r rune)
type ClientUiState struct {
	eventChannel chan termbox.Event
	appLog       *log.Logger
	keyListener  KeyListener
}

func NewTetrisClientUi(kl KeyListener) TetrisClientUi {
	return ClientUiState{
		eventChannel: make(chan termbox.Event, 1),
		appLog:       GetFileLogger().Logger,
		keyListener:  kl,
	}
}

func (r ClientUiState) NewGame() {

	err := termbox.Init()

	if err != nil {
		log.Panicf("Initialise terminal failed: %s", err)
	}

	termbox.HideCursor()

	go r.readKey()
	go r.listenKeyPress()

	defer termbox.Close()

	drawBorder(0, 0, 40, 20)


	b := Block{
		X:     5,
		Y:     5,
		Color: 5,
	}
	r.drawBoardBlock(b)

	termbox.Flush()
	time.Sleep(5 * time.Second)
	termbox.SetInputMode(termbox.InputEsc)
}

func (r ClientUiState) Update(gm GameState) {

	for _, bl := range gm.Blocks {
		r.drawBoardBlock(bl)
	}
	for _, bl := range gm.NextPiece {
		r.drawNextPieceBlock(bl)
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

	char := rune(event.Ch)

	r.keyListener(char)

	switch char {
	case ROTATELEFT:
		r.appLog.Println("UI ->rotate left")
	case ROTATERIGHT:
		r.appLog.Println("UI ->rotate right")
	case MOVELEFT:
		r.appLog.Println("UI ->move left")
	case MOVERIGHT:
		r.appLog.Println("UI ->move right")
	case DROP:
		r.appLog.Println("UI ->drop")
	case DOWN:
		r.appLog.Println("UI ->down")
	}
}

func (r ClientUiState) drawBoardBlock(block Block) {
	termbox.SetCell(block.X/2, block.Y, BLOCK, termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(block.X/2+1, block.Y, BLOCK, termbox.ColorDefault, termbox.ColorDefault)
}

func (r ClientUiState) drawNextPieceBlock(block Block) {
	termbox.SetCell(block.X/2, block.Y, BLOCK, termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(block.X/2+1, block.Y, BLOCK, termbox.ColorDefault, termbox.ColorDefault)
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
