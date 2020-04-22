package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	pf "github.com/pauljeremyturner/dockerised-tetris/protofiles"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var port = ":50051"

type debugServer struct {
}

func (s *debugServer) StartGame(in *pf.NewGameRequest, stream pf.StartGame_StartGameServer) error {

	for i := 0; i < 100; i++ {

		time.Sleep(1 * time.Second)

		fmt.Println("Sending stream message")

		if err := stream.Send(&pf.GameUpdateResponse{
			Uuid:       in.Uuid,
			PlayerName: in.PlayerName,
			GameOver:   false,
			Score:      uint32(i),
			Duration:   0,
			Squares: []*pf.Square{&pf.Square{
				X:     uint32(i%10) * 2,
				Y:     uint32(i%10) * 2,
				Color: 0,
			}},
			NextPiece: []*pf.Square{&pf.Square{
				X:     uint32(3),
				Y:     uint32(3),
				Color: 0,
			}},
		}); err != nil {
			fmt.Printf("something went wrong", err)
		}
	}
	return nil
}

func (s *debugServer) Move(ctx context.Context, in *pf.MoveRequest) (*pf.MoveResponse, error) {

	fmt.Printf("Received move %s\n", in)

	return &pf.MoveResponse{}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pf.RegisterMoveServer(s, &debugServer{})
	pf.RegisterStartGameServer(s, &debugServer{})
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
