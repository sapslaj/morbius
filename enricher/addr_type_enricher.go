package enricher

import (
	"log"
	"net/netip"
)

type AddrTypeEnricherConfig struct {
}

type addrTypeEnricherPrefix struct {
	cidr   string
	prefix netip.Prefix
	typ    string
}

type AddrTypeEnricher struct {
	Config   *AddrTypeEnricherConfig
	prefixes []addrTypeEnricherPrefix
}

func NewAddrTypeEnricher(config *AddrTypeEnricherConfig) AddrTypeEnricher {
	if config == nil {
		config = &AddrTypeEnricherConfig{}
	}
	e := AddrTypeEnricher{
		Config: config,
	}

	e.addPrefix("255.255.255.255/32", "limited-broadcast")
	e.addPrefix("192.0.0.0/29", "ds-lite")
	e.addPrefix("192.0.0.0/24", "ietf-protocol-assignments")
	e.addPrefix("192.0.2.0/24", "documentation")
	e.addPrefix("192.88.99.0/24", "6to4-relay-anycast")
	e.addPrefix("198.51.100.0/24", "documentation")
	e.addPrefix("203.0.113.0/24", "documentation")
	e.addPrefix("233.252.0.0/24", "mcast-test-net")
	e.addPrefix("192.168.0.0/16", "private")
	e.addPrefix("169.254.0.0/16", "link-local")
	e.addPrefix("198.18.0.0/15", "benchmarking")
	e.addPrefix("172.16.0.0/12", "private")
	e.addPrefix("100.0.0.0/10", "cgn")
	e.addPrefix("127.0.0.0/8", "loopback")
	e.addPrefix("10.0.0.0/8", "private")
	e.addPrefix("0.0.0.0/8", "current-network")
	e.addPrefix("224.0.0.0/4", "multicast")
	e.addPrefix("240.0.0.0/4", "reserved")
	e.addPrefix("::/128", "unspecified")
	e.addPrefix("::1/128", "loopback")
	e.addPrefix("100::/64", "discard-only")
	e.addPrefix("::ffff:0:0/96", "ipv4-mapped")
	e.addPrefix("::ffff:0:0:0/96", "ipv4-translated")
	e.addPrefix("64:ff9b::/96", "ipv4-ipv6-translation-global")
	e.addPrefix("64:ff9b:1::/48", "ipv4-ipv6-translation-private")
	e.addPrefix("2001:2::/48", "benchmarking")
	e.addPrefix("2001:0000::/32", "teredo")
	e.addPrefix("2001:db8::/32", "documentation")
	e.addPrefix("2001:20::/28", "orchidv2")
	e.addPrefix("2001:10::/28", "orchid")
	e.addPrefix("2001::/23", "ietf-protocol-assignments")
	e.addPrefix("2002::/16", "6to4")
	e.addPrefix("fe80::/10", "link-local")
	e.addPrefix("::/8", "reserved-by-ietf")
	e.addPrefix("ff00::/8", "multicast")
	e.addPrefix("fc00::/7", "ula")
	return e
}

func (e *AddrTypeEnricher) addPrefix(cidr, typ string) {
	prefix, err := netip.ParsePrefix(cidr)
	if err != nil {
		panic(err)
	}
	e.prefixes = append(e.prefixes, addrTypeEnricherPrefix{
		cidr:   cidr,
		prefix: prefix,
		typ:    typ,
	})
}

func (e *AddrTypeEnricher) Process(msg map[string]interface{}) map[string]interface{} {
	msg = e.add(msg, "src_addr", "src_addr_type")
	msg = e.add(msg, "dst_addr", "dst_addr_type")
	msg = e.add(msg, "src_addr_encap", "src_addr_encap_type")
	msg = e.add(msg, "dst_addr_encap", "dst_addr_encap_type")
	return msg
}

func (e *AddrTypeEnricher) add(msg map[string]interface{}, originalField string, targetField string) map[string]interface{} {
	addrRaw, ok := msg[originalField]
	if !ok {
		return msg
	}
	addr, ok := addrRaw.(string)
	if !ok {
		return msg
	}

	netipAddr, err := netip.ParseAddr(addr)
	if err != nil {
		log.Printf("error parsing IP address `%s`: %v", addr, err)
		return msg
	}

	typ := "global"
	for _, p := range e.prefixes {
		if p.prefix.Contains(netipAddr) {
			typ = p.typ
			break
		}
	}
	msg[targetField] = typ

	return msg
}
