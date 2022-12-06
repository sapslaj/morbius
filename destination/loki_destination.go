package destination

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/go-kit/log"
	"github.com/grafana/dskit/backoff"
	"github.com/grafana/dskit/flagext"
	"github.com/prometheus/common/model"

	"github.com/sapslaj/morbius/lokiclient"
	"github.com/sapslaj/morbius/lokiclient/api"
	lokiflag "github.com/sapslaj/morbius/lokiclient/flagext"
	"github.com/sapslaj/morbius/lokiclient/logproto"
)

type LokiDestinationConfig struct {
	PushURL        string
	StaticLabels   map[string]string
	DynamicLabels  []string
	BatchWait      time.Duration
	BatchSize      int
	MakeLokiSuffer bool
}

type LokiDestination struct {
	Config *LokiDestinationConfig
	client lokiclient.Client
}

func NewLokiDestination(config *LokiDestinationConfig) LokiDestination {
	if config == nil {
		config = &LokiDestinationConfig{}
	}
	if config.PushURL == "" {
		config.PushURL = "http://localhost:3100/api/prom/push"
	}
	if config.StaticLabels == nil {
		config.StaticLabels = map[string]string{
			"job": "netflow",
		}
	}
	if config.DynamicLabels == nil {
		config.DynamicLabels = make([]string, 0)
	}
	if config.BatchWait == 0 {
		config.BatchWait = lokiclient.BatchWait
	}
	if config.BatchSize == 0 {
		config.BatchSize = lokiclient.BatchSize
	}

	url, err := url.Parse(config.PushURL)
	if err != nil {
		panic(err)
	}

	externalLabels := make(model.LabelSet)
	for k, v := range config.StaticLabels {
		externalLabels[model.LabelName(k)] = model.LabelValue(v)
	}
	clientConfig := lokiclient.Config{
		URL:       flagext.URLValue{URL: url},
		BatchWait: config.BatchWait,
		BatchSize: config.BatchSize,
		ExternalLabels: lokiflag.LabelSet{
			LabelSet: externalLabels,
		},
		BackoffConfig: backoff.Config{
			MaxBackoff: lokiclient.MaxBackoff,
			MaxRetries: lokiclient.MaxRetries,
			MinBackoff: lokiclient.MinBackoff,
		},
		Timeout: lokiclient.Timeout,
	}
	client, err := lokiclient.New(clientConfig, []string{}, 0, log.NewLogfmtLogger(os.Stdout))
	if err != nil {
		panic(err)
	}
	d := LokiDestination{
		Config: config,
		client: client,
	}
	return d
}

func (d *LokiDestination) Publish(msg map[string]interface{}) {
	result, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	labelSet := make(model.LabelSet)
	if d.Config.MakeLokiSuffer {
		for key, value := range msg {
			labelSet[model.LabelName(key)] = model.LabelValue(fmt.Sprint(value))
		}
	} else {
		for _, key := range d.Config.DynamicLabels {
			value, ok := msg[key]
			if !ok {
				continue
			}
			labelSet[model.LabelName(key)] = model.LabelValue(fmt.Sprint(value))
		}
	}
	d.client.Chan() <- api.Entry{
		Labels: labelSet,
		Entry: logproto.Entry{
			Timestamp: time.Now(),
			Line:      string(result),
		},
	}
}
