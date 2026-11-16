package log

import (
	"io"
	"log"
	"os"
	"sync"
)

var errorLog = log.New(os.Stdout, "\033[31m[error]\033[0m ", log.LstdFlags|log.Lshortfile)
var infoLog = log.New(os.Stdout, "\033[34m[info ]\033[0m ", log.LstdFlags|log.Lshortfile)
var loggers = []*log.Logger{errorLog, infoLog}
var mu sync.Mutex

var Error = errorLog.Println
var Errorf = errorLog.Printf
var Info = infoLog.Println
var Infof = infoLog.Printf

const (
	InfoLevel = iota
	ErrorLevel
	Disabled
)

func SetLevel(level int) {
	mu.Lock()
	defer mu.Unlock()
	for _, logger := range loggers {
		logger.SetOutput(os.Stdout)
	}
	if ErrorLevel < level {
		errorLog.SetOutput(io.Discard)
	}
	if InfoLevel < level {
		infoLog.SetOutput(io.Discard)
	}
}
