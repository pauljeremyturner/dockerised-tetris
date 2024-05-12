package client

import (
	"context"
	"github.com/google/uuid"
	"github.com/nsf/termbox-go"
	"go.uber.org/zap"
	"io"
	"log"

	"github.com/pauljeremyturner/dockerised-tetris/protogen"
	"github.com/pauljeremyturner/dockerised-tetris/shared"
	"google.golang.org/grpc"
)

type ProtoClientState struct {
	playerSession ClientSession
	tetrisClient  protogen.GameClient
	moveClient    protogen.MoveClient
	sugar         *zap.SugaredLogger
}

var (
	address = "localhost:50051"
)

type ProtoClient interface {
	ListenToMove()
	ReceiveStream(uuid uuid.UUID, playerName string)
}

func NewTetrisProto(session *ClientSession, suagr *zap.SugaredLogger) ProtoClient {

	conn, err := grpc.Dial(address, grpc.WithInsecure())

	if err != nil {

		log.Fatalf("Did not connect: %v", err)
	}

	return &ProtoClientState{
		sugar:         suagr,
		playerSession: *session,
		tetrisClient:  protogen.NewGameClient(conn),
		moveClient:    protogen.NewMoveClient(conn),
	}
}

func (pcs ProtoClientState) ListenToMove() {

	for mt := range pcs.playerSession.MoveChannel {

		in := &protogen.MoveRequest{
			Uuid: pcs.playerSession.Uuid.String(),
			Move: moveTypeToProto(mt),
		}
		pcs.moveClient.Move(context.Background(), in)

	}
}

func (pcs ProtoClientState) ReceiveStream(uuid uuid.UUID, playerName string) {
	request := &protogen.NewGameRequest{
		Uuid:       uuid.String(),
		PlayerName: playerName,
	}

	stream, err := pcs.tetrisClient.StartGame(context.TODO(), request)
	if err != nil {

	}
	// Listen to the stream of messages
	for {
		gameUpdate, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
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

func convertColor(ce protogen.Square_ColorEnum) termbox.Attribute {
	var c termbox.Attribute
	switch ce {
	case protogen.Square_MAGENTA:
		c = termbox.ColorMagenta
	case protogen.Square_CYAN:
		c = termbox.ColorCyan
	case protogen.Square_YELLOW:
		c = termbox.ColorYellow
	case protogen.Square_BLUE:
		c = termbox.ColorBlue
	case protogen.Square_RED:
		c = termbox.ColorRed
	case protogen.Square_WHITE:
		c = termbox.ColorWhite
	case protogen.Square_BLACK:
		c = termbox.ColorBlack
	default:
		c = termbox.ColorDefault
	}
	return c
}

func moveTypeToProto(mt shared.MoveType) protogen.MoveRequest_MoveEnum {

	switch mt {
	case shared.DOWN:
		return protogen.MoveRequest_DOWN
	case shared.DROP:
		return protogen.MoveRequest_DROP
	case shared.ROTATERIGHT:
		return protogen.MoveRequest_ROTATERIGHT
	case shared.ROTATELEFT:
		return protogen.MoveRequest_ROTATELEFT
	case shared.MOVERIGHT:
		return protogen.MoveRequest_MOVERIGHT
	case shared.MOVELEFT:
		return protogen.MoveRequest_MOVELEFT
	default:
		panic("fixme")
	}

}
