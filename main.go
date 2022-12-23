package main

import (
	"flag"

	"github.com/kr/pretty"
	"github.com/sapslaj/morbius/config"

	_ "net/http/pprof"
)

func main() {
	configFile := flag.String("config-file", "config.yaml", "Path to the config file")
	printConfig := flag.Bool("print-config", false, "Print configuration before starting")
	flag.Parse()

	c := config.NewFromFile(*configFile)
	server := c.BuildServer()
	logger := server.Config.Logger

	if *printConfig {
		logger.Printf("%# v", pretty.Formatter(c))
	}

	if !server.IsRunnable() {
		panic("Server is not runnable. Check server configuration.")
	}

	if !server.Config.NoFunAllowed {
		logger.Printf("It's Morbin' Time!")

		if server.Config.NetFlowV5.Enable {
			logger.Printf("v5:\t%s:%d", server.Config.NetFlowV5.Addr, server.Config.NetFlowV5.Port)
		}
		if server.Config.NetFlowV9.Enable {
			logger.Printf("v9:\t%s:%d", server.Config.NetFlowV9.Addr, server.Config.NetFlowV9.Port)
		}
		if server.Config.SFlow.Enable {
			logger.Printf("sFlow:\t%s:%d", server.Config.SFlow.Addr, server.Config.SFlow.Port)
		}
		if server.Config.HTTP.Enable {
			logger.Printf("http:\t%s:%d", server.Config.HTTP.Addr, server.Config.HTTP.Port)
		}
	}

	server.RunAll()
}
