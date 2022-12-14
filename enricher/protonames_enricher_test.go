package enricher_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sapslaj/morbius/enricher"
)

func TestProtonamesEnricher(t *testing.T) {
	type test struct {
		desc  string
		input map[string]interface{}
		want  map[string]interface{}
	}

	pe := enricher.NewProtonamesEnricher(nil)

	tests := []test{
		{
			desc:  "Does not modify message if protocol not defined",
			input: map[string]interface{}{"other": 69},
			want:  map[string]interface{}{"other": 69},
		},
		{
			desc:  "Adds a protocol name for a known protocol",
			input: map[string]interface{}{"proto": 1},
			want:  map[string]interface{}{"proto": 1, "protocol_name": "ICMP"},
		},
		{
			desc:  "Does not add protocol name if protocol number is unknown",
			input: map[string]interface{}{"proto": 69},
			want:  map[string]interface{}{"proto": 69},
		},
		{
			desc:  "Adds an etype name for a known etype",
			input: map[string]interface{}{"ethernet_type": 2048},
			want:  map[string]interface{}{"ethernet_type": 2048, "ethernet_type_name": "IPv4"},
		},
		{
			desc:  "Does not add etype name if etype is unknown",
			input: map[string]interface{}{"ethernet_type": 69},
			want:  map[string]interface{}{"ethernet_type": 69},
		},
		{
			desc:  "Handles proto_encap",
			input: map[string]interface{}{"proto": 47, "proto_encap": 1},
			want:  map[string]interface{}{"proto": 47, "protocol_name": "GRE", "proto_encap": 1, "protocol_encap_name": "ICMP"},
		},
		{
			desc:  "Handles ethernet_type_encap",
			input: map[string]interface{}{"ethernet_type": 2048, "ethernet_type_encap": 34525},
			want:  map[string]interface{}{"ethernet_type": 2048, "ethernet_type_name": "IPv4", "ethernet_type_encap": 34525, "ethernet_type_encap_name": "IPv6"},
		},
	}

	for _, tc := range tests {
		got := pe.Process(tc.input)
		if diff := cmp.Diff(tc.want, got); diff != "" {
			t.Fatalf("\"%s\":\n%s", tc.desc, diff)
		}
	}
}
