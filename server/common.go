package server

import (
	goflowpb "github.com/cloudflare/goflow/v3/pb"
)

type Logger interface {
	Printf(string, ...interface{})
	Errorf(string, ...interface{})
	Warnf(string, ...interface{})
	Warn(...interface{})
	Error(...interface{})
	Debug(...interface{})
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Fatalf(string, ...interface{})
	Fatal(...interface{})
}

type Transport interface {
	Publish([]*goflowpb.FlowMessage)
	PublishMessage(msg map[string]interface{})
}
