package main

import (
	"github.com/pauljeremyturner/dockerised-tetris/protogen"
	server "github.com/pauljeremyturner/dockerised-tetris/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func main() {

	l, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	runtime.GOMAXPROCS(2)
	rand.Seed(time.Now().UnixNano())

	tetris := server.NewTetris(l.Sugar())
	ps := server.NewProtoServer(tetris)
	s := grpc.NewServer()
	protogen.RegisterGameServer(s, ps)
	protogen.RegisterMoveServer(s, ps)
	reflection.Register(s)

	lis, err := net.Listen("tcp", ":50051") // Change the port as needed
	if err != nil {
		panic(err)
	}

	go func() {
		l.Sugar().Infof("gRPC server listening on %s", lis.Addr())
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM)

	<-sigChan

	s.GracefulStop()

	l.Sugar().Info("gRPC server gracefully stopped")

}
