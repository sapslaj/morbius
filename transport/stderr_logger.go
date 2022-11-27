package transport

import (
	"log"
	"os"
)

// Implementation of github.com/cloudflare/goflow/v3/utils.Logger
// (https://pkg.go.dev/github.com/cloudflare/goflow/v3/utils#Logger)
type StderrLogger struct {
}

func (l *StderrLogger) makeLogger() *log.Logger {
	return log.New(os.Stderr, "", 0)
}

func (l *StderrLogger) Printf(s string, i ...interface{}) {
	l.makeLogger().Printf(s, i...)
}

func (l *StderrLogger) Errorf(s string, i ...interface{}) {
	l.makeLogger().Printf(s, i...)
}

func (l *StderrLogger) Warnf(s string, i ...interface{}) {
	l.makeLogger().Printf(s, i...)
}

func (l *StderrLogger) Debugf(s string, i ...interface{}) {
	l.makeLogger().Printf(s, i...)
}

func (l *StderrLogger) Infof(s string, i ...interface{}) {
	l.makeLogger().Printf(s, i...)
}

func (l *StderrLogger) Fatalf(s string, i ...interface{}) {
	l.makeLogger().Fatalf(s, i...)
}

func (l *StderrLogger) Warn(i ...interface{}) {
	l.makeLogger().Print(i...)
}

func (l *StderrLogger) Error(i ...interface{}) {
	l.makeLogger().Print(i...)
}

func (l *StderrLogger) Debug(i ...interface{}) {
	l.makeLogger().Print(i...)
}

func (l *StderrLogger) Fatal(i ...interface{}) {
	l.makeLogger().Fatal(i...)
}
