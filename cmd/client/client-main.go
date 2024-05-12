package main

import (
	"github.com/google/uuid"
	"github.com/pauljeremyturner/dockerised-tetris/client"
	"github.com/pauljeremyturner/dockerised-tetris/shared"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"runtime"
	"strconv"
	"time"
)

var tp client.ProtoClient
var ui client.TetrisUi
var clientSession *client.ClientSession

func main() {

	sugar, err := initialiseFileLogger()
	if err != nil {
		panic(err)
	}

	runtime.GOMAXPROCS(2)
	uuid, _ := uuid.NewRandom()
	clientSession = &client.ClientSession{
		Uuid:               uuid,
		PlayerName:         "paul",
		MoveChannel:        make(chan shared.MoveType, 10),
		BoardUpdateChannel: make(chan client.GameState, 10),
	}

	ui = client.NewTetrisUi(clientSession, sugar)
	tp = client.NewTetrisProto(clientSession, sugar)

	go tp.ReceiveStream(uuid, "paul")
	go tp.ListenToMove()

	ui.StartGame()

}

func initialiseFileLogger() (*zap.SugaredLogger, error) {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	unixTimeString := strconv.FormatInt(time.Now().Unix(), 10)
	wd, err := os.Getwd()
	if err != nil {
		return nil, client.NewClientSystemError("could not get WD", err)
	}
	logfileName := wd + "/" + unixTimeString
	logFile, err := os.OpenFile(logfileName, os.O_CREATE, 0777)
	if err != nil {
		return nil, client.NewClientSystemError("could not open logfile", err)
	}
	fileOutput := zapcore.AddSync(logFile)

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		fileOutput,
		zap.InfoLevel,
	)

	logger := zap.New(core)

	return logger.Sugar(), nil
}
