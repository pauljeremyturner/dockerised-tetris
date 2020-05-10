package main

import (
	"github.com/google/uuid"
	"github.com/pauljeremyturner/dockerised-tetris/client"
	"github.com/pauljeremyturner/dockerised-tetris/shared"
	"runtime"
)

var tp client.ProtoClient
var ui client.TetrisUi
var clientSession *client.ClientSession

func main() {
	runtime.GOMAXPROCS(2)
	uuid, _ := uuid.NewRandom()
	clientSession = &client.ClientSession{
		Uuid:               uuid,
		PlayerName:         "paul",
		MoveChannel:        make(chan shared.MoveType, 10),
		BoardUpdateChannel: make(chan client.GameState, 10),
	}

	ui = client.NewTetrisUi(clientSession)
	tp = client.NewTetrisProto(clientSession)

	go tp.ReceiveStream(uuid, "paul")
	go tp.ListenToMove()
	go ui.ListenToBoardUpdates()

	ui.StartGame()
}
