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
			input: map[string]interface{}{"protocol": 69},
			want:  map[string]interface{}{"protocol": 69},
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
	}

	for _, tc := range tests {
		got := pe.Process(tc.input)
		if !cmp.Equal(tc.want, got) {
			t.Fatalf("\"%s\": expected: %v, got: %v", tc.desc, tc.want, got)
		}
	}
}
