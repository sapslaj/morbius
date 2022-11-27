package destination

import (
	"encoding/json"
	"time"

	"github.com/sapslaj/morbius/promtail"
)

type LokiDestination struct {
	client *promtail.JSONClient
}

func NewLokiDestination() LokiDestination {
	client, err := promtail.NewClientJson(promtail.ClientConfig{
		// Yes I know this is a deprecated endpoint. The real one is buggy with JSON input.
		// https://github.com/grafana/loki/issues/4837
		PushURL:            "http://localhost:3100/api/prom/push",
		Labels:             "{job=\"netflow\"}",
		BatchWait:          1 * time.Second,
		BatchEntriesNumber: 10000,
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
