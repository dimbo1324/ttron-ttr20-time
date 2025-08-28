package logging

import (
	"log"
	"os"
)

func New(path string) *log.Logger {
	if path == "" {
		return log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)
	}

	fileObj, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		l := log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)
		l.Printf("cannot open log file %s: %v — logging to stdout", path, err)
		return l
	}
	return log.New(fileObj, "", log.LstdFlags|log.Lmicroseconds)
}
