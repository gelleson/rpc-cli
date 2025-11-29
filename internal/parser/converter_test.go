package parser

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func TestConvertCtyToGo(t *testing.T) {
	tests := []struct {
		name string
		val  cty.Value
		want any
	}{
		{
			name: "null value",
			val:  cty.NullVal(cty.String),
			want: nil,
		},
		{
			name: "string value",
			val:  cty.StringVal("hello"),
			want: "hello",
		},
		{
			name: "integer value",
			val:  cty.NumberIntVal(42),
			want: 42,
		},
		{
			name: "float value",
			val:  cty.NumberFloatVal(3.14),
			want: 3.14,
		},
		{
			name: "bool value true",
			val:  cty.BoolVal(true),
			want: true,
		},
		{
			name: "bool value false",
			val:  cty.BoolVal(false),
			want: false,
		},
		{
			name: "list of strings",
			val: cty.ListVal([]cty.Value{
				cty.StringVal("a"),
				cty.StringVal("b"),
				cty.StringVal("c"),
			}),
			want: []any{"a", "b", "c"},
		},
		{
			name: "tuple with mixed types",
			val: cty.TupleVal([]cty.Value{
				cty.StringVal("hello"),
				cty.NumberIntVal(42),
				cty.BoolVal(true),
			}),
			want: []any{"hello", 42, true},
		},
		{
			name: "map of strings",
			val: cty.MapVal(map[string]cty.Value{
				"key1": cty.StringVal("value1"),
				"key2": cty.StringVal("value2"),
			}),
			want: map[string]any{
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			name: "object with mixed types",
			val: cty.ObjectVal(map[string]cty.Value{
				"name":   cty.StringVal("John"),
				"age":    cty.NumberIntVal(30),
				"active": cty.BoolVal(true),
			}),
			want: map[string]any{
				"name":   "John",
				"age":    30,
				"active": true,
			},
		},
		{
			name: "nested structures",
			val: cty.ObjectVal(map[string]cty.Value{
				"user": cty.ObjectVal(map[string]cty.Value{
					"name": cty.StringVal("Alice"),
					"tags": cty.ListVal([]cty.Value{
						cty.StringVal("admin"),
						cty.StringVal("user"),
					}),
				}),
			}),
			want: map[string]any{
				"user": map[string]any{
					"name": "Alice",
					"tags": []any{"admin", "user"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertCtyToGo(tt.val)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertCtyToGo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertCtyList(t *testing.T) {
	val := cty.ListVal([]cty.Value{
		cty.NumberIntVal(1),
		cty.NumberIntVal(2),
		cty.NumberIntVal(3),
	})

	got := convertCtyList(val)
	want := []any{1, 2, 3}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("convertCtyList() = %v, want %v", got, want)
	}
}

func TestConvertCtyMap(t *testing.T) {
	val := cty.MapVal(map[string]cty.Value{
		"a": cty.NumberIntVal(1),
		"b": cty.NumberIntVal(2),
	})

	got := convertCtyMap(val)
	want := map[string]any{
		"a": 1,
		"b": 2,
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("convertCtyMap() = %v, want %v", got, want)
	}
}

func TestConvertCtyToGo_LargeNumber(t *testing.T) {
	// Test with a large number that doesn't fit in int64
	bigNum := new(big.Float).SetFloat64(1e20)
	val := cty.NumberVal(bigNum)

	got := ConvertCtyToGo(val)

	// For very large numbers, it depends on whether they're represented as int or float
	// The implementation converts based on BigFloat.IsInt()
	switch got.(type) {
	case int, float64:
		// Both are acceptable for large numbers
	default:
		t.Errorf("Expected int or float64 for large number, got %T", got)
	}
}
