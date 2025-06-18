package logger

import (
	"log"
	"os"
)

var (
	Info  = log.New(os.Stdout, "INFO: ", log.Ldate|log.Lmicroseconds|log.LUTC|log.Lshortfile)
	Warn  = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Lmicroseconds|log.LUTC|log.Lshortfile)
	Error = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Lmicroseconds|log.LUTC|log.Lshortfile)
)
