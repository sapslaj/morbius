package destination_test

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sapslaj/morbius/destination"
)

func TestStdoutDestination(t *testing.T) {
	type test struct {
		format string
		skip   string
		input  map[string]interface{}
		want   string
	}

	tests := []test{
		{
			format: "",
			input: map[string]interface{}{
				"str": "nice",
				"int": 69,
			},
			want: "{\"int\":69,\"str\":\"nice\"}\n",
		},
		{
			format: "json",
			input: map[string]interface{}{
				"str": "nice",
				"int": 69,
			},
			want: "{\"int\":69,\"str\":\"nice\"}\n",
		},
		{
			format: "logfmt",
			input: map[string]interface{}{
				"str": "nice",
				"int": 69,
			},
			want: "str=nice int=69\n",
		},
	}

	for _, tc := range tests {
		if tc.input == nil {
			t.Logf("\"%s\": skip (unimplmented)", tc.format)
			continue
		}
		if tc.skip != "" {
			t.Logf("\"%s\": skip (%s)", tc.format, tc.skip)
			continue
		}
		var output bytes.Buffer
		d := destination.NewStdoutDestination(&destination.StdoutDestinationConfig{
			Format: tc.format,
		})
		d.Writer = &output
		d.Publish(tc.input)
		got := output.String()
		if diff := cmp.Diff(tc.want, got); diff != "" {
			t.Logf("\"%s\":\n%s", tc.format, diff)
			t.Fail()
		}
	}
}
