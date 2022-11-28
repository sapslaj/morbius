package destination

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/sapslaj/morbius/promtail"
)

type LokiDestinationConfig struct {
	PushURL            string
	Labels             map[string]string
	BatchWait          time.Duration
	BatchEntriesNumber int
}

type LokiDestination struct {
	Config *LokiDestinationConfig
	client *promtail.JSONClient
}

func NewLokiDestination(config *LokiDestinationConfig) LokiDestination {
	if config == nil {
		config = &LokiDestinationConfig{}
	}
	if config.PushURL == "" {
		// Yes I know this is a deprecated endpoint. The real one is buggy with JSON input.
		// https://github.com/grafana/loki/issues/4837
		config.PushURL = "http://localhost:3100/api/prom/push"
	}
	if config.Labels == nil {
		config.Labels = map[string]string{
			"job": "netflow",
		}
	}
	if config.BatchWait == 0 {
		config.BatchWait = 1 * time.Second
	}
	if config.BatchEntriesNumber == 0 {
		config.BatchEntriesNumber = 10000
	}
	client, err := promtail.NewClientJson(promtail.ClientConfig{
		PushURL:            config.PushURL,
		Labels:             makeLokiLabelString(config.Labels),
		BatchWait:          config.BatchWait,
		BatchEntriesNumber: config.BatchEntriesNumber,
	})
	if err != nil {
		panic(err)
	}
	return LokiDestination{
		client: client,
	}
}

func (d *LokiDestination) Publish(msg map[string]interface{}) {
	result, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	d.client.Send(string(result))
}

func makeLokiLabelString(labels map[string]string) string {
	var pairs []string
	for k, v := range labels {
		pairs = append(pairs, fmt.Sprintf("%s=\"%s\"", k, v))
	}
	return "{" + strings.Join(pairs, ",") + "}"
}
