
package main

import (
	"log"

	pf "github.com/pauljeremyturner/dockerised-tetris/protofiles"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	c := pf.NewMakeBoardClient(conn)

	r, err := c.MakeBoard(context.Background(), &pf.NewBoardRequest{PlayerName: "foobar" })

	if err != nil {
		log.Fatalf("New board failed %v", err)
	}
	log.Printf("New Board confirmed: %t", r.PlayerName)
}