package main

import (
	server "github.com/pauljeremyturner/dockerised-tetris/server"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	tetris := server.NewTetris()
	server.StartProtoServer(tetris)
}
