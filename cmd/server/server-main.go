package main

import (
	server "github.com/pauljeremyturner/dockerised-tetris/server"
	"github.com/pauljeremyturner/dockerised-tetris/shared"
	"sync"
)

const (
	port = ":50051"
)

var serverLog *shared.Logger

type TetrisGame struct{}

type tetrisServer struct {
	activeGames sync.Map
}

func main() {
	serverLog = server.GetFileLogger()

}
