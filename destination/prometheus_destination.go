package destination

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/mitchellh/hashstructure/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sapslaj/morbius/syncmap"
)

type PrometheusDestinationConfig struct {
	Namespace    string
	MetricLabels []string
	CountBytes   bool
	CountPackets bool
}

type prometheusDestinationMetric struct {
	bytes       uint64
	packets     uint64
	lastReceive time.Time
	labels      prometheus.Labels
}

type PrometheusDestination struct {
	Config            *PrometheusDestinationConfig
	byteCounterDesc   *prometheus.Desc
	packetCounterDesc *prometheus.Desc
	metricStore       *syncmap.Map[uint64, *prometheusDestinationMetric]
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
	metricStore, err := syncmap.NewMap[uint64, *prometheusDestinationMetric]()
	if err != nil {
		panic(err)
	}
	d := PrometheusDestination{
		Config:      config,
		metricStore: &metricStore,
	}
	if d.Config.CountBytes {
		d.byteCounterDesc = prometheus.NewDesc(
			fmt.Sprintf("%s_bytes", d.Config.Namespace),
			"Total number of bytes",
			d.Config.MetricLabels,
			nil,
		)
	}
	if d.Config.CountPackets {
		d.packetCounterDesc = prometheus.NewDesc(
			fmt.Sprintf("%s_packets", d.Config.Namespace),
			"Total number of packets",
			d.Config.MetricLabels,
			nil,
		)
	}
	prometheus.MustRegister(&d)
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

	bytes := uint64(msg["bytes"].(int))
	packets := uint64(msg["packets"].(int))

	hash, err := hashstructure.Hash(promLabels, hashstructure.FormatV2, nil)
	if err != nil {
		panic(err)
	}

	metric, loaded := d.metricStore.Load(hash)
	if !loaded {
		metric = &prometheusDestinationMetric{
			labels: promLabels,
		}
		d.metricStore.Store(hash, metric)
	}
	metric.lastReceive = time.Now()
	if d.Config.CountBytes {
		atomic.AddUint64(&metric.bytes, bytes)
	}
	if d.Config.CountPackets {
		atomic.AddUint64(&metric.packets, packets)
	}
}

func (d *PrometheusDestination) Describe(ch chan<- *prometheus.Desc) {
	if d.byteCounterDesc != nil {
		ch <- d.byteCounterDesc
	}
	if d.packetCounterDesc != nil {
		ch <- d.packetCounterDesc
	}
}

func (d *PrometheusDestination) Collect(ch chan<- prometheus.Metric) {
	d.metricStore.Range(func(key uint64, metric *prometheusDestinationMetric) bool {
		labelValues := make([]string, 0)
		for _, key := range d.Config.MetricLabels {
			labelValues = append(labelValues, metric.labels[key])
		}
		if d.byteCounterDesc != nil {
			ch <- prometheus.MustNewConstMetric(d.byteCounterDesc, prometheus.CounterValue, float64(metric.bytes), labelValues...)
		}
		if d.packetCounterDesc != nil {
			ch <- prometheus.MustNewConstMetric(d.packetCounterDesc, prometheus.CounterValue, float64(metric.packets), labelValues...)
		}
		return true
	})
}
