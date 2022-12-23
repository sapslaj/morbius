package config_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sapslaj/morbius/config"
)

func TestMapGetFunc(t *testing.T) {
	type input struct {
		m   map[string]any
		key string
		f   func(v any, present bool) string
	}
	type test struct {
		desc  string
		skip  string
		input input
		want  string
	}

	tests := []test{
		{
			desc: "Passed value and true for existing key",
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
		{
			desc: "Passed nil value and false for non-existing key",
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

	for _, tc := range tests {
		if tc.input.m == nil {
			t.Logf("\"%s\": skip (unimplmented)", tc.desc)
			continue
		}
		if tc.skip != "" {
			t.Logf("\"%s\": skip (%s)", tc.desc, tc.skip)
			continue
		}
		got := config.MapGetFunc(tc.input.m, tc.input.key, tc.input.f)
		if diff := cmp.Diff(tc.want, got); diff != "" {
			t.Logf("\"%s\":\n%s", tc.desc, diff)
			t.Fail()
		}
	}
}

func TestMapGetDefault(t *testing.T) {
	type input struct {
		m   map[string]any
		key string
		def string
	}
	type test struct {
		desc  string
		skip  string
		input input
		want  string
	}

	tests := []test{
		{
			desc: "Returns expected value when key is present",
			input: input{
				m: map[string]any{
					"test": "hello",
				},
				key: "test",
				def: "oops",
			},
			want: "hello",
		},
		{
			desc: "Returns default value if key is not present",
			input: input{
				m: map[string]any{
					"test": "hello",
				},
				key: "nope",
				def: "oops",
			},
			want: "oops",
		},
		{
			desc: "Returns default value if value is wrong type",
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

	for _, tc := range tests {
		if tc.input.m == nil {
			t.Logf("\"%s\": skip (unimplmented)", tc.desc)
			continue
		}
		if tc.skip != "" {
			t.Logf("\"%s\": skip (%s)", tc.desc, tc.skip)
			continue
		}
		got := config.MapGetDefault(tc.input.m, tc.input.key, tc.input.def)
		if diff := cmp.Diff(tc.want, got); diff != "" {
			t.Logf("\"%s\":\n%s", tc.desc, diff)
			t.Fail()
		}
	}
}
