package enricher_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sapslaj/morbius/enricher"
)

func TestProtonamesEnricher(t *testing.T) {
	t.Parallel()
	type test struct {
		input map[string]interface{}
		want  map[string]interface{}
	}

	tests := map[string]test{
		"does not modify message if protocol not defined": {
			input: map[string]interface{}{"other": 69},
			want:  map[string]interface{}{"other": 69},
		},
		"adds a protocol name for a known protocol": {
			input: map[string]interface{}{"proto": 1},
			want:  map[string]interface{}{"proto": 1, "protocol_name": "ICMP"},
		},
		"does not add protocol name if protocol number is unknown": {
			input: map[string]interface{}{"proto": 69},
			want:  map[string]interface{}{"proto": 69},
		},
		"adds an etype name for a known etype": {
			input: map[string]interface{}{"ethernet_type": 2048},
			want:  map[string]interface{}{"ethernet_type": 2048, "ethernet_type_name": "IPv4"},
		},
		"does not add etype name if etype is unknown": {
			input: map[string]interface{}{"ethernet_type": 69},
			want:  map[string]interface{}{"ethernet_type": 69},
		},
		"handles proto_encap": {
			input: map[string]interface{}{"proto": 47, "proto_encap": 1},
			want:  map[string]interface{}{"proto": 47, "protocol_name": "GRE", "proto_encap": 1, "protocol_encap_name": "ICMP"},
		},
		"handles ethernet_type_encap": {
			input: map[string]interface{}{"ethernet_type": 2048, "ethernet_type_encap": 34525},
			want:  map[string]interface{}{"ethernet_type": 2048, "ethernet_type_name": "IPv4", "ethernet_type_encap": 34525, "ethernet_type_encap_name": "IPv6"},
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			pe := enricher.NewProtonamesEnricher(nil)
			got := pe.Process(tc.input)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf("\"%s\":\n%s", name, diff)
			}
		})
	}
}
