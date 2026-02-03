package parser

import (
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/source-c/go-gql2j/internal/errors"
)

// Supported directive names.
const (
	DirectiveSkip       = "skip"
	DirectiveJavaName   = "javaName"
	DirectiveJavaType   = "javaType"
	DirectiveDeprecated = "deprecated"
	DirectiveAnnotation = "annotation"
	DirectiveConstraint = "constraint"
	DirectiveLombok     = "lombok"
	DirectiveCollection = "collection"
)

// extractDirectives converts AST directives to our DirectiveDef format.
func extractDirectives(astDirectives ast.DirectiveList, source string) []*DirectiveDef {
	if len(astDirectives) == 0 {
		return nil
	}

	directives := make([]*DirectiveDef, 0, len(astDirectives))
	for _, d := range astDirectives {
		directive := &DirectiveDef{
			Name:      d.Name,
			Arguments: extractDirectiveArguments(d.Arguments),
		}
		if d.Position != nil {
			directive.Location = &errors.Location{
				File:   source,
				Line:   d.Position.Line,
				Column: d.Position.Column,
			}
		}
		directives = append(directives, directive)
	}

	return directives
}

// extractDirectiveArguments converts AST arguments to a map.
func extractDirectiveArguments(args ast.ArgumentList) map[string]interface{} {
	if len(args) == 0 {
		return nil
	}

	result := make(map[string]interface{})
	for _, arg := range args {
		result[arg.Name] = valueToInterface(arg.Value)
	}
	return result
}

// valueToInterface converts an AST value to a Go interface.
func valueToInterface(v *ast.Value) interface{} {
	if v == nil {
		return nil
	}

	switch v.Kind {
	case ast.StringValue, ast.BlockValue:
		return v.Raw
	case ast.IntValue:
		// Parse as int64 to handle large numbers
		var n int64
		for _, c := range v.Raw {
			n = n*10 + int64(c-'0')
		}
		return n
	case ast.FloatValue:
		// Simple float parsing
		var f float64
		var decimal float64 = 0.1
		var pastDecimal bool
		for _, c := range v.Raw {
			if c == '.' {
				pastDecimal = true
				continue
			}
			if pastDecimal {
				f += decimal * float64(c-'0')
				decimal *= 0.1
			} else {
				f = f*10 + float64(c-'0')
			}
		}
		return f
	case ast.BooleanValue:
		return v.Raw == "true"
	case ast.EnumValue:
		return v.Raw
	case ast.ListValue:
		list := make([]interface{}, 0, len(v.Children))
		for _, child := range v.Children {
			list = append(list, valueToInterface(child.Value))
		}
		return list
	case ast.ObjectValue:
		obj := make(map[string]interface{})
		for _, child := range v.Children {
			obj[child.Name] = valueToInterface(child.Value)
		}
		return obj
	case ast.NullValue:
		return nil
	case ast.Variable:
		// Variables shouldn't appear in schema definitions
		return nil
	default:
		return v.Raw
	}
}

// SkipDirectiveInfo extracts information from @skip directive.
type SkipDirectiveInfo struct {
	Applied bool
}

// ExtractSkipDirective extracts @skip directive info.
func ExtractSkipDirective(directives []*DirectiveDef) *SkipDirectiveInfo {
	for _, d := range directives {
		if d.Name == DirectiveSkip {
			return &SkipDirectiveInfo{Applied: true}
		}
	}
	return nil
}

// JavaNameDirectiveInfo extracts information from @javaName directive.
type JavaNameDirectiveInfo struct {
	Name string
}

// ExtractJavaNameDirective extracts @javaName directive info.
func ExtractJavaNameDirective(directives []*DirectiveDef) *JavaNameDirectiveInfo {
	for _, d := range directives {
		if d.Name == DirectiveJavaName {
			name := d.GetArgumentString("name")
			if name != "" {
				return &JavaNameDirectiveInfo{Name: name}
			}
		}
	}
	return nil
}

// JavaTypeDirectiveInfo extracts information from @javaType directive.
type JavaTypeDirectiveInfo struct {
	Type    string
	Imports []string
}

