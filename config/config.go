package config

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/sapslaj/morbius/destination"
	"github.com/sapslaj/morbius/enricher"
	"github.com/sapslaj/morbius/server"
	"github.com/sapslaj/morbius/transport"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Server    *server.ServerConfig   `yaml:"server"`
	Transport map[string]interface{} `yaml:"transport"` // TODO: better way of handling this (see Config.BuildTransport)
	Enrichers struct {
		AddrType    *enricher.AddrTypeEnricherConfig    `yaml:"addr_type"`
		MaxmindDB   *enricher.MaxmindDBEnricherConfig   `yaml:"maxmind_db"`
		NetDB       *enricher.NetDBEnricherConfig       `yaml:"netdb"`
		ProtoNames  *enricher.ProtonamesEnricherConfig  `yaml:"proto_names"`
		RDNS        *enricher.RDNSEnricherConfig        `yaml:"rdns"`
		FieldMapper *enricher.FieldMapperEnricherConfig `yaml:"field_mapper"`
	} `yaml:"enrichers"`
	Destinations struct {
		Discard       *destination.DiscardDestinationConfig      `yaml:"discard"`
		Elasticsearch *destination.ElasticseachDestinationConfig `yaml:"elasticsearch"`
		Loki          *destination.LokiDestinationConfig         `yaml:"loki"`
		Prometheus    *destination.PrometheusDestinationConfig   `yaml:"prometheus"`
		Stdout        *destination.StdoutDestinationConfig       `yaml:"stdout"`
	} `yaml:"destinations"`
}

func NewFromFile(filename string) *Config {
	f, err := os.Open(filename)
	if err != nil {
		panic(fmt.Errorf("Config: error opening file %s: %w", filename, err))
	}
	return NewFromReader(f)
}

func NewFromReader(r io.ReadCloser) *Config {
	defer r.Close()
	b, err := io.ReadAll(r)
	if err != nil {
		panic(fmt.Errorf("Config: error reading: %w", err))
	}
	return NewFromBytes(b)
}

func NewFromString(s string) *Config {
	return NewFromBytes([]byte(s))
}

func NewFromBytes(b []byte) *Config {
	c := &Config{}
	err := yaml.Unmarshal(b, c)
	if err != nil {
		panic(fmt.Errorf("Config: error unmarshalling YAML: %w", err))
	}
	return c
}

func (c *Config) BuildEnrichers() []enricher.Enricher {
	var enrichers []enricher.Enricher
	if c.Enrichers.AddrType != nil {
		addrTypeEnricher := enricher.NewAddrTypeEnricher(c.Enrichers.AddrType)
		enrichers = append(enrichers, &addrTypeEnricher)
	}
	if c.Enrichers.MaxmindDB != nil {
		maxmindDBEnricher := enricher.NewMaxmindDBEnricher(c.Enrichers.MaxmindDB)
		enrichers = append(enrichers, &maxmindDBEnricher)
	}
	if c.Enrichers.NetDB != nil {
		netdbEnricher := enricher.NewNetDBEnricher(c.Enrichers.NetDB)
		enrichers = append(enrichers, &netdbEnricher)
	}
	if c.Enrichers.ProtoNames != nil {
		protnamesEnricher := enricher.NewProtonamesEnricher(c.Enrichers.ProtoNames)
		enrichers = append(enrichers, &protnamesEnricher)
	}
	if c.Enrichers.RDNS != nil {
		rdnsEnricher := enricher.NewRDNSEnricher(c.Enrichers.RDNS)
		enrichers = append(enrichers, &rdnsEnricher)
	}
	if c.Enrichers.FieldMapper != nil {
		fieldMapperEnricher := enricher.NewFieldMapperEnricher(c.Enrichers.FieldMapper)
		enrichers = append(enrichers, &fieldMapperEnricher)
	}
	return enrichers
}

