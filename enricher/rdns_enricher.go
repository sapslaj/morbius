package enricher

import (
	"net"
	"sort"
	"sync"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/prometheus/client_golang/prometheus"
)

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

type RDNSEnricherConfig struct {
	EnableCache bool
	CacheSize   int
	CacheOnly   bool
}

type RDNSEnricher struct {
	Config            *RDNSEnricherConfig
	cache             *lru.Cache[string, string]
	cacheLookupStatus sync.Map
}

func NewRDNSEnricher(config *RDNSEnricherConfig) RDNSEnricher {
	if config == nil {
		config = &RDNSEnricherConfig{}
	}
	if config.EnableCache && config.CacheSize == 0 {
		config.CacheSize = 128
	}
	var cache *lru.Cache[string, string]
	var err error
	if config.EnableCache {
		cache, err = lru.New[string, string](config.CacheSize)
		if err != nil {
			panic(err)
		}
	}
	return RDNSEnricher{
		Config: config,
		cache:  cache,
	}
}

func (e *RDNSEnricher) Process(msg map[string]interface{}) map[string]interface{} {
	if e.Config.EnableCache {
		MetricRDNSCacheSize.Set(float64(e.cache.Len()))
	}
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
	var value string

	if e.Config.EnableCache {
		value, ok := e.cache.Get(addr)
		if ok {
			MetricRDNSCacheHits.Inc()
			if value != "" {
				msg[targetField] = value
			}
			return msg
		}
		MetricRDNSCacheMisses.Inc()
	}

	if e.Config.CacheOnly {
		go func(addr string) {
			_, inProgress := e.cacheLookupStatus.Load(addr)
			if !inProgress {
				e.cacheLookupStatus.Store(addr, true)
				e.lookup(addr)
				e.cacheLookupStatus.Delete(addr)
			}
		}(addr)
		return msg
	}

	value, ok = e.lookup(addr)
	if !ok {
		return msg
	}
	if value != "" {
		msg[targetField] = value
	}
	return msg
}

func (e *RDNSEnricher) lookup(addr string) (string, bool) {
	names, err := net.LookupAddr(addr)
	if err != nil {
		MetricRDNSLookups.With(prometheus.Labels{"status": "error"}).Inc()
		// TODO: only return ok = true on general lookup failures
	} else if len(names) == 0 {
		MetricRDNSLookups.With(prometheus.Labels{"status": "empty"}).Inc()
	} else {
		MetricRDNSLookups.With(prometheus.Labels{"status": "success"}).Inc()
	}

	value := ""
	if len(names) > 0 {
		sort.Strings(names)
		value = names[0]
	}
	if e.Config.EnableCache {
		e.cache.Add(addr, value)
	}
	return value, true
}
