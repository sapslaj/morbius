package server

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/cloudflare/goflow/v3/utils"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sapslaj/morbius/transport"
)

type ServerPortConfig struct {
	Enable    bool   `yaml:"enable"`
	Port      int    `yaml:"port"`
	Addr      string `yaml:"addr"`
	Workers   int    `yaml:"workers"`
	ReusePort bool   `yaml:"reuse_port"`
}

func mergeDefaultServerPortConfig(in *ServerPortConfig, port int) *ServerPortConfig {
	if in == nil {
		in = &ServerPortConfig{}
	}
	if in.Port == 0 {
		in.Port = port
	}
	if in.Addr == "" {
		in.Addr = "0.0.0.0"
	}
	if in.Workers == 0 {
		in.Workers = 1
	}
	return in
}

type ServerConfig struct {
	Transport Transport
	Logger    Logger
	NetFlowV5 *ServerPortConfig `yaml:"netflowv5"`
	NetFlowV9 *ServerPortConfig `yaml:"netflowv9"`
	SFlow     *ServerPortConfig `yaml:"sflow"`
	HTTP      *ServerPortConfig `yaml:"http"`
}

type Server struct {
	Config *ServerConfig
}

func NewServerWithTransportAndLogger(config ServerConfig, transport Transport, logger Logger) *Server {
	config.Transport = transport
	config.Logger = logger
	return NewServer(config)
}

func NewServer(config ServerConfig) *Server {
	config.NetFlowV5 = mergeDefaultServerPortConfig(config.NetFlowV5, 2055)
	config.NetFlowV9 = mergeDefaultServerPortConfig(config.NetFlowV9, 2056)
	config.SFlow = mergeDefaultServerPortConfig(config.SFlow, 6343)
	config.HTTP = mergeDefaultServerPortConfig(config.HTTP, 6060)
	if config.Logger == nil {
		config.Logger = &transport.StderrLogger{}
	}
	s := &Server{
		Config: &config,
	}
	return s
}

func (s *Server) IsRunnable() bool {
	if s.Config.NetFlowV5.Enable || s.Config.NetFlowV9.Enable || s.Config.SFlow.Enable {
		return true
	}
	return false
}

func (s *Server) RunAll() {
	var wg sync.WaitGroup

	if s.Config.NetFlowV5.Enable {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Config.Logger.Fatal(s.RunNetFlowV5())
		}()
	}

	if s.Config.NetFlowV9.Enable {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Config.Logger.Fatal(s.RunNetFlowV9())
		}()
	}

	if s.Config.SFlow.Enable {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Config.Logger.Fatal(s.RunSFlow())
		}()
	}

	if s.Config.HTTP.Enable {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Config.Logger.Fatal(s.RunHTTP())
		}()
	}

	wg.Wait()
	panic("fuck this shouldn't happen")
}

func (s *Server) RunNetFlowV5() error {
	state := utils.StateNFLegacy{
		Transport: s.Config.Transport,
		Logger:    s.Config.Logger,
	}
	return state.FlowRoutine(
		s.Config.NetFlowV5.Workers,
		s.Config.NetFlowV5.Addr,
		s.Config.NetFlowV5.Port,
		s.Config.NetFlowV5.ReusePort,
	)
}

func (s *Server) RunNetFlowV9() error {
	state := utils.StateNetFlow{
		Transport: s.Config.Transport,
		Logger:    s.Config.Logger,
	}
	return state.FlowRoutine(
		s.Config.NetFlowV9.Workers,
		s.Config.NetFlowV9.Addr,
		s.Config.NetFlowV9.Port,
		s.Config.NetFlowV9.ReusePort,
	)
}

func (s *Server) RunSFlow() error {
	state := utils.StateSFlow{
		Transport: s.Config.Transport,
		Logger:    s.Config.Logger,
	}
	return state.FlowRoutine(
		s.Config.SFlow.Workers,
		s.Config.SFlow.Addr,
		s.Config.SFlow.Port,
		s.Config.SFlow.ReusePort,
	)
}

func (s *Server) RunHTTP() error {
	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(
		fmt.Sprintf("%s:%d", s.Config.HTTP.Addr, s.Config.HTTP.Port),
		nil,
	)
}
