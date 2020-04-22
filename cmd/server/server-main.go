package main

import (
	pf "github.com/pauljeremyturner/dockerised-tetris/protofiles"
	"golang.org/x/net/context"
)

const (
	port = ":50051"
)

type server struct{}

func (s *server) GameMove(ctx context.Context, in *pf.GameMove) (*pf.GameMoveResponse, error) {

	return &pf.GameMoveResponse{}, nil
}

/*

func (s *server) MakeBoard(ctx context.Context, in *pf.NewBoardRequest) (*pf.BoardState, error) {

	return &pf.BoardState{
		PlayerName: "foo",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pf.RegisterMakeBoardServer(s, &server{})
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
*/
