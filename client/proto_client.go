package client

import (
	"context"
	"github.com/google/uuid"
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

func NewTetrisProto(session ClientSession) ProtoClient {

	conn, err := grpc.Dial(address, grpc.WithInsecure())

	if err != nil {

		log.Fatalf("Did not connect: %v", err)
	}

	return ProtoClientState{
		appLog:        GetFileLogger(),
		playerSession: session,
		startGameClient : pf.NewStartGameClient(conn),
		moveClient : pf.NewMoveClient(conn),
	}
}




func (pcs ProtoClientState) ListenToMove() {

	appLog.Println("Listen to Moves")

	for mt := range pcs.playerSession.MoveChannel {

		appLog.Println(string(mt))
		in := &pf.MoveRequest{
			Uuid: pcs.playerSession.Uuid.String(),
			Move: moveTypeToProto(mt),
		}
		pcs.moveClient.Move(context.Background(), in)

	}

	return

}

func (pcs ProtoClientState) ReceiveStream(uuid uuid.UUID, playerName string) {
	request := &pf.NewGameRequest{
		Uuid: uuid.String(),
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
			appLog.Fatalf("%v.StartGame(_) = _, %v", pcs.startGameClient, err)
		}

		gs := &GameState{
			Pixels:    make([]shared.Pixel, len(gameUpdate.Squares), len(gameUpdate.Squares)),
			NextPiece: nil,
			GameOver:  gameUpdate.GameOver,
			Score:     int(gameUpdate.Score),
			Duration:  gameUpdate.Duration,
		}
		for _, sq := range gameUpdate.Squares {
			pixel := shared.Pixel{X: int(sq.X), Y: int(sq.Y), Color: int(sq.Color)}
			gs.Pixels = append(gs.Pixels, pixel)
		}

		pcs.playerSession.BoardUpdateChannel <- *gs
	}

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
