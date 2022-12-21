package destination

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

type PrometheusDestinationConfig struct {
	Namespace    string
	MetricLabels []string
	CountBytes   bool
	CountPackets bool
}

type PrometheusDestination struct {
	Config        *PrometheusDestinationConfig
	byteCounter   *prometheus.CounterVec
	packetCounter *prometheus.CounterVec
}

func NewPrometheusDestination(config *PrometheusDestinationConfig) PrometheusDestination {
	if config == nil {
		config = &PrometheusDestinationConfig{
			CountBytes:   true,
			CountPackets: true,
		}
	}
	if config.Namespace == "" {
		config.Namespace = "netflow"
	}
	d := PrometheusDestination{
		Config: config,
	}
	if d.Config.CountBytes {
		d.byteCounter = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: d.Config.Namespace,
				Name:      "bytes",
				Help:      "Total number of bytes",
			},
			d.Config.MetricLabels,
		)
		prometheus.MustRegister(d.byteCounter)
	}
	if d.Config.CountPackets {
		d.packetCounter = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: d.Config.Namespace,
				Name:      "packets",
				Help:      "Total number of packets",
			},
			d.Config.MetricLabels,
		)
		prometheus.MustRegister(d.packetCounter)
	}
	return d
}

func (d *PrometheusDestination) Publish(msg map[string]interface{}) {
	promLabels := make(prometheus.Labels, 0)
	for _, label := range d.Config.MetricLabels {
		if value, ok := msg[label]; ok {
			promLabels[label] = fmt.Sprint(value)
		} else {
			// Must set a value otherwise it gets angry
			promLabels[label] = ""
		}
	}
	if d.byteCounter != nil {
		bytes := msg["bytes"].(int)
		d.byteCounter.With(promLabels).Add(float64(bytes))
	}
	if d.packetCounter != nil {
		packets := msg["packets"].(int)
		d.packetCounter.With(promLabels).Add(float64(packets))
	}
}
