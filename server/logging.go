package server

import (
	"github.com/pauljeremyturner/dockerised-tetris/shared"
	"log"
	"os"
	"sync"
	"time"
)

var appLog *shared.Logger
var once sync.Once

func GetFileLogger() *shared.Logger {
	once.Do(func() {
		appLog = createLogger("tetris-server.log")
	})
	return appLog
}

func createLogger(fname string) *shared.Logger {
	file, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)

	if err != nil {
		log.Panicf("Could not start logger %s", err)
	}

	return &shared.Logger{
		Logger: log.New(file, time.Now().Format(time.RFC3339), log.Lshortfile),
	}
}
