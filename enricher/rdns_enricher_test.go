package enricher_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sapslaj/morbius/enricher"
)

func TestRDNSEnricher(t *testing.T) {
	type test struct {
		desc  string
		input map[string]interface{}
		want  map[string]interface{}
	}

	pe := enricher.NewRDNSEnricher(nil)

	tests := []test{
		{
			desc:  "Does not modify message if an address field is not defined",
			input: map[string]interface{}{"other": 69},
			want:  map[string]interface{}{"other": 69},
		},
		{
			desc:  "Adds hostname field when IPv4 is resolvable",
			input: map[string]interface{}{"src_addr": "1.1.1.1"},
			want:  map[string]interface{}{"src_addr": "1.1.1.1", "src_hostname": "one.one.one.one."},
		},
		{
			desc:  "Omits hostname field when IPv4 is not resolvable",
			input: map[string]interface{}{"src_addr": "30.1.1.1"},
			want:  map[string]interface{}{"src_addr": "30.1.1.1"},
		},
		{
			desc:  "Adds hostname field when IPv6 is resolvable",
			input: map[string]interface{}{"src_addr": "2606:4700:4700::1111"},
			want:  map[string]interface{}{"src_addr": "2606:4700:4700::1111", "src_hostname": "one.one.one.one."},
		},
		{
			desc:  "Omits hostname field when IPv6 is not resolvable",
			input: map[string]interface{}{"src_addr": "2001::404"},
			want:  map[string]interface{}{"src_addr": "2001::404"},
		},
		{
			desc:  "Adds src_hostname when src_addr is set",
			input: map[string]interface{}{"src_addr": "1.1.1.1"},
			want:  map[string]interface{}{"src_addr": "1.1.1.1", "src_hostname": "one.one.one.one."},
		},
		{
			desc:  "Adds dst_hostname when dst_addr is set",
			input: map[string]interface{}{"dst_addr": "1.1.1.1"},
			want:  map[string]interface{}{"dst_addr": "1.1.1.1", "dst_hostname": "one.one.one.one."},
		},
		{
			desc:  "Adds src_hostname_encap when src_addr_encap is set",
			input: map[string]interface{}{"src_addr_encap": "1.1.1.1"},
			want:  map[string]interface{}{"src_addr_encap": "1.1.1.1", "src_hostname_encap": "one.one.one.one."},
		},
		{
			desc:  "Adds dst_hostname_encap when dst_addr_encap is set",
			input: map[string]interface{}{"dst_addr_encap": "1.1.1.1"},
			want:  map[string]interface{}{"dst_addr_encap": "1.1.1.1", "dst_hostname_encap": "one.one.one.one."},
		},
	}

	for _, tc := range tests {
		got := pe.Process(tc.input)
		if diff := cmp.Diff(tc.want, got); diff != "" {
			t.Fatalf("\"%s\":\n%s", tc.desc, diff)
		}
	}
}
