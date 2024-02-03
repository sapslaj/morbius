package config_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sapslaj/morbius/config"
)

func TestMapGetFunc(t *testing.T) {
	t.Parallel()
	type input struct {
		m   map[string]any
		key string
		f   func(v any, present bool) string
	}
	type test struct {
		skip  string
		input input
		want  string
	}

	tests := map[string]test{
		"passed value and true for existing key": {
			input: input{
				m: map[string]any{
					"test": "hello",
				},
				key: "test",
				f: func(v any, present bool) string {
					if !present {
						t.Error("test key should be present")
					}
					return v.(string)
				},
			},
			want: "hello",
		},
		"passed nil value and false for non-existing key": {
			input: input{
				m: map[string]any{
					"test": "hello",
				},
				key: "nope",
				f: func(v any, present bool) string {
					if present {
						t.Error("nope key should not be present")
					}
					if v != nil {
						t.Errorf("value is not nil, got: %v", v)
					}
					return ""
				},
			},
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if tc.input.m == nil {
				t.Logf("\"%s\": skip (unimplmented)", name)
				return
			}
			if tc.skip != "" {
				t.Logf("\"%s\": skip (%s)", name, tc.skip)
				return
			}
			got := config.MapGetFunc(tc.input.m, tc.input.key, tc.input.f)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Logf("\"%s\":\n%s", name, diff)
				t.Fail()
			}
		})
	}
}

func TestMapGetDefault(t *testing.T) {
	t.Parallel()
	type input struct {
		m   map[string]any
		key string
		def string
	}
	type test struct {
		skip  string
		input input
		want  string
	}

	tests := map[string]test{
		"returns expected value when key is present": {
			input: input{
				m: map[string]any{
					"test": "hello",
				},
				key: "test",
				def: "oops",
			},
			want: "hello",
		},
		"returns default value if key is not present": {
			input: input{
				m: map[string]any{
					"test": "hello",
				},
				key: "nope",
				def: "oops",
			},
			want: "oops",
		},
		"returns default value if value is wrong type": {
			input: input{
				m: map[string]any{
					"test": 69,
				},
				key: "test",
				def: "oops",
			},
			want: "oops",
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if tc.input.m == nil {
				t.Logf("\"%s\": skip (unimplmented)", name)
				return
			}
			if tc.skip != "" {
				t.Logf("\"%s\": skip (%s)", name, tc.skip)
				return
			}
			got := config.MapGetDefault(tc.input.m, tc.input.key, tc.input.def)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Logf("\"%s\":\n%s", name, diff)
				t.Fail()
			}
		})
	}
}
