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

	originXBoard = 2
	originYBoard = 2

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

type TetrisUi struct {
	eventChannel  chan termbox.Event
	playerSession ClientSession
	appLog        *Logger
}

func NewTetrisUi(cs *ClientSession) TetrisUi {
	return TetrisUi{
		eventChannel:  make(chan termbox.Event, 10),
		appLog:        GetFileLogger(),
		playerSession: *cs,
	}
}

func (r TetrisUi) StartGame() {

	err := termbox.Init()

	if err != nil {
		appLog.Panicf("Initialise terminal failed: %s", err)
	}

	go r.readKey()
	go r.listenKeyPress()
	go r.ListenToBoardUpdates()

	defer termbox.Close()

	drawBorder(0, 0, 40, 20)

	termbox.Flush()
	time.Sleep(5 * time.Minute)
	termbox.SetInputMode(termbox.InputEsc)
}

func (r TetrisUi) ListenToBoardUpdates() {

	for gm := range r.playerSession.BoardUpdateChannel {

		r.appLog.Printf("Board Update: %s", gm.String())

		if gm.GameOver {
			r.writeMessage("GAME OVER", 0, 5, termbox.ColorWhite)
			termbox.Flush()
			break
		}

		for x := 0; x < 50; x++ {
			for y := 1; y < 30; y++ {
				r.clearBoardPixel(Pixel{x, y, 0})
			}
		}

		for _, p := range gm.Pixels {
			r.drawBoardPixel(p)
		}
		for _, p := range gm.NextPiece {
			r.drawNextPieceBlock(p)
		}
		termbox.Flush()
	}
}

func (r TetrisUi) String() string {
	return "tetris" //do more here!
}

func (r TetrisUi) readKey() {
	switch ev := termbox.PollEvent(); ev.Type {
	case termbox.EventKey:
		r.eventChannel <- ev
	}
}

func (r TetrisUi) listenKeyPress() {
	for e := range r.eventChannel {
		r.onKeyPress(e)
		r.readKey()
	}
}

func (r TetrisUi) onKeyPress(event termbox.Event) {

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

func (r TetrisUi) clearBoardPixel(p Pixel) {

	//r.appLog.Println("draw board pixel ", p)

	termbox.SetCell(originXBoard+(2*p.X), originYBoard+p.Y, ' ', termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(originXBoard+(2*p.X+1), originYBoard+p.Y, ' ', termbox.ColorDefault, termbox.ColorDefault)
}
func (r TetrisUi) drawBoardPixel(p Pixel) {

	//r.appLog.Println("draw board pixel ", p)
	var c termbox.Attribute
	switch p.Color {
	case 1:
		c = termbox.ColorMagenta
	case 2:
		c = termbox.ColorRed
	case 3:
		c = termbox.ColorGreen
	case 4:
		c = termbox.ColorCyan
	case 5:
		c = termbox.ColorWhite
	case 6:
		c = termbox.ColorYellow
	case 7:
		c = termbox.ColorBlue
	default:
		c = termbox.ColorDefault
	}

	termbox.SetCell(originXBoard+(2*p.X), originYBoard+p.Y, ' ', termbox.ColorDefault, c)
	termbox.SetCell(originXBoard+(2*p.X+1), originYBoard+p.Y, ' ', termbox.ColorDefault, c)
}

func (r TetrisUi) writeMessage(message string, x int, y int, color termbox.Attribute) {

	for _, char := range message {
		termbox.SetCell(x, y, char, termbox.ColorDefault, color)
		x++
	}

}

func (r TetrisUi) drawNextPieceBlock(pixel Pixel) {
	termbox.SetCell(pixel.X/2, pixel.Y, BLOCK, termbox.ColorDefault, termbox.ColorDefault)
	//termbox.SetCell(pixel.X/2+1, pixel.Y, BLOCK, termbox.ColorDefault, termbox.ColorDefault)
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
