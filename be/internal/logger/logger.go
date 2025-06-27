package logger

import (
	"go.uber.org/zap"
)

// zapWrapper adapts zap's SugaredLogger to the stdlib logger interface used in
// the rest of the codebase.
type zapWrapper struct {
	logFunc   func(args ...interface{})
	logfFunc  func(format string, args ...interface{})
	fatalFunc func(args ...interface{})
}

func (z zapWrapper) Println(v ...interface{})               { z.logFunc(v...) }
func (z zapWrapper) Printf(format string, v ...interface{}) { z.logfFunc(format, v...) }
func (z zapWrapper) Fatalln(v ...interface{})               { z.fatalFunc(v...) }

var (
	Info  zapWrapper
	Warn  zapWrapper
	Error zapWrapper
)

func init() {
	l, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	s := l.Sugar()
	Info = zapWrapper{s.Info, s.Infof, s.Fatal}
	Warn = zapWrapper{s.Warn, s.Warnf, s.Fatal}
	Error = zapWrapper{s.Error, s.Errorf, s.Fatal}
}
