package enricher

import (
	"net"
	"sort"

	lru "github.com/hashicorp/golang-lru/v2"
)

type RDNSEnricher struct {
	cache *lru.Cache[string, string]
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
		msg[targetField] = value
		return msg
	}
	names, err := net.LookupAddr(addr)
	if err != nil {
		// log.Printf("error with net.LookupAddr: %v", err)
		return msg
	}
	if len(names) == 0 {
		return msg
	}
	sort.Strings(names)
	value = names[0]
	e.cache.Add(addr, value)
	msg[targetField] = value
	return msg
}
