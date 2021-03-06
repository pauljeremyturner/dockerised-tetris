package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/pauljeremyturner/dockerised-tetris/shared"

	"github.com/google/uuid"
	pf "github.com/pauljeremyturner/dockerised-tetris/protofiles"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var port = ":50051"

type protoServer struct {
	tetris tetris
}

type ProtoServer interface {
	CreateNewGame() uuid.UUID
	Move(uuid uuid.UUID)
}

func StartProtoServer(t *tetris) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pf.RegisterMoveServer(s, &protoServer{tetris: *t})
	pf.RegisterStartGameServer(s, &protoServer{tetris: *t})
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func (s *protoServer) StartGame(in *pf.NewGameRequest, stream pf.StartGame_StartGameServer) error {

	u, _ := uuid.Parse(in.Uuid)
	ss := s.tetris.StartNewGame(Player{uuid: u, playerName: in.PlayerName})
	s.tetris.activeGames.Store(u, ss)

	for gs := range ss.gameQueue {

		gameUpdateResponse := &pf.GameUpdateResponse{
			Uuid:       in.Uuid,
			PlayerName: in.PlayerName,
			GameOver:   gs.GameOver,
			Lines:      uint32(gs.LineCount),
			Pieces:     uint32(gs.PieceCount),
			Duration:   gs.Duration,
			Squares:    pixelsToSquares(gs.Pixels),
			NextPiece:  pixelsToSquares(gs.NextPiece),
		}

		GetFileLogger().Println("Receive board update, sending", gameUpdateResponse.String())


		if err := stream.Send(gameUpdateResponse); err != nil {
			fmt.Printf("something went wrong %s", err)
		}

	}

	GetFileLogger().Println("Game update channel closed, game over")

	return nil
}

func (s *protoServer) Move(ctx context.Context, in *pf.MoveRequest) (*pf.MoveResponse, error) {
	u, _ := uuid.Parse(in.Uuid)

	s.tetris.EnqueueMove(u, protoMoveTypeToModel(in.GetMove()))

	return &pf.MoveResponse{}, nil
}

func protoMoveTypeToModel(me pf.MoveRequest_MoveEnum) shared.MoveType {

	switch me {
	case pf.MoveRequest_ROTATELEFT:
		return shared.ROTATELEFT
	case pf.MoveRequest_ROTATERIGHT:
		return shared.ROTATERIGHT
	case pf.MoveRequest_MOVELEFT:
		return shared.MOVELEFT
	case pf.MoveRequest_MOVERIGHT:
		return shared.MOVERIGHT
	case pf.MoveRequest_DROP:
		return shared.DROP
	case pf.MoveRequest_DOWN:
		return shared.DOWN
	default:
		panic("oops") //fixme
	}
}

func pixelsToSquares(pixels []Pixel) []*pf.Square {
	squares := make([]*pf.Square, 0)

	for _, p := range pixels {
		squares = append(squares, &pf.Square{
			X:     uint32(p.X),
			Y:     uint32(p.Y),
			Color: colorToProto(p.Color),
		})
	}

	return squares
}

func colorToProto(c int) pf.Square_ColorEnum {

	switch c {
	case 0:
		return pf.Square_MAGENTA
	case 1:
		return pf.Square_CYAN
	case 2:
		return pf.Square_YELLOW
	case 3:
		return pf.Square_BLUE
	case 4:
		return pf.Square_GREEN
	case 5:
		return pf.Square_RED
	case 6:
		return pf.Square_WHITE
	default:
		return pf.Square_BLACK
	}


}