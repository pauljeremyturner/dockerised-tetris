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

type protoServerState struct {
	tetris Tetris
}

type ProtoServer interface {
	CreateNewGame() uuid.UUID
	Move(uuid uuid.UUID)
}

func StartProtoServer(t Tetris) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pf.RegisterMoveServer(s, &protoServerState{tetris: t})
	pf.RegisterStartGameServer(s, &protoServerState{tetris: t})
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func (s *protoServerState) StartGame(in *pf.NewGameRequest, stream pf.StartGame_StartGameServer) error {

	u, _ := uuid.Parse(in.Uuid)
	ss := s.tetris.StartNewGame(Player{uuid: u, playerName: in.PlayerName})

	for gs := range ss.gameQueue {

		GetFileLogger().Println("Receive board update")

		if err := stream.Send(&pf.GameUpdateResponse{
			Uuid:       in.Uuid,
			PlayerName: in.PlayerName,
			GameOver:   false,
			Score:      uint32(gs.Score),
			Duration:   0,
			Squares:    pixelsToSquares(gs.Pixels),
			NextPiece:  pixelsToSquares(gs.NextPiece),
		}); err != nil {
			fmt.Printf("something went wrong %s", err)
		}
	}
	return nil
}

func (s *protoServerState) Move(ctx context.Context, in *pf.MoveRequest) (*pf.MoveResponse, error) {
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
	squares := make([]*pf.Square, len(pixels))

	for _, p := range pixels {
		squares = append(squares, &pf.Square{
			X:     uint32(p.X),
			Y:     uint32(p.Y),
			Color: uint32(p.Color),
		})
	}

	return squares
}