func (c *Config) BuildDestinations() []destination.Destination {
	var destinations []destination.Destination
	if c.Destinations.Discard != nil {
		discardDestination := destination.NewDiscardDestination(c.Destinations.Discard)
		destinations = append(destinations, &discardDestination)
	}
	if c.Destinations.Elasticsearch != nil {
		elasticsearchDestination := destination.NewElasticsearchDestination(c.Destinations.Elasticsearch)
		destinations = append(destinations, &elasticsearchDestination)
	}
	if c.Destinations.Loki != nil {
		lokiDestination := destination.NewLokiDestination(c.Destinations.Loki)
		destinations = append(destinations, &lokiDestination)
	}
	if c.Destinations.Prometheus != nil {
		prometheusDestination := destination.NewPrometheusDestination(c.Destinations.Prometheus)
		destinations = append(destinations, &prometheusDestination)
	}
	if c.Destinations.Stdout != nil {
		stdoutDestionation := destination.NewStdoutDestination(c.Destinations.Stdout)
		destinations = append(destinations, &stdoutDestionation)
	}
	return destinations
}

func (c *Config) renderTemplate(v string, s any) (string, error) {
	var buf bytes.Buffer
	tmpl, err := template.New("config").Funcs(sprig.FuncMap()).Parse(v)
	if err != nil {
		return "", fmt.Errorf("Config.renderTemplate: unable to parse template: %w", err)
	}
	err = tmpl.Execute(&buf, s)
	if err != nil {
		return "", fmt.Errorf("Config.renderTemplate: unable to execute template: %w", err)
	}
	return buf.String(), nil
}

func (c *Config) BuildTransport() server.Transport {
	var t server.Transport
	enrichers := c.BuildEnrichers()
	destinations := c.BuildDestinations()
	tplValues := struct {
		NumCPU int
	}{
		NumCPU: runtime.NumCPU(),
	}
	parallelizeDestinations := MapGetDefault(c.Transport, "parallelize_destinations", false)
	switch MapGetDefault(c.Transport, "dispatch_method", "worker_pool") {
	case "linear":
		t = transport.NewLinearTransport(parallelizeDestinations, destinations, enrichers)
	case "worker_pool":
		workerCount := MapGetFunc(c.Transport, "worker_count", func(v any, present bool) int {
			if !present {
				return tplValues.NumCPU
			}
			switch value := v.(type) {
			case int:
				return value
			case string:
				value, err := c.renderTemplate(value, tplValues)
				if err != nil {
					panic(err)
				}
				valueInt, err := strconv.Atoi(value)
				if err != nil {
					panic(err)
				}
				return valueInt
			default:
				panic(fmt.Errorf("config: BuildTransport: unable to parse worker_count: invalid type %T", value))
			}
		})
		messageBuffer := MapGetFunc(c.Transport, "message_buffer", func(v any, present bool) int {
			if !present {
				return tplValues.NumCPU * 4
			}
			switch value := v.(type) {
			case int:
				return value
			case string:
				value, err := c.renderTemplate(value, tplValues)
				if err != nil {
					panic(err)
				}
				valueInt, err := strconv.Atoi(value)
				if err != nil {
					panic(err)
				}
				return valueInt
			default:
				panic(fmt.Errorf("config: BuildTransport: unable to parse message_buffer: invalid type %T", value))
			}
		})
		t = transport.NewWorkerPoolTransport(parallelizeDestinations, destinations, enrichers, workerCount, messageBuffer)
	case "goroutine":
		maxGoroutines := MapGetFunc(c.Transport, "max_goroutines", func(v any, present bool) int {
			if !present {
				return 0
			}
			switch value := v.(type) {
			case int:
				return value
			case string:
				value, err := c.renderTemplate(value, tplValues)
				if err != nil {
					panic(err)
				}
				valueInt, err := strconv.Atoi(value)
				if err != nil {
					panic(err)
				}
				return valueInt
			default:
				panic(fmt.Errorf("config: BuildTransport: unable to parse max_goroutines: invalid type %T", value))
			}
		})
		t = transport.NewGoroutineTransport(parallelizeDestinations, destinations, enrichers, int64(maxGoroutines))
	}
	return t
}

func (c *Config) BuildServer() server.Server {
	transport := c.BuildTransport()
	if c.Server == nil {
		c.Server = &server.ServerConfig{}
	}
	return *server.NewServerWithTransportAndLogger(*c.Server, transport, nil)
}
