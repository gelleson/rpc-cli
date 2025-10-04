package parser

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

// AttributeDecoder handles decoding HCL attributes to Go types
type AttributeDecoder struct{}

// NewAttributeDecoder creates a new AttributeDecoder
func NewAttributeDecoder() *AttributeDecoder {
	return &AttributeDecoder{}
}

// DecodeString decodes an HCL attribute to a string
func (d *AttributeDecoder) DecodeString(attr *hcl.Attribute, target *string) error {
	val, diags := attr.Expr.Value(nil)
	if diags.HasErrors() {
		return fmt.Errorf("failed to decode string: %s", diags.Error())
	}

	if val.Type() != cty.String {
		return fmt.Errorf("expected string, got %s", val.Type().FriendlyName())
	}

	*target = val.AsString()
	return nil
}

// DecodeInt decodes an HCL attribute to an integer
func (d *AttributeDecoder) DecodeInt(attr *hcl.Attribute, target *int) error {
	val, diags := attr.Expr.Value(nil)
	if diags.HasErrors() {
		return fmt.Errorf("failed to decode int: %s", diags.Error())
	}

	if val.Type() != cty.Number {
		return fmt.Errorf("expected number, got %s", val.Type().FriendlyName())
	}

	bf := val.AsBigFloat()
	i, _ := bf.Int64()
	*target = int(i)
	return nil
}

// DecodeStringMap decodes an HCL attribute to a map[string]string
func (d *AttributeDecoder) DecodeStringMap(attr *hcl.Attribute, target *map[string]string) error {
	val, diags := attr.Expr.Value(nil)
	if diags.HasErrors() {
		return fmt.Errorf("failed to decode map: %s", diags.Error())
	}

	if !val.Type().IsMapType() && !val.Type().IsObjectType() {
		return fmt.Errorf("expected map or object, got %s", val.Type().FriendlyName())
	}

	*target = make(map[string]string)
	it := val.ElementIterator()
	for it.Next() {
		key, elemVal := it.Element()
		(*target)[key.AsString()] = elemVal.AsString()
	}

	return nil
}
