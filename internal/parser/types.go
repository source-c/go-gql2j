package parser

import (
	"github.com/source-c/go-gql2j/internal/errors"
)

// TypeKind represents the kind of GraphQL type.
type TypeKind string

const (
	TypeKindObject      TypeKind = "OBJECT"
	TypeKindInterface   TypeKind = "INTERFACE"
	TypeKindInputObject TypeKind = "INPUT_OBJECT"
	TypeKindEnum        TypeKind = "ENUM"
	TypeKindUnion       TypeKind = "UNION"
	TypeKindScalar      TypeKind = "SCALAR"
)

// TypeDef represents a parsed GraphQL type definition.
type TypeDef struct {
	Name        string
	Kind        TypeKind
	Description string
	Fields      []*FieldDef
	EnumValues  []*EnumValueDef
	Interfaces  []string
	Directives  []*DirectiveDef
	Location    *errors.Location
}

// FieldDef represents a parsed field definition.
type FieldDef struct {
	Name         string
	Description  string
	Type         *TypeRef
	Arguments    []*ArgumentDef
	Directives   []*DirectiveDef
	DefaultValue interface{}
	Location     *errors.Location
}

// TypeRef represents a type reference (can be named, list, or non-null).
type TypeRef struct {
	Name    string   // For named types
	Elem    *TypeRef // For list types
	NonNull bool
}

// IsNamed returns true if this is a named type (not a list).
func (t *TypeRef) IsNamed() bool {
	return t.Elem == nil
}

// IsList returns true if this is a list type.
func (t *TypeRef) IsList() bool {
	return t.Elem != nil
}

// NamedType returns the innermost named type.
func (t *TypeRef) NamedType() string {
	if t.Elem != nil {
		return t.Elem.NamedType()
	}
	return t.Name
}

// EnumValueDef represents a parsed enum value.
type EnumValueDef struct {
	Name        string
	Description string
	Directives  []*DirectiveDef
	Location    *errors.Location
}

// ArgumentDef represents a parsed argument definition.
type ArgumentDef struct {
	Name         string
	Description  string
	Type         *TypeRef
	DefaultValue interface{}
	Location     *errors.Location
}

// DirectiveDef represents a parsed directive.
type DirectiveDef struct {
	Name      string
	Arguments map[string]interface{}
	Location  *errors.Location
}

// Schema represents the complete parsed schema.
type Schema struct {
	Types      map[string]*TypeDef
	Directives map[string]*DirectiveDefinition
}

// DirectiveDefinition represents the definition of a directive.
type DirectiveDefinition struct {
	Name        string
	Description string
	Arguments   []*ArgumentDef
	Locations   []string
}

// GetType returns a type by name.
func (s *Schema) GetType(name string) *TypeDef {
	if s.Types == nil {
		return nil
	}
	return s.Types[name]
}

// ObjectTypes returns all object types.
func (s *Schema) ObjectTypes() []*TypeDef {
	return s.TypesByKind(TypeKindObject)
}

// InterfaceTypes returns all interface types.
func (s *Schema) InterfaceTypes() []*TypeDef {
	return s.TypesByKind(TypeKindInterface)
}

// InputTypes returns all input object types.
func (s *Schema) InputTypes() []*TypeDef {
	return s.TypesByKind(TypeKindInputObject)
}

// EnumTypes returns all enum types.
func (s *Schema) EnumTypes() []*TypeDef {
	return s.TypesByKind(TypeKindEnum)
}

// TypesByKind returns all types of a specific kind.
func (s *Schema) TypesByKind(kind TypeKind) []*TypeDef {
	var result []*TypeDef
	for _, t := range s.Types {
		if t.Kind == kind {
			result = append(result, t)
		}
	}
	return result
}

// HasDirective checks if the type has a specific directive.
func (t *TypeDef) HasDirective(name string) bool {
	for _, d := range t.Directives {
		if d.Name == name {
			return true
		}
	}
	return false
}

// GetDirective returns the directive with the given name, or nil.
func (t *TypeDef) GetDirective(name string) *DirectiveDef {
	for _, d := range t.Directives {
		if d.Name == name {
			return d
		}
	}
	return nil
}

// HasDirective checks if the field has a specific directive.
func (f *FieldDef) HasDirective(name string) bool {
	for _, d := range f.Directives {
		if d.Name == name {
			return true
		}
	}
	return false
}

// GetDirective returns the directive with the given name, or nil.
func (f *FieldDef) GetDirective(name string) *DirectiveDef {
	for _, d := range f.Directives {
		if d.Name == name {
			return d
		}
	}
	return nil
}

// HasDirective checks if the enum value has a specific directive.
func (e *EnumValueDef) HasDirective(name string) bool {
	for _, d := range e.Directives {
		if d.Name == name {
			return true
		}
	}
	return false
}

// GetDirective returns the directive with the given name, or nil.
func (e *EnumValueDef) GetDirective(name string) *DirectiveDef {
	for _, d := range e.Directives {
		if d.Name == name {
			return d
		}
	}
	return nil
}

// GetArgumentString returns a directive argument as a string.
func (d *DirectiveDef) GetArgumentString(name string) string {
	if d.Arguments == nil {
		return ""
	}
	if v, ok := d.Arguments[name]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// GetArgumentInt returns a directive argument as an int.
func (d *DirectiveDef) GetArgumentInt(name string) (int, bool) {
	if d.Arguments == nil {
		return 0, false
	}
	if v, ok := d.Arguments[name]; ok {
		switch n := v.(type) {
		case int:
			return n, true
		case int64:
			return int(n), true
		case float64:
			return int(n), true
		}
	}
	return 0, false
}

// GetArgumentBool returns a directive argument as a bool.
func (d *DirectiveDef) GetArgumentBool(name string) (bool, bool) {
	if d.Arguments == nil {
		return false, false
	}
	if v, ok := d.Arguments[name]; ok {
		if b, ok := v.(bool); ok {
			return b, true
		}
	}
	return false, false
}

// GetArgumentStringSlice returns a directive argument as a string slice.
func (d *DirectiveDef) GetArgumentStringSlice(name string) []string {
	if d.Arguments == nil {
		return nil
	}
	if v, ok := d.Arguments[name]; ok {
		if arr, ok := v.([]interface{}); ok {
			result := make([]string, 0, len(arr))
			for _, item := range arr {
				if s, ok := item.(string); ok {
					result = append(result, s)
				}
			}
			return result
		}
	}
	return nil
}
