package main

import (
	"fmt"
	"sync"

	"github.com/cloudflare/goflow/v3/utils"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sapslaj/morbius/destination"
	"github.com/sapslaj/morbius/enricher"
	"github.com/sapslaj/morbius/transport"

	"net/http"
	_ "net/http/pprof"
)

func main() {
	host := "0.0.0.0"
	v5port := 2055
	v9port := 2056
	sFlowPort := 6343
	httpPort := 6060

	var enrichers []enricher.Enricher
	var destinations []destination.Destination

	protnamesEnricher := enricher.NewProtonamesEnricher()
	enrichers = append(enrichers, &protnamesEnricher)
	rdnsEnricher := enricher.NewRDNSEnricher()
	enrichers = append(enrichers, &rdnsEnricher)

	lokiDestination := destination.NewLokiDestination()
	destinations = append(destinations, &lokiDestination)
	elasticsearchDestination := destination.NewElasticsearchDestination()
	destinations = append(destinations, &elasticsearchDestination)
	// stdoutDestination := destination.StdoutDestination{}
	// destinations = append(destinations, &stdoutDestination)

	logger := &transport.StderrLogger{}
	transport := &transport.Transport{
		Enrichers:    enrichers,
		Destinations: destinations,
	}

	logger.Printf("It's Morbin' Time!")
	logger.Printf("v5:\t%s:%d", host, v5port)
	logger.Printf("v9:\t%s:%d", host, v9port)
	logger.Printf("sFlow:\t%s:%d", host, sFlowPort)
	logger.Printf("http:\t%s:%d", host, httpPort)

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

	wg.Add(1)
	go func() {
		defer wg.Done()
		http.Handle("/metrics", promhttp.Handler())
		logger.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", host, httpPort), nil))
	}()

	wg.Wait()
}
