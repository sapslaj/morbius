package enricher

import (
	"net"
	"sort"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/prometheus/client_golang/prometheus"
)

type RDNSEnricher struct {
	cache *lru.Cache[string, string]
}

var (
	MetricRDNSCacheSize = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "rdns_cache_size",
			Help: "size of RDNS enricher LRU cache",
		},
	)
	MetricRDNSCacheHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "rdns_cache_hits",
			Help: "Number of RDNS enricher LRU cache hits",
		},
	)
	MetricRDNSCacheMisses = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "rdns_cache_misses",
			Help: "Number of RDNS enricher LRU cache misses",
		},
	)
	MetricRDNSLookups = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rdns_lookups",
			Help: "Number of RDNS enricher lookups",
		},
		[]string{"status"},
	)
)

func init() {
	prometheus.MustRegister(MetricRDNSCacheSize)
	prometheus.MustRegister(MetricRDNSCacheHits)
	prometheus.MustRegister(MetricRDNSCacheMisses)
	prometheus.MustRegister(MetricRDNSLookups)
}

func NewRDNSEnricher() RDNSEnricher {
	cache, err := lru.New[string, string](2048)
	if err != nil {
		panic(err)
	}
	return RDNSEnricher{
		cache: cache,
	}
}

func (e *RDNSEnricher) Process(msg map[string]interface{}) map[string]interface{} {
	MetricRDNSCacheSize.Set(float64(e.cache.Len()))
	msg = e.add(msg, "src_addr", "src_hostname")
	msg = e.add(msg, "dst_addr", "dst_hostname")
	return msg
}

func (e *RDNSEnricher) add(msg map[string]interface{}, originalField string, targetField string) map[string]interface{} {
	addrRaw, ok := msg[originalField]
	if !ok {
		return msg
	}
	addr := addrRaw.(string)
	value, ok := e.cache.Get(addr)
	if ok {
		MetricRDNSCacheHits.Inc()
		msg[targetField] = value
		return msg
	}
	MetricRDNSCacheMisses.Inc()
	names, err := net.LookupAddr(addr)
	if err != nil {
		MetricRDNSLookups.With(prometheus.Labels{"status": "error"}).Inc()
		return msg
	}
	if len(names) == 0 {
		MetricRDNSLookups.With(prometheus.Labels{"status": "empty"}).Inc()
		return msg
	}
	MetricRDNSLookups.With(prometheus.Labels{"status": "success"}).Inc()
	sort.Strings(names)
	value = names[0]
	e.cache.Add(addr, value)
	msg[targetField] = value
	return msg
}
