package destination

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/mitchellh/hashstructure/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sapslaj/morbius/syncmap"
)

var (
	MetricPrometheusMetricStoreSize = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "prom_metric_store_size",
			Help: "size of the Prometheus metric store",
		},
	)
	MetricPrometheusMetricStoreGCSeconds = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "prom_metric_store_gc_seconds",
			Help:    "Seconds and count of metric store GC job",
			Buckets: []float64{0, .0001, .00025, .0005, .00075, .001, .0025, .005, .0075, .01, .025, .05, .075, .1, .25, .5, .75, 1},
		},
	)
	MetricPrometheusMetricStoreLastGCSeconds = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "prom_metric_store_last_gc_seconds",
			Help: "Duration of last metric store GC job run",
		},
	)
	MetricPrometheusMetricStoreEvictionCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "prom_metric_store_eviction_count",
			Help: "Number of metrics evicted based on visibility timeout",
		},
	)
	MetricPrometheusMetricStoreLastEvictionCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "prom_metric_store_last_eviction_count",
			Help: "Number of metrics evicted based on visibility timeout on the last GC job run",
		},
	)
)

func init() {
	prometheus.MustRegister(MetricPrometheusMetricStoreSize)
	prometheus.MustRegister(MetricPrometheusMetricStoreGCSeconds)
	prometheus.MustRegister(MetricPrometheusMetricStoreLastGCSeconds)
	prometheus.MustRegister(MetricPrometheusMetricStoreEvictionCount)
	prometheus.MustRegister(MetricPrometheusMetricStoreLastEvictionCount)
}

type PrometheusDestinationConfig struct {
	Namespace           string
	MetricLabels        []string
	CountBytes          bool
	CountPackets        bool
	CountFlows          bool
	ObserveFlowDuration bool
	FlowDurationBuckets []float64
	VisibilityTimeout   time.Duration
	GCInterval          time.Duration
}

type prometheusDestinationMetric struct {
	labels       prometheus.Labels
	lastReceive  time.Time
	bytes        uint64
	packets      uint64
	flows        uint64
	flowDuration prometheus.Histogram
}

type PrometheusDestination struct {
	Config            *PrometheusDestinationConfig
	byteCounterDesc   *prometheus.Desc
	packetCounterDesc *prometheus.Desc
	flowCounterDesc   *prometheus.Desc
	flowDurationDesc  *prometheus.Desc
	metricStore       *syncmap.Map[uint64, *prometheusDestinationMetric]
}

func NewPrometheusDestination(config *PrometheusDestinationConfig) PrometheusDestination {
	if config == nil {
		config = &PrometheusDestinationConfig{
			CountBytes:          true,
			CountPackets:        true,
			CountFlows:          true,
			ObserveFlowDuration: true,
		}
	}
	if config.Namespace == "" {
		config.Namespace = "netflow"
	}
	if config.VisibilityTimeout == 0 {
		config.VisibilityTimeout = 5 * time.Minute
	}
	if config.GCInterval == 0 {
		config.GCInterval = 15 * time.Second
	}
	if config.FlowDurationBuckets == nil {
		config.FlowDurationBuckets = prometheus.DefBuckets
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
	if d.Config.CountFlows {
		d.flowCounterDesc = prometheus.NewDesc(
			fmt.Sprintf("%s_flows", d.Config.Namespace),
			"Total number of packets",
			d.Config.MetricLabels,
			nil,
		)
	}
	if d.Config.ObserveFlowDuration {
		d.flowDurationDesc = prometheus.NewDesc(
			fmt.Sprintf("%s_flow_duration_seconds", d.Config.Namespace),
			"Duration of flow based on start and end timestamps",
			d.Config.MetricLabels,
			nil,
		)
	}
	prometheus.MustRegister(&d)
	d.startMetricStoreGC()
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

	hash, err := hashstructure.Hash(promLabels, hashstructure.FormatV2, nil)
	if err != nil {
		panic(err)
	}

	metric, loaded := d.metricStore.Load(hash)
	if !loaded {
		metric = &prometheusDestinationMetric{
			labels: promLabels,
		}
		if d.Config.ObserveFlowDuration {
			metric.flowDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
				Namespace:   d.Config.Namespace,
				Name:        "flow_duration_seconds",
				Buckets:     d.Config.FlowDurationBuckets,
				ConstLabels: promLabels,
			})
		}
		d.metricStore.Store(hash, metric)
	}
	metric.lastReceive = time.Now()
	if d.Config.CountBytes {
		if bytes, ok := msg["bytes"].(int); ok {
			atomic.AddUint64(&metric.bytes, uint64(bytes))
		}
	}
	if d.Config.CountPackets {
		if packets, ok := msg["packets"].(int); ok {
			atomic.AddUint64(&metric.packets, uint64(packets))
		}
	}
	if d.Config.CountFlows {
		atomic.AddUint64(&metric.flows, uint64(1))
	}
	if d.Config.ObserveFlowDuration {
		flowDuration, ok := func(msg map[string]interface{}) (time.Duration, bool) {
			flowStartRaw, ok := msg["time_flow_start"]
			if !ok {
				return 0, ok
			}
			flowStart, ok := flowStartRaw.(int)
			if !ok {
				return 0, ok
			}
			flowEndRaw, ok := msg["time_flow_end"]
			if !ok {
				return 0, ok
			}
			flowEnd, ok := flowEndRaw.(int)
			if !ok {
				return 0, ok
			}
			return time.Duration(flowEnd-flowStart) * time.Second, true
		}(msg)
		if ok {
			metric.flowDuration.Observe(flowDuration.Seconds())
		}
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
	d.metricStore.Range(func(hash uint64, metric *prometheusDestinationMetric) bool {
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
		if d.flowCounterDesc != nil {
			ch <- prometheus.MustNewConstMetric(d.flowCounterDesc, prometheus.CounterValue, float64(metric.flows), labelValues...)
		}
		if d.flowDurationDesc != nil {
			ch <- metric.flowDuration
		}
		return true
	})
}

func (d *PrometheusDestination) metricStoreGC() {
	storeSize := 0
	evictions := 0
	d.metricStore.Range(func(hash uint64, metric *prometheusDestinationMetric) bool {
		threshold := time.Now().Add(-d.Config.VisibilityTimeout)
		if metric.lastReceive.Before(threshold) {
			d.metricStore.Delete(hash)
			evictions++
			MetricPrometheusMetricStoreEvictionCount.Inc()
		} else {
			storeSize++
		}
		return true
	})
	MetricPrometheusMetricStoreLastEvictionCount.Set(float64(evictions))
	MetricPrometheusMetricStoreSize.Set(float64(storeSize))
}

func (d *PrometheusDestination) startMetricStoreGC() {
	go func() {
		for {
			start := time.Now()
			d.metricStoreGC()
			duration := time.Since(start)
			MetricPrometheusMetricStoreLastGCSeconds.Set(duration.Seconds())
			MetricPrometheusMetricStoreGCSeconds.Observe(duration.Seconds())
			timeUntilNext := d.Config.GCInterval - duration
			if timeUntilNext > 0 {
				time.Sleep(timeUntilNext)
			}
		}
	}()
}
