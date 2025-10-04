package parser

import "github.com/zclconf/go-cty/cty"

// ConvertCtyToGo converts a cty.Value to native Go types
// This handles all HCL value types including primitives, lists, and maps
func ConvertCtyToGo(val cty.Value) interface{} {
	if val.IsNull() {
		return nil
	}

	valType := val.Type()

	// Handle primitive types
	switch {
	case valType == cty.String:
		return val.AsString()

	case valType == cty.Number:
		bf := val.AsBigFloat()
		if bf.IsInt() {
			i, _ := bf.Int64()
			return int(i)
		}
		f, _ := bf.Float64()
		return f

	case valType == cty.Bool:
		return val.True()

	case valType.IsListType() || valType.IsTupleType():
		return convertCtyList(val)

	case valType.IsMapType() || valType.IsObjectType():
		return convertCtyMap(val)

	default:
		// Fallback - return string representation
		return val.AsString()
	}
}

// convertCtyList converts a cty list or tuple to a Go slice
func convertCtyList(val cty.Value) []interface{} {
	var result []interface{}
	it := val.ElementIterator()
	for it.Next() {
		_, elemVal := it.Element()
		result = append(result, ConvertCtyToGo(elemVal))
	}
	return result
}

// convertCtyMap converts a cty map or object to a Go map
func convertCtyMap(val cty.Value) map[string]interface{} {
	result := make(map[string]interface{})
	it := val.ElementIterator()
	for it.Next() {
		key, elemVal := it.Element()
		keyStr := key.AsString()
		result[keyStr] = ConvertCtyToGo(elemVal)
	}
	return result
}
