package main

import (
	"github.com/google/uuid"
	"github.com/pauljeremyturner/dockerised-tetris/client"
	"github.com/pauljeremyturner/dockerised-tetris/shared"
)

var tp client.ProtoClient
var ui client.TetrisClientUi
var clientSession client.ClientSession

func main() {

	uuid, _ := uuid.NewRandom()
	clientSession = client.ClientSession{
		Uuid:               uuid,
		PlayerName:         "paul",
		MoveChannel:        make(chan shared.MoveType, 10),
		BoardUpdateChannel: make(chan client.GameState, 10),
	}

	ui = client.NewTetrisClientUi(clientSession)

	tp = client.NewTetrisProto(clientSession)

	go tp.ReceiveStream(uuid, "paul")

	go tp.ListenToMove()
	go ui.ListenToBoardUpdates()
	ui.StartGame()
}
