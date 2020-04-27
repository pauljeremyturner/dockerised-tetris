package client

import (
	"log"
	"os"
	"sync"
	"time"
)

type Logger struct {
	filename string
	*log.Logger
}

var appLog *Logger
var once sync.Once

func GetFileLogger() *Logger {
	once.Do(func() {
		appLog = createLogger("tetris-client.log")
	})
	return appLog
}

func createLogger(fname string) *Logger {
	file, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)

	if err != nil {
		log.Panicf("Could not start logger %s", err)
	}

	return &Logger{
		Logger:   log.New(file, time.Now().Format(time.RFC3339), log.Lshortfile),
	}
}
