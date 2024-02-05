package enricher_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/sapslaj/morbius/enricher"
)

func TestNetDBEnricher(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		skip   string
		config *enricher.NetDBEnricherConfig
		input  map[string]interface{}
		want   map[string]interface{}
	}{
		"does not modify message is enriched field is not defined": {
			input: map[string]interface{}{"other": 69},
			want:  map[string]interface{}{"other": 69},
		},
		"sets protocol_name if protocol number is known": {
			input: map[string]interface{}{"proto": 1},
			want:  map[string]interface{}{"proto": 1, "protocol_name": "icmp"},
		},
		"does not set protocol_name if protocol number is unknown": {
			input: map[string]interface{}{"proto": 69},
			want:  map[string]interface{}{"proto": 69},
		},
		"sets ethernet_type_name if EtherType is known": {
			input: map[string]interface{}{"ethernet_type": 2048},
			want:  map[string]interface{}{"ethernet_type": 2048, "ethernet_type_name": "IPv4"},
		},
		"does not set ethernet_type_name if EtherType is unknown": {
			input: map[string]interface{}{"ethernet_type": 69},
			want:  map[string]interface{}{"ethernet_type": 69},
		},
		"handles proto_encap": {
			input: map[string]interface{}{"proto": 47, "proto_encap": 1},
			want:  map[string]interface{}{"proto": 47, "protocol_name": "gre", "proto_encap": 1, "protocol_encap_name": "icmp"},
		},
		"handles ethernet_type_encap": {
			input: map[string]interface{}{"ethernet_type": 2048, "ethernet_type_encap": 34525},
			want:  map[string]interface{}{"ethernet_type": 2048, "ethernet_type_name": "IPv4", "ethernet_type_encap": 34525, "ethernet_type_encap_name": "IPv6"},
		},
		"sets src_service_name": {
			input: map[string]interface{}{"proto": 6, "src_port": 80},
			want:  map[string]interface{}{"proto": 6, "protocol_name": "tcp", "src_port": 80, "src_service_name": "http", "service_name": "http"},
		},
		"sets dst_service_name": {
			input: map[string]interface{}{"proto": 6, "dst_port": 80},
			want:  map[string]interface{}{"proto": 6, "protocol_name": "tcp", "dst_port": 80, "dst_service_name": "http", "service_name": "http"},
		},
		"sets service_name": {
			input: map[string]interface{}{"proto": 6, "src_port": 42069, "dst_port": 80},
			want:  map[string]interface{}{"proto": 6, "protocol_name": "tcp", "src_port": 42069, "dst_port": 80, "dst_service_name": "http", "service_name": "http"},
		},
		"supports custom protocols": {
			config: &enricher.NetDBEnricherConfig{
				Protocols: &enricher.NetDBEnricherConfigConfig{
					SourceInline: []string{"nice 69"},
				},
			},
			input: map[string]interface{}{"proto": 69},
			want:  map[string]interface{}{"proto": 69, "protocol_name": "nice"},
		},
		"supports custom services": {
			config: &enricher.NetDBEnricherConfig{
				Services: &enricher.NetDBEnricherConfigConfig{
					SourceInline: []string{"nice 69/tcp"},
				},
			},
			input: map[string]interface{}{"proto": 6, "src_port": 42069, "dst_port": 69},
			want:  map[string]interface{}{"proto": 6, "protocol_name": "tcp", "src_port": 42069, "dst_port": 69, "dst_service_name": "nice", "service_name": "nice"},
		},
		"supports custom services with custom protocols": {
			config: &enricher.NetDBEnricherConfig{
				Protocols: &enricher.NetDBEnricherConfigConfig{
					SourceInline: []string{"nice 69"},
				},
				Services: &enricher.NetDBEnricherConfigConfig{
					SourceInline: []string{"nice 69/nice"},
				},
			},
			input: map[string]interface{}{"proto": 69, "src_port": 42069, "dst_port": 69},
			want:  map[string]interface{}{"proto": 69, "protocol_name": "nice", "src_port": 42069, "dst_port": 69, "dst_service_name": "nice", "service_name": "nice"},
		},
		"supports custom EtherTypes": {
			config: &enricher.NetDBEnricherConfig{
				EtherTypes: &enricher.NetDBEnricherConfigConfig{
					SourceInline: []string{"nice 69 # nice"},
				},
			},
			input: map[string]interface{}{"ethernet_type": 0x69},
			want:  map[string]interface{}{"ethernet_type": 0x69, "ethernet_type_name": "nice"},
		},
		"supports name aliases": {
			config: &enricher.NetDBEnricherConfig{
				EtherTypes: &enricher.NetDBEnricherConfigConfig{
					BuiltIn: true,
					NameAliases: map[string]string{
						"IPv4": "custom EtherType alias",
					},
				},
				Protocols: &enricher.NetDBEnricherConfigConfig{
					BuiltIn: true,
					NameAliases: map[string]string{
						"tcp": "custom protocol alias",
					},
				},
				Services: &enricher.NetDBEnricherConfigConfig{
					BuiltIn: true,
					NameAliases: map[string]string{
						"http": "custom service alias",
					},
				},
			},
			input: map[string]interface{}{"ethernet_type": 2048, "proto": 6, "src_port": 80, "dst_port": 80},
			want: map[string]interface{}{
				"ethernet_type":      2048,
				"ethernet_type_name": "custom EtherType alias",
				"proto":              6,
				"protocol_name":      "custom protocol alias",
				"src_port":           80,
				"src_service_name":   "custom service alias",
				"dst_port":           80,
				"dst_service_name":   "custom service alias",
				"service_name":       "custom service alias",
			},
		},
		"aliases can use aliases": {
			config: &enricher.NetDBEnricherConfig{
				EtherTypes: &enricher.NetDBEnricherConfigConfig{
					BuiltIn: true,
					NameAliases: map[string]string{
						"802.1q": "an alias of an alias!",
					},
				},
			},
			input: map[string]interface{}{"ethernet_type": 33024},
			want:  map[string]interface{}{"ethernet_type": 33024, "ethernet_type_name": "an alias of an alias!"},
		},
		"uses extra built-in EtherTypes": {
			input: map[string]interface{}{"ethernet_type": 0x88bf},
			want:  map[string]interface{}{"ethernet_type": 0x88bf, "ethernet_type_name": "mikrotik-romon"},
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
			e := enricher.NewNetDBEnricher(tc.config)
			got := e.Process(tc.input)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Logf("\"%s\":\n%s", name, diff)
				t.Fail()
			}
		})
	}
}
