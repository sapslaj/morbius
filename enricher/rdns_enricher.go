package enricher

import (
	"net"
	"sort"
)

type RDNSEnricher struct {
}

func NewRDNSEnricher() RDNSEnricher {
	return RDNSEnricher{}
}

func (e *RDNSEnricher) Process(msg map[string]interface{}) map[string]interface{} {
	msg = e.add(msg, "src_addr", "src_hostname")
	msg = e.add(msg, "dst_addr", "dst_hostname")
	return msg
}

func (e *RDNSEnricher) add(msg map[string]interface{}, originalField string, targetField string) map[string]interface{} {
	addr, ok := msg[originalField]
	if !ok {
		return msg
	}
	names, err := net.LookupAddr(addr.(string))
	if err != nil {
		// log.Printf("error with net.LookupAddr: %v", err)
		return msg
	}
	if len(names) == 0 {
		return msg
	}
	sort.Strings(names)
	msg[targetField] = names[0]
	return msg
}
