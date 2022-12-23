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
	MetricPrometheusStoreSize = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "prom_store_size",
			Help: "size of the Prometheus metric store",
		},
		[]string{"store"},
	)
	MetricPrometheusStoreGCSeconds = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "prom_store_gc_seconds",
			Help:    "Seconds and count of metric store GC job",
			Buckets: []float64{0, .0001, .00025, .0005, .00075, .001, .0025, .005, .0075, .01, .025, .05, .075, .1, .25, .5, .75, 1},
		},
		[]string{"store"},
	)
	MetricPrometheusStoreLastGCSeconds = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "prom_store_last_gc_seconds",
			Help: "Duration of last metric store GC job run",
		},
		[]string{"store"},
	)
	MetricPrometheusStoreEvictionCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "prom_store_eviction_count",
			Help: "Number of metrics evicted based on visibility timeout",
		},
		[]string{"store"},
	)
	MetricPrometheusStoreLastEvictionCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "prom_store_last_eviction_count",
			Help: "Number of metrics evicted based on visibility timeout on the last GC job run",
		},
		[]string{"store"},
	)
)

func init() {
	prometheus.MustRegister(MetricPrometheusStoreSize)
	prometheus.MustRegister(MetricPrometheusStoreGCSeconds)
	prometheus.MustRegister(MetricPrometheusStoreLastGCSeconds)
	prometheus.MustRegister(MetricPrometheusStoreEvictionCount)
	prometheus.MustRegister(MetricPrometheusStoreLastEvictionCount)
}

type PrometheusDestinationConfig struct {
	Namespace           string
	MetricLabels        []string
	VisibilityTimeout   time.Duration
	GCInterval          time.Duration
	CountBytes          bool
	CountPackets        bool
	CountFlows          bool
	ObserveFlowDuration bool
	FlowDurationBuckets []float64
	ExportIpInfo        bool
	IpInfoLabels        []string
}

type prometheusDestinationMetric struct {
	labels       prometheus.Labels
	lastReceive  time.Time
	bytes        uint64
	packets      uint64
	flows        uint64
	flowDuration prometheus.Histogram
}

type prometheusDestinationIpInfo struct {
	labels      prometheus.Labels
	lastReceive time.Time
}

