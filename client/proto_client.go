package client

import (
	pf "github.com/pauljeremyturner/dockerised-tetris/protofiles"
	"google.golang.org/grpc"
	"log"
)

var (
	address = "localhost:50051"
)

type UpdateBoard func(gs GameState)

type TetrisProto interface {
	Move(r rune)
}

func NewTetrisProto(ub UpdateBoard) TetrisProto {
	return ProtoClientState{
		updateBoard: ub,
		appLog:      GetFileLogger().Logger,
	}
}

type ProtoClientState struct {
	updateBoard UpdateBoard
	appLog      *log.Logger
}

func (r ProtoClientState) Move(char rune) {

	r.appLog.Printf("send move to proto %s", char)

}

func (pcs ProtoClientState) ReceiveStream(client pf.StartGameClient, request *pf.NewGameRequest) {
	/*
		stream, err := client.StartGame(context.Background(), request)
		if err != nil {
			log.Fatalf("%v.StartGame(_) = _, %v", client, err)
		}
		// Listen to the stream of messages
		for {
			gameUpdate, err := stream.Recv()
			if err == io.EOF {
				// If there are no more messages, get out of loop
				break
			}
			if err != nil {
				log.Fatalf("%v.StartGame(_) = _, %v", client, err)
			}

			fmt.Println(gameUpdate)
	*/
	for i := 0; i < 100; i++ {
		gs := GameState{
			Blocks:    nil,
			NextPiece: nil,
			GameOver:  false,
			Lines:     0,
			Duration:  0,
		}

		pcs.updateBoard(gs)
	}

}

func stuff() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	//client := pf.NewStartGameClient(conn)

	// Contact the server and print out its response.
	//ReceiveStream(client, &pf.NewGameRequest{})
}