// ExtractJavaTypeDirective extracts @javaType directive info.
func ExtractJavaTypeDirective(directives []*DirectiveDef) *JavaTypeDirectiveInfo {
	for _, d := range directives {
		if d.Name == DirectiveJavaType {
			typeStr := d.GetArgumentString("type")
			if typeStr != "" {
				return &JavaTypeDirectiveInfo{
					Type:    typeStr,
					Imports: d.GetArgumentStringSlice("imports"),
				}
			}
		}
	}
	return nil
}

// DeprecatedDirectiveInfo extracts information from @deprecated directive.
type DeprecatedDirectiveInfo struct {
	Reason string
}

// ExtractDeprecatedDirective extracts @deprecated directive info.
func ExtractDeprecatedDirective(directives []*DirectiveDef) *DeprecatedDirectiveInfo {
	for _, d := range directives {
		if d.Name == DirectiveDeprecated {
			return &DeprecatedDirectiveInfo{
				Reason: d.GetArgumentString("reason"),
			}
		}
	}
	return nil
}

// AnnotationDirectiveInfo extracts information from @annotation directive.
type AnnotationDirectiveInfo struct {
	Value   string
	Imports []string
}

// ExtractAnnotationDirectives extracts all @annotation directive info.
func ExtractAnnotationDirectives(directives []*DirectiveDef) []*AnnotationDirectiveInfo {
	var result []*AnnotationDirectiveInfo
	for _, d := range directives {
		if d.Name == DirectiveAnnotation {
			value := d.GetArgumentString("value")
			if value != "" {
				result = append(result, &AnnotationDirectiveInfo{
					Value:   value,
					Imports: d.GetArgumentStringSlice("imports"),
				})
			}
		}
	}
	return result
}

// ConstraintDirectiveInfo extracts information from @constraint directive.
type ConstraintDirectiveInfo struct {
	MinLength *int
	MaxLength *int
	Min       *int
	Max       *int
	Pattern   string
	NotNull   bool
	NotBlank  bool
	Email     bool
}

// ExtractConstraintDirective extracts @constraint directive info.
func ExtractConstraintDirective(directives []*DirectiveDef) *ConstraintDirectiveInfo {
	for _, d := range directives {
		if d.Name == DirectiveConstraint {
			info := &ConstraintDirectiveInfo{}
			if v, ok := d.GetArgumentInt("minLength"); ok {
				info.MinLength = &v
			}
			if v, ok := d.GetArgumentInt("maxLength"); ok {
				info.MaxLength = &v
			}
			if v, ok := d.GetArgumentInt("min"); ok {
				info.Min = &v
			}
			if v, ok := d.GetArgumentInt("max"); ok {
				info.Max = &v
			}
			info.Pattern = d.GetArgumentString("pattern")
			if v, ok := d.GetArgumentBool("notNull"); ok {
				info.NotNull = v
			}
			if v, ok := d.GetArgumentBool("notBlank"); ok {
				info.NotBlank = v
			}
			if v, ok := d.GetArgumentBool("email"); ok {
				info.Email = v
			}
			return info
		}
	}
	return nil
}

// LombokDirectiveInfo extracts information from @lombok directive.
type LombokDirectiveInfo struct {
	Exclude []string
	Include []string
}

// ExtractLombokDirective extracts @lombok directive info.
func ExtractLombokDirective(directives []*DirectiveDef) *LombokDirectiveInfo {
	for _, d := range directives {
		if d.Name == DirectiveLombok {
			return &LombokDirectiveInfo{
				Exclude: d.GetArgumentStringSlice("exclude"),
				Include: d.GetArgumentStringSlice("include"),
			}
		}
	}
	return nil
}

// CollectionDirectiveInfo extracts information from @collection directive.
type CollectionDirectiveInfo struct {
	Type string
}

// ExtractCollectionDirective extracts @collection directive info.
func ExtractCollectionDirective(directives []*DirectiveDef) *CollectionDirectiveInfo {
	for _, d := range directives {
		if d.Name == DirectiveCollection {
			typeStr := d.GetArgumentString("type")
			if typeStr != "" {
				return &CollectionDirectiveInfo{Type: typeStr}
			}
		}
	}
	return nil
}