type PrometheusDestination struct {
	Config            *PrometheusDestinationConfig
	byteCounterDesc   *prometheus.Desc
	packetCounterDesc *prometheus.Desc
	flowCounterDesc   *prometheus.Desc
	flowDurationDesc  *prometheus.Desc
	ipInfoDesc        *prometheus.Desc
	metricStore       *syncmap.Map[uint64, *prometheusDestinationMetric]
	ipInfoStore       *syncmap.Map[string, *prometheusDestinationIpInfo]
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
	if config.IpInfoLabels == nil {
		config.IpInfoLabels = []string{"addr"}
	}
	metricStore, err := syncmap.NewMap[uint64, *prometheusDestinationMetric]()
	if err != nil {
		panic(err)
	}
	ipInfoStore, err := syncmap.NewMap[string, *prometheusDestinationIpInfo]()
	if err != nil {
		panic(err)
	}
	d := PrometheusDestination{
		Config:      config,
		metricStore: &metricStore,
		ipInfoStore: &ipInfoStore,
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
	if d.Config.ExportIpInfo {
		d.ipInfoDesc = prometheus.NewDesc(
			fmt.Sprintf("%s_ip_info", d.Config.Namespace),
			"Info metric for information about a particular IP address",
			d.Config.IpInfoLabels,
			nil,
		)
	}
	prometheus.MustRegister(&d)
	d.startMetricStoreGC()
	d.startIpStoreGC()
	return d
}

func (d *PrometheusDestination) Publish(msg map[string]interface{}) {
	d.storeFlowMetricsFromMsg(msg)

	if d.Config.ExportIpInfo {
		d.storeIpInfoFromMsg(msg, "src_addr", "src_")
		d.storeIpInfoFromMsg(msg, "src_addr", "dst_")
		d.storeIpInfoFromMsg(msg, "src_addr_encap", "src_encap_")
		d.storeIpInfoFromMsg(msg, "dst_addr_encap", "dst_encap_")
	}
}

func (d *PrometheusDestination) storeFlowMetricsFromMsg(msg map[string]interface{}) {
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

func (d *PrometheusDestination) storeIpInfoFromMsg(msg map[string]interface{}, addrField string, prefix string) {
	addrRaw, ok := msg[addrField]
	if !ok {
		return
	}
	addr, ok := addrRaw.(string)
	if !ok {
		return
	}
	promLabels := make(prometheus.Labels, 0)
	for _, label := range d.Config.IpInfoLabels {
		if value, ok := msg[prefix+label]; ok {
			promLabels[label] = fmt.Sprint(value)
		} else {
			// Must set a value otherwise it gets angry
			promLabels[label] = ""
		}
	}

	metric, loaded := d.ipInfoStore.Load(addr)
	if !loaded {
		metric = &prometheusDestinationIpInfo{
			labels: promLabels,
		}
		d.ipInfoStore.Store(addr, metric)
	}
	metric.lastReceive = time.Now()
}

func (d *PrometheusDestination) Describe(ch chan<- *prometheus.Desc) {
	if d.byteCounterDesc != nil {
		ch <- d.byteCounterDesc
	}
	if d.packetCounterDesc != nil {
		ch <- d.packetCounterDesc
	}
	if d.flowCounterDesc != nil {
		ch <- d.flowCounterDesc
	}
	if d.flowDurationDesc != nil {
		ch <- d.flowDurationDesc
	}
	if d.ipInfoDesc != nil {
		ch <- d.ipInfoDesc
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

	if d.ipInfoDesc != nil {
		d.ipInfoStore.Range(func(addr string, metric *prometheusDestinationIpInfo) bool {
			labelValues := make([]string, 0)
			for _, key := range d.Config.IpInfoLabels {
				labelValues = append(labelValues, metric.labels[key])
			}
			ch <- prometheus.MustNewConstMetric(d.ipInfoDesc, prometheus.GaugeValue, 1.0, labelValues...)
			return true
		})
	}
}

func (d *PrometheusDestination) metricStoreGC() {
	storeSize := 0
	evictions := 0
	d.metricStore.Range(func(hash uint64, metric *prometheusDestinationMetric) bool {
		threshold := time.Now().Add(-d.Config.VisibilityTimeout)
		if metric.lastReceive.Before(threshold) {
			d.metricStore.Delete(hash)
			evictions++
			MetricPrometheusStoreEvictionCount.WithLabelValues("metrics").Inc()
		} else {
			storeSize++
		}
		return true
	})
	MetricPrometheusStoreLastEvictionCount.WithLabelValues("metrics").Set(float64(evictions))
	MetricPrometheusStoreSize.WithLabelValues("metrics").Set(float64(storeSize))
}

func (d *PrometheusDestination) startMetricStoreGC() {
	go func() {
		MetricPrometheusStoreEvictionCount.WithLabelValues("metrics").Add(0.0)
		MetricPrometheusStoreLastEvictionCount.WithLabelValues("metrics").Set(0.0)
		MetricPrometheusStoreSize.WithLabelValues("metrics").Set(0.0)
		for {
			start := time.Now()
			d.metricStoreGC()
			duration := time.Since(start)
			MetricPrometheusStoreLastGCSeconds.WithLabelValues("metrics").Set(duration.Seconds())
			MetricPrometheusStoreGCSeconds.WithLabelValues("metrics").Observe(duration.Seconds())
			timeUntilNext := d.Config.GCInterval - duration
			if timeUntilNext > 0 {
				time.Sleep(timeUntilNext)
			}
		}
	}()
}

func (d *PrometheusDestination) ipStoreGC() {
	storeSize := 0
	evictions := 0
	d.ipInfoStore.Range(func(addr string, metric *prometheusDestinationIpInfo) bool {
		threshold := time.Now().Add(-d.Config.VisibilityTimeout)
		if metric.lastReceive.Before(threshold) {
			d.ipInfoStore.Delete(addr)
			evictions++
			MetricPrometheusStoreEvictionCount.WithLabelValues("ipinfo").Inc()
		} else {
			storeSize++
		}
		return true
	})
	MetricPrometheusStoreLastEvictionCount.WithLabelValues("ipinfo").Set(float64(evictions))
	MetricPrometheusStoreSize.WithLabelValues("ipinfo").Set(float64(storeSize))
}

func (d *PrometheusDestination) startIpStoreGC() {
	go func() {
		MetricPrometheusStoreEvictionCount.WithLabelValues("ipinfo").Add(0.0)
		MetricPrometheusStoreLastEvictionCount.WithLabelValues("ipinfo").Set(0.0)
		MetricPrometheusStoreSize.WithLabelValues("ipinfo").Set(0.0)
		for {
			start := time.Now()
			d.ipStoreGC()
			duration := time.Since(start)
			MetricPrometheusStoreLastGCSeconds.WithLabelValues("ipinfo").Set(duration.Seconds())
			MetricPrometheusStoreGCSeconds.WithLabelValues("ipinfo").Observe(duration.Seconds())
			timeUntilNext := d.Config.GCInterval - duration
			if timeUntilNext > 0 {
				time.Sleep(timeUntilNext)
			}
		}
	}()
}
