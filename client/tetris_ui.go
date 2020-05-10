package client

import (
	"fmt"
	"time"

	"github.com/nsf/termbox-go"
	"github.com/pauljeremyturner/dockerised-tetris/shared"
)

const (

	originXBoard = 2
	originYBoard = 10

	originXNextPiece = 15
	originYNextPiece = 6

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
	board         shared.Board
}

func NewTetrisUi(cs *ClientSession) TetrisUi {
	return TetrisUi{
		eventChannel:  make(chan termbox.Event, 10),
		appLog:        GetFileLogger(),
		playerSession: *cs,
		board: shared.Board{
			Height: shared.BOARDSIZEY,
			Width:  shared.BOARDSIZEX,
		},
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

	drawBorder(0, 0, 35, 30)
	r.writeMessage("tetris://", 2, 2, termbox.ColorWhite)

	termbox.Flush()
	time.Sleep(5 * time.Minute)
	termbox.SetInputMode(termbox.InputEsc)
}

func (r TetrisUi) ListenToBoardUpdates() {

	for gm := range r.playerSession.BoardUpdateChannel {

		r.appLog.Printf("Board Update: %s", gm.String())

		if gm.GameOver {
			r.writeMessage("GAME OVER", 2, 5, termbox.ColorWhite)
			termbox.Flush()
			break
		}

		for x := 0; x < r.board.Width; x++ {
			for y := 0; y < r.board.Height; y++ {
				r.drawBoardPixel(Pixel{X:x, Y:y, Color: termbox.ColorDefault})
			}
		}

		for x := 0; x < 4; x++ {
			for y := 0; y < 4; y++ {
				r.drawNextPiecePixel(Pixel{X:x, Y:y, Color: termbox.ColorDefault})
			}
		}

		r.writeMessage(fmt.Sprintf("player: %s", r.playerSession.PlayerName), 2, 3, termbox.ColorWhite)
		r.writeMessage(fmt.Sprintf("pieces: %d", gm.Pieces), 2, 4, termbox.ColorWhite)
		r.writeMessage(fmt.Sprintf("lines: %d", gm.Lines), 2, 5, termbox.ColorWhite)
		r.writeMessage(fmt.Sprintf("next piece:"), 2, 6, termbox.ColorWhite)

		for _, p := range gm.Pixels {
			r.drawBoardPixel(p)
		}
		for _, p := range gm.NextPiece {
			r.drawNextPiecePixel(p)
		}
		termbox.Flush()
	}
}

func (r TetrisUi) String() string {
	return fmt.Sprintf("Tetris client: playerSession: %s, board: %s", r.playerSession.String(), r.board.String())
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


func (r TetrisUi) drawBoardPixel(p Pixel) {

	//r.appLog.Println("draw board pixel ", p)
	termbox.SetCell(originXBoard+(2*p.X), originYBoard+p.Y, ' ', p.Color, p.Color)
	termbox.SetCell(originXBoard+(2*p.X+1), originYBoard+p.Y, ' ', p.Color, p.Color)
}

func (r TetrisUi) writeMessage(message string, x int, y int, color termbox.Attribute) {

	for _, char := range message {
		termbox.SetCell(x, y, char, termbox.ColorBlack, color)
		x++
	}

}

func (r TetrisUi) drawNextPiecePixel(p Pixel) {
	r.appLog.Println("draw next piece pixel ", p)


	termbox.SetCell(originXNextPiece+(2*p.X), originYNextPiece+p.Y, ' ', p.Color, p.Color)
	termbox.SetCell(originXNextPiece+(2*p.X+1), originYNextPiece+p.Y, ' ', p.Color, p.Color)
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
