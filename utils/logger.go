package utils

import (
	"log"
	"os"
)

const (
	ErrorLevel int = iota
	WarnLevel
	InfoLevel
	DebugLevel
)

var (
	logger   *log.Logger
	logLevel int
)

func LogDebugf(format string, args ...interface{}) {
	if logLevel >= DebugLevel {
		logger.Printf(format, args...)
	}

}

func LogInfof(format string, args ...interface{}) {
	if logLevel >= InfoLevel {
		logger.Printf(format, args...)
	}
}

func LogWarnf(format string, args ...interface{}) {
	if logLevel >= WarnLevel {
		logger.Printf(format, args...)
	}
}

func LogErrorf(format string, args ...interface{}) {
	if logLevel >= ErrorLevel {
		logger.Printf(format, args...)
	}
}

func SetConfig(level int, logfile string) {
	logLevel = level
	file, err := os.OpenFile(logfile, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm|os.ModeTemporary)
	if err != nil {
		log.Print(logfile)
		log.Fatalln(err.Error())
	}
	logger = log.New(file, "", log.LstdFlags)
}
