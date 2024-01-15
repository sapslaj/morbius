package enricher_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sapslaj/morbius/enricher"
)

func TestFieldMapperEnricher(t *testing.T) {
	tests := map[string]struct {
		input  map[string]interface{}
		want   map[string]interface{}
		config *enricher.FieldMapperEnricherConfig
	}{
		"does not modify message if config is nil": {
			input:  map[string]interface{}{"other": 69},
			want:   map[string]interface{}{"other": 69},
			config: nil,
		},
		"does not modify message if config doesn't have any mappings": {
			input: map[string]interface{}{"other": 69},
			want:  map[string]interface{}{"other": 69},
			config: &enricher.FieldMapperEnricherConfig{
				Fields: []*enricher.FieldMapperEnricherFieldConfig{},
			},
		},
		"does not modify message if config mappings don't map to fields in message": {
			input: map[string]interface{}{"other": 69},
			want:  map[string]interface{}{"other": 69},
			config: &enricher.FieldMapperEnricherConfig{
				Fields: []*enricher.FieldMapperEnricherFieldConfig{
					{
						SourceField: "bogus",
						TargetField: "bogouss",
						Mapping: map[any]any{
							"a": "b",
						},
					},
				},
			},
		},
		"basic key-value mappings": {
			input: map[string]interface{}{"src_addr": "1.1.1.1", "dst_addr": "2.2.2.2", "in_interface": 1, "out_interface": 2},
			want:  map[string]interface{}{"src_addr": "1.1.1.1", "dst_addr": "2.2.2.2", "in_interface": 1, "out_interface": 2, "in_interface_name": "igb1", "out_interface_name": "igb2"},
			config: &enricher.FieldMapperEnricherConfig{
				Fields: []*enricher.FieldMapperEnricherFieldConfig{
					{
						SourceField: "in_interface",
						TargetField: "in_interface_name",
						Mapping: map[any]any{
							0: "igb0",
							1: "igb1",
							2: "igb2",
						},
					},
					{
						SourceField: "out_interface",
						TargetField: "out_interface_name",
						Mapping: map[any]any{
							0: "igb0",
							1: "igb1",
							2: "igb2",
						},
					},
				},
			},
		},
		"key-value mapping with matching field but no matching mapping": {
			input: map[string]interface{}{"src_addr": "1.1.1.1", "dst_addr": "2.2.2.2", "sampler_address": "10.0.0.69"},
			want:  map[string]interface{}{"src_addr": "1.1.1.1", "dst_addr": "2.2.2.2", "sampler_address": "10.0.0.69"},
			config: &enricher.FieldMapperEnricherConfig{
				Fields: []*enricher.FieldMapperEnricherFieldConfig{
					{
						SourceField: "sampler_address",
						TargetField: "sampler_name",
						Mapping: map[any]any{
							"10.0.0.1": "router1",
							"10.0.0.2": "router2",
						},
					},
				},
			},
		},
		"basic template": {
			input: map[string]interface{}{"src_addr": "1.1.1.1", "dst_addr": "2.2.2.2", "in_interface": 1, "out_interface": 2},
			want:  map[string]interface{}{"src_addr": "1.1.1.1", "dst_addr": "2.2.2.2", "in_interface": 1, "out_interface": 2, "in_interface_name": "eth1", "out_interface_name": "eth2"},
			config: &enricher.FieldMapperEnricherConfig{
				Fields: []*enricher.FieldMapperEnricherFieldConfig{
					{
						SourceField: "in_interface",
						TargetField: "in_interface_name",
						Template:    `eth{{ .SourceField }}`,
					},
					{
						SourceField: "out_interface",
						TargetField: "out_interface_name",
						Template:    `eth{{ .SourceField }}`,
					},
				},
			},
		},
		"template with complex logic": {
			input: map[string]interface{}{"src_addr": "1.1.1.1", "dst_addr": "2.2.2.2", "sampler_address": "10.0.0.2", "sampler_name": "router2", "in_interface": 1, "out_interface": 2},
			want:  map[string]interface{}{"src_addr": "1.1.1.1", "dst_addr": "2.2.2.2", "sampler_address": "10.0.0.2", "sampler_name": "router2", "in_interface": 1, "out_interface": 2, "in_interface_name": "eth1", "out_interface_name": "eth1.1"},
			config: &enricher.FieldMapperEnricherConfig{
				Fields: []*enricher.FieldMapperEnricherFieldConfig{
					{
						SourceField: "sampler_address",
						TargetField: "sampler_name",
						Mapping: map[any]any{
							"10.0.0.1": "router1",
							"10.0.0.2": "router2",
						},
					},
					{
						SourceField: "in_interface",
						TargetField: "in_interface_name",
						Mapping: map[any]any{
							"router2": map[any]any{
								0: "eth0",
								1: "eth1",
								2: "eth1.1",
								3: "eth1.2",
								4: "eth2",
							},
						},
						Template: `
							{{ if eq .Msg.sampler_name "router1" }}
							igb{{ .SourceField }}
							{{ else if eq .Msg.sampler_name "router2" }}
							{{ index .Mapping.router2 .SourceField }}
							{{ end }}
						`,
					},
					{
						SourceField: "out_interface",
						TargetField: "out_interface_name",
						Mapping: map[any]any{
							"router2": map[any]any{
								0: "eth0",
								1: "eth1",
								2: "eth1.1",
								3: "eth1.2",
								4: "eth2",
							},
						},
						Template: `
							{{ if eq .Msg.sampler_name "router1" }}
							igb{{ .SourceField }}
							{{ else if eq .Msg.sampler_name "router2" }}
							{{ index .Mapping.router2 .SourceField }}
							{{ end }}
						`,
					},
				},
			},
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			e := enricher.NewFieldMapperEnricher(tc.config)
			got := e.Process(tc.input)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf("\"%s\":\n%s", name, diff)
			}
		})
	}
}
