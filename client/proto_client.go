package client

import (
	"context"
	"github.com/google/uuid"
	"github.com/nsf/termbox-go"
	"io"
	"log"

	pf "github.com/pauljeremyturner/dockerised-tetris/protofiles"
	"github.com/pauljeremyturner/dockerised-tetris/shared"
	"google.golang.org/grpc"
)

type ProtoClientState struct {
	appLog          *Logger
	playerSession   ClientSession
	startGameClient pf.StartGameClient
	moveClient      pf.MoveClient
}

var (
	address = "localhost:50051"
)

type ProtoClient interface {
	ListenToMove()
	ReceiveStream(uuid uuid.UUID, playerName string)
}

func NewTetrisProto(session *ClientSession) ProtoClient {

	conn, err := grpc.Dial(address, grpc.WithInsecure())

	if err != nil {

		log.Fatalf("Did not connect: %v", err)
	}

	return &ProtoClientState{
		appLog:          GetFileLogger(),
		playerSession:   *session,
		startGameClient: pf.NewStartGameClient(conn),
		moveClient:      pf.NewMoveClient(conn),
	}
}

func (pcs ProtoClientState) ListenToMove() {

	appLog.Println("Listen to Moves")

	for mt := range pcs.playerSession.MoveChannel {

		in := &pf.MoveRequest{
			Uuid: pcs.playerSession.Uuid.String(),
			Move: moveTypeToProto(mt),
		}
		appLog.Println("Sending Move, ", in)
		pcs.moveClient.Move(context.Background(), in)

	}
}

func (pcs ProtoClientState) ReceiveStream(uuid uuid.UUID, playerName string) {
	request := &pf.NewGameRequest{
		Uuid:       uuid.String(),
		PlayerName: playerName,
	}
	stream, err := pcs.startGameClient.StartGame(context.Background(), request)
	if err != nil {
		appLog.Fatalf("%v.StartGame(_) = _, %v", pcs.startGameClient, err)
	}
	// Listen to the stream of messages
	for {
		gameUpdate, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			appLog.Println("StartGame state, error", pcs.startGameClient, err)
			break
		}

		gs := &GameState{
			Pixels:    make([]Pixel, 0),
			NextPiece: make([]Pixel, 0),
			GameOver:  gameUpdate.GameOver,
			Lines:     int(gameUpdate.Lines),
			Pieces:    int(gameUpdate.Pieces),
			Duration:  gameUpdate.Duration,
		}

		appLog.Println("Game over?", gs.GameOver)

		for _, sq := range gameUpdate.Squares {
			pixel := Pixel{X: int(sq.X), Y: int(sq.Y), Color: convertColor(sq.Color)}
			gs.Pixels = append(gs.Pixels, pixel)
		}

		for _, sq := range gameUpdate.NextPiece {
			pixel := Pixel{X: int(sq.X), Y: int(sq.Y), Color: convertColor(sq.Color)}
			gs.NextPiece = append(gs.NextPiece, pixel)
		}

		pcs.playerSession.BoardUpdateChannel <- *gs
	}

}

func convertColor(ce pf.Square_ColorEnum) termbox.Attribute {
	var c termbox.Attribute
	switch ce {
	case pf.Square_MAGENTA:
		c = termbox.ColorMagenta
	case pf.Square_CYAN:
		c = termbox.ColorCyan
	case pf.Square_YELLOW:
		c = termbox.ColorYellow
	case pf.Square_BLUE:
		c = termbox.ColorBlue
	case pf.Square_RED:
		c = termbox.ColorRed
	case pf.Square_WHITE:
		c = termbox.ColorWhite
	case pf.Square_BLACK:
		c = termbox.ColorBlack
	default:
		c = termbox.ColorDefault
	}
	return c
}

func moveTypeToProto(mt shared.MoveType) pf.MoveRequest_MoveEnum {

	switch mt {
	case shared.DOWN:
		return pf.MoveRequest_DOWN
	case shared.DROP:
		return pf.MoveRequest_DROP
	case shared.ROTATERIGHT:
		return pf.MoveRequest_ROTATERIGHT
	case shared.ROTATELEFT:
		return pf.MoveRequest_ROTATELEFT
	case shared.MOVERIGHT:
		return pf.MoveRequest_MOVERIGHT
	case shared.MOVELEFT:
		return pf.MoveRequest_MOVELEFT
	default:
		panic("fixme")
	}

}
