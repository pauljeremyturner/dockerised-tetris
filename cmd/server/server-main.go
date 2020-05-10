package main

import (
	server "github.com/pauljeremyturner/dockerised-tetris/server"
	"math/rand"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(2)
	rand.Seed(time.Now().UnixNano())
	tetris := server.NewTetris()
	server.StartProtoServer(tetris)
}
