package client

import (
	"log"
	"os"
	"sync"
	"time"
)

type logger struct {
	filename string
	*log.Logger
}

var appLog *logger
var once sync.Once

func GetFileLogger() *logger {
	once.Do(func() {
		appLog = createLogger("log.log")
	})
	return appLog
}

func createLogger(fname string) *logger {
	file, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)

	if err != nil {
		log.Panicf("Could not start logger %s", err)
	}

	return &logger{
		filename: fname,
		Logger:   log.New(file, time.Now().Format(time.RFC3339), log.Lshortfile),
	}
}
