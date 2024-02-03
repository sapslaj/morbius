package enricher_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/sapslaj/morbius/enricher"
)

func TestAddrTypeEnricher(t *testing.T) {
	t.Parallel()
	type test struct {
		skip  string
		input map[string]interface{}
		want  map[string]interface{}
	}

	tests := map[string]test{
		"does not modify message if an address field is not defined": {
			input: map[string]interface{}{"other": 69},
			want:  map[string]interface{}{"other": 69},
		},
		"adds src_addr_type when src_addr is set": {
			input: map[string]interface{}{"src_addr": "1.1.1.1"},
			want:  map[string]interface{}{"src_addr": "1.1.1.1", "src_addr_type": "global"},
		},
		"adds dst_addr_type when dst_addr is set": {
			input: map[string]interface{}{"dst_addr": "1.1.1.1"},
			want:  map[string]interface{}{"dst_addr": "1.1.1.1", "dst_addr_type": "global"},
		},
		"adds src_addr_encap_type when src_addr_encap is set": {
			input: map[string]interface{}{"src_addr_encap": "1.1.1.1"},
			want:  map[string]interface{}{"src_addr_encap": "1.1.1.1", "src_addr_encap_type": "global"},
		},
		"adds dst_addr_encap_type when dst_addr_encap is set": {
			input: map[string]interface{}{"dst_addr_encap": "1.1.1.1"},
			want:  map[string]interface{}{"dst_addr_encap": "1.1.1.1", "dst_addr_encap_type": "global"},
		},
		"does nothing on empty string for address": {
			input: map[string]interface{}{"src_addr": ""},
			want:  map[string]interface{}{"src_addr": ""},
		},
		"correctly classifies a global IPv4 address": {
			input: map[string]interface{}{"src_addr": "1.1.1.1"},
			want:  map[string]interface{}{"src_addr": "1.1.1.1", "src_addr_type": "global"},
		},
		"correctly classifies a 'this host on this network' IPv4 address": {
			input: map[string]interface{}{"src_addr": "0.0.0.1"},
			want:  map[string]interface{}{"src_addr": "0.0.0.1", "src_addr_type": "current-network"},
		},
		"correctly classifies a shared/CGN IPv4 address": {
			input: map[string]interface{}{"src_addr": "100.0.4.20"},
			want:  map[string]interface{}{"src_addr": "100.0.4.20", "src_addr_type": "cgn"},
		},
		"correctly classifies an RFC1918 private class A IPv4 address": {
			input: map[string]interface{}{"src_addr": "10.69.4.20"},
			want:  map[string]interface{}{"src_addr": "10.69.4.20", "src_addr_type": "private"},
		},
		"correctly classifies an RFC1918 private class B IPv4 address": {
			input: map[string]interface{}{"src_addr": "172.17.4.20"},
			want:  map[string]interface{}{"src_addr": "172.17.4.20", "src_addr_type": "private"},
		},
		"correctly classifies an RFC1918 private class C IPv4 address": {
			input: map[string]interface{}{"src_addr": "192.168.4.20"},
			want:  map[string]interface{}{"src_addr": "192.168.4.20", "src_addr_type": "private"},
		},
		"correctly classifies a loopback IPv4 address": {
			input: map[string]interface{}{"src_addr": "127.69.4.20"},
			want:  map[string]interface{}{"src_addr": "127.69.4.20", "src_addr_type": "loopback"},
		},
		"correctly classifies a link-local IPv4 address": {
			input: map[string]interface{}{"src_addr": "169.254.4.20"},
			want:  map[string]interface{}{"src_addr": "169.254.4.20", "src_addr_type": "link-local"},
		},
		"correctly classifies a IETF protocol assignment IPv4 address": {
			input: map[string]interface{}{"src_addr": "192.0.0.69"},
			want:  map[string]interface{}{"src_addr": "192.0.0.69", "src_addr_type": "ietf-protocol-assignments"},
		},
		"correctly classifies a TEST-NET-1 IPv4 address": {
			input: map[string]interface{}{"src_addr": "192.0.2.69"},
			want:  map[string]interface{}{"src_addr": "192.0.2.69", "src_addr_type": "documentation"},
		},
		"correctly classifies a 6to4 relay anycast IPv4 address": {
			input: map[string]interface{}{"src_addr": "192.88.99.69"},
			want:  map[string]interface{}{"src_addr": "192.88.99.69", "src_addr_type": "6to4-relay-anycast"},
		},
		"correctly classifies a benchmarking address": {
			input: map[string]interface{}{"src_addr": "198.18.4.20"},
			want:  map[string]interface{}{"src_addr": "198.18.4.20", "src_addr_type": "benchmarking"},
		},
		"correctly classifies a TEST-NET-2 IPv4 address": {
			input: map[string]interface{}{"src_addr": "198.51.100.69"},
			want:  map[string]interface{}{"src_addr": "198.51.100.69", "src_addr_type": "documentation"},
		},
		"correctly classifies a TEST-NET-3 IPv4 address": {
			input: map[string]interface{}{"src_addr": "203.0.113.69"},
			want:  map[string]interface{}{"src_addr": "203.0.113.69", "src_addr_type": "documentation"},
		},
		"correctly classifies a multicast IPv4 address": {
			input: map[string]interface{}{"src_addr": "224.69.4.20"},
			want:  map[string]interface{}{"src_addr": "224.69.4.20", "src_addr_type": "multicast"},
		},
		"correctly classifies a MCAST-TEST-NET IPv4 address": {
			input: map[string]interface{}{"src_addr": "233.252.0.69"},
			want:  map[string]interface{}{"src_addr": "233.252.0.69", "src_addr_type": "mcast-test-net"},
		},
		"correctly classifies a reserved IPv4 address": {
			input: map[string]interface{}{"src_addr": "240.69.4.20"},
			want:  map[string]interface{}{"src_addr": "240.69.4.20", "src_addr_type": "reserved"},
		},
		"correctly classifies a limited broadcast IPv4 address": {
			input: map[string]interface{}{"src_addr": "255.255.255.255"},
			want:  map[string]interface{}{"src_addr": "255.255.255.255", "src_addr_type": "limited-broadcast"},
		},
		"correctly classifies a loopback IPv6 address": {
			input: map[string]interface{}{"src_addr": "::1"},
			want:  map[string]interface{}{"src_addr": "::1", "src_addr_type": "loopback"},
		},
		"correctly classifies an unspecified IPv6 address": {
			input: map[string]interface{}{"src_addr": "::"},
			want:  map[string]interface{}{"src_addr": "::", "src_addr_type": "unspecified"},
		},
		"correctly classifies a reserved IPv6 address": {
			input: map[string]interface{}{"src_addr": "64::1"},
			want:  map[string]interface{}{"src_addr": "64::1", "src_addr_type": "reserved-by-ietf"},
		},
		"correctly classifies a global IPv4-IPv6 Translation IPv6 address": {
			input: map[string]interface{}{"src_addr": "64:ff9b::1"},
			want:  map[string]interface{}{"src_addr": "64:ff9b::1", "src_addr_type": "ipv4-ipv6-translation-global"},
		},
		"correctly classifies a private IPv4-IPv6 Translation IPv6 address": {
			input: map[string]interface{}{"src_addr": "64:ff9b:1::1"},
			want:  map[string]interface{}{"src_addr": "64:ff9b:1::1", "src_addr_type": "ipv4-ipv6-translation-private"},
		},
		"correctly classifies an IPv4-mapped IPv6 address": {
			input: map[string]interface{}{"src_addr": "::ffff:0:1"},
			want:  map[string]interface{}{"src_addr": "::ffff:0:1", "src_addr_type": "ipv4-mapped"},
		},
		"correctly classifies an IPv4-translated IPv6 address": {
			input: map[string]interface{}{"src_addr": "::ffff:0:0:1"},
			want:  map[string]interface{}{"src_addr": "::ffff:0:0:1", "src_addr_type": "ipv4-translated"},
		},
		"correctly classifies a discard-only IPv6 address": {
			input: map[string]interface{}{"src_addr": "100::1"},
			want:  map[string]interface{}{"src_addr": "100::1", "src_addr_type": "discard-only"},
		},
		"correctly classifies a Teredo IPv6 address": {
			input: map[string]interface{}{"src_addr": "2001::1"},
			want:  map[string]interface{}{"src_addr": "2001::1", "src_addr_type": "teredo"},
		},
		"correctly classifies a benchmarking IPv6 address": {
			input: map[string]interface{}{"src_addr": "2001:2::1"},
			want:  map[string]interface{}{"src_addr": "2001:2::1", "src_addr_type": "benchmarking"},
		},
		"correctly classifies a ORCHID IPv6 address": {
			input: map[string]interface{}{"src_addr": "2001:10::1"},
			want:  map[string]interface{}{"src_addr": "2001:10::1", "src_addr_type": "orchid"},
		},
		"correctly classifies a ORCHIDv2 IPv6 address": {
			input: map[string]interface{}{"src_addr": "2001:20::1"},
			want:  map[string]interface{}{"src_addr": "2001:20::1", "src_addr_type": "orchidv2"},
		},
		"correctly classifies a documentation IPv6 address": {
			input: map[string]interface{}{"src_addr": "2001:db8::1"},
			want:  map[string]interface{}{"src_addr": "2001:db8::1", "src_addr_type": "documentation"},
		},
		"correctly classifies a 6to4 IPv6 address": {
			input: map[string]interface{}{"src_addr": "2002::1"},
			want:  map[string]interface{}{"src_addr": "2002::1", "src_addr_type": "6to4"},
		},
		"correctly classifies a ULA IPv6 address": {
			input: map[string]interface{}{"src_addr": "fd00::1"},
			want:  map[string]interface{}{"src_addr": "fd00::1", "src_addr_type": "ula"},
		},
		"correctly classifies a link-local IPv6 address": {
			input: map[string]interface{}{"src_addr": "fe80::1"},
			want:  map[string]interface{}{"src_addr": "fe80::1", "src_addr_type": "link-local"},
		},
		"correctly classifies a multicast IPv6 address": {
			input: map[string]interface{}{"src_addr": "ff00::1"},
			want:  map[string]interface{}{"src_addr": "ff00::1", "src_addr_type": "multicast"},
		},
		"correctly classifies an IETF protocol assignments IPv6 address": {
			skip:  "collision with Teredo range",
			input: map[string]interface{}{"src_addr": "2001::1"},
			want:  map[string]interface{}{"src_addr": "2001::1", "src_addr_type": "ietf-protocol-assignments"},
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if tc.input == nil {
				t.Logf("\"%s\": skip (unimplmented)", name)
				return
			}
			if tc.skip != "" {
				t.Logf("\"%s\": skip (%s)", name, tc.skip)
				return
			}
			e := enricher.NewAddrTypeEnricher(&enricher.AddrTypeEnricherConfig{})
			got := e.Process(tc.input)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Logf("\"%s\":\n%s", name, diff)
				t.Fail()
			}
		})
	}
}
