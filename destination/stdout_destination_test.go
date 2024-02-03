package destination_test

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sapslaj/morbius/destination"
)

func TestStdoutDestination(t *testing.T) {
	t.Parallel()
	type test struct {
		format string
		skip   string
		input  map[string]interface{}
		want   string
	}

	tests := map[string]test{
		"default format": {
			format: "",
			input: map[string]interface{}{
				"str": "nice",
				"int": 69,
			},
			want: "{\"int\":69,\"str\":\"nice\"}\n",
		},
		"json format": {
			format: "json",
			input: map[string]interface{}{
				"str": "nice",
				"int": 69,
			},
			want: "{\"int\":69,\"str\":\"nice\"}\n",
		},
		"logfmt format": {
			format: "logfmt",
			input: map[string]interface{}{
				"str": "nice",
				"int": 69,
			},
			want: "str=nice int=69\n",
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
			var output bytes.Buffer
			d := destination.NewStdoutDestination(&destination.StdoutDestinationConfig{
				Format: tc.format,
			})
			d.Writer = &output
			d.Publish(tc.input)
			got := output.String()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Logf("\"%s\":\n%s", name, diff)
				t.Fail()
			}
		})
	}
}
