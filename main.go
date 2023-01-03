package main

import (
	"flag"
	"os"

	"github.com/kr/pretty"
	"github.com/sapslaj/morbius/config"

	_ "net/http/pprof"
)

func envWithDefault(key string, def string) string {
	value := os.Getenv(key)
	if value == "" {
		return def
	}
	return value
}

func main() {
	// Cannot use default flagset due to other packages (somewhat infuriatingly)
	// registering their own flags.
	f := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	configFile := f.String("config-file", envWithDefault("MORBIUS_CONFIG_FILE", "config.yaml"), "Path to the config file")
	printConfig := f.Bool("print-config", false, "Print configuration before starting")
	f.Parse(os.Args[1:])

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
