package main

import (
	"sync"

	"github.com/cloudflare/goflow/v3/utils"
	"github.com/sapslaj/morbius/destination"
	"github.com/sapslaj/morbius/enricher"
	"github.com/sapslaj/morbius/transport"
)

func main() {
	host := "0.0.0.0"
	v5port := 2055
	v9port := 2056
	sFlowPort := 6343

	logger := &transport.StderrLogger{}
	protnamesEnricher := enricher.NewProtonamesEnricher()
	lokiDestination := destination.NewLokiDestination()
	elasticsearchDestination := destination.NewElasticsearchDestination()
	// stdoutDestination := StdoutDestination{}
	transport := &transport.Transport{
		Enrichers: []enricher.Enricher{
			&protnamesEnricher,
		},
		Destinations: []destination.Destination{
			// &stdoutDestination,
			&lokiDestination,
			&elasticsearchDestination,
		},
	}

	logger.Printf("It's Morbin' Time!")
	logger.Printf("v5:\t%s:%d", host, v5port)
	logger.Printf("v9:\t%s:%d", host, v9port)
	logger.Printf("sFlow:\t%s:%d", host, sFlowPort)

	sNF5 := utils.StateNFLegacy{
		Transport: transport,
		Logger:    logger,
	}

	sNF9 := utils.StateNetFlow{
		Transport: transport,
		Logger:    logger,
	}

	sSflow := utils.StateSFlow{
		Transport: transport,
		Logger:    logger,
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Fatal(sNF5.FlowRoutine(1, host, v5port, false))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Fatal(sNF9.FlowRoutine(1, host, v9port, false))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Fatal(sSflow.FlowRoutine(1, host, sFlowPort, false))
	}()

	wg.Wait()
}
