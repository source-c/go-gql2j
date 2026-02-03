package generator

import (
	"strings"
	"unicode"

	"github.com/source-c/go-gql2j/internal/config"
	"github.com/source-c/go-gql2j/internal/parser"
)

// NamingHelper handles naming conventions.
type NamingHelper struct {
	config *config.NamingConfig
}

// NewNamingHelper creates a new naming helper.
func NewNamingHelper(cfg *config.NamingConfig) *NamingHelper {
	return &NamingHelper{config: cfg}
}

// GetTypeName returns the Java type name for a type definition.
func (n *NamingHelper) GetTypeName(typeDef *parser.TypeDef) string {
	// Check for @javaName directive
	if javaName := parser.ExtractJavaNameDirective(typeDef.Directives); javaName != nil {
		return javaName.Name
	}

	name := typeDef.Name

	// Apply naming conventions from config
	switch typeDef.Kind {
	case parser.TypeKindInterface:
		if n.config.InterfacePrefix != "" {
			name = n.config.InterfacePrefix + name
		}
	case parser.TypeKindObject, parser.TypeKindInputObject:
		if n.config.ClassSuffix != "" {
			name = name + n.config.ClassSuffix
		}
	}

	return name
}

// GetFieldName returns the Java field name for a field definition.
func (n *NamingHelper) GetFieldName(field *parser.FieldDef) string {
	// Check for @javaName directive
	if javaName := parser.ExtractJavaNameDirective(field.Directives); javaName != nil {
		return javaName.Name
	}

	name := field.Name

	// Apply naming convention
	switch n.config.FieldCase {
	case config.FieldCaseSnake:
		name = toSnakeCase(name)
	case config.FieldCaseCamel:
		fallthrough
	default:
		name = toCamelCase(name)
	}

	return name
}

// GetEnumValueName returns the Java enum value name.
func (n *NamingHelper) GetEnumValueName(enumValue *parser.EnumValueDef) string {
	// Check for @javaName directive
	if javaName := parser.ExtractJavaNameDirective(enumValue.Directives); javaName != nil {
		return javaName.Name
	}

	// Enum values are typically uppercase in Java
	return enumValue.Name
}

// GetGetterName returns the getter method name for a field.
func (n *NamingHelper) GetGetterName(fieldName string, isBoolean bool) string {
	prefix := "get"
	if isBoolean {
		// Use "is" prefix for boolean fields if not already starting with "is"
		if !strings.HasPrefix(strings.ToLower(fieldName), "is") {
			prefix = "is"
		}
	}
	return prefix + capitalizeFirst(fieldName)
}

// GetSetterName returns the setter method name for a field.
func (n *NamingHelper) GetSetterName(fieldName string) string {
	return "set" + capitalizeFirst(fieldName)
}

// capitalizeFirst capitalizes the first letter of a string.
func capitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

// toCamelCase converts a string to camelCase.
func toCamelCase(s string) string {
	if s == "" {
		return s
	}

	// Handle snake_case input
	if strings.Contains(s, "_") {
		parts := strings.Split(s, "_")
		result := strings.ToLower(parts[0])
		for i := 1; i < len(parts); i++ {
			if parts[i] != "" {
				result += capitalizeFirst(strings.ToLower(parts[i]))
			}
		}
		return result
	}

	// Already camelCase or PascalCase, ensure first char is lowercase
	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}

// toSnakeCase converts a string to snake_case.
func toSnakeCase(s string) string {
	if s == "" {
		return s
	}

	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// ToPascalCase converts a string to PascalCase.
func ToPascalCase(s string) string {
	if s == "" {
		return s
	}

	// Handle snake_case input
	if strings.Contains(s, "_") {
		parts := strings.Split(s, "_")
		var result string
		for _, part := range parts {
			if part != "" {
				result += capitalizeFirst(strings.ToLower(part))
			}
		}
		return result
	}

	// Already camelCase or PascalCase, ensure first char is uppercase
	return capitalizeFirst(s)
}

// IsJavaKeyword checks if a word is a Java reserved keyword.
func IsJavaKeyword(word string) bool {
	keywords := map[string]bool{
		"abstract": true, "assert": true, "boolean": true, "break": true,
		"byte": true, "case": true, "catch": true, "char": true,
		"class": true, "const": true, "continue": true, "default": true,
		"do": true, "double": true, "else": true, "enum": true,
		"extends": true, "final": true, "finally": true, "float": true,
		"for": true, "goto": true, "if": true, "implements": true,
		"import": true, "instanceof": true, "int": true, "interface": true,
		"long": true, "native": true, "new": true, "package": true,
		"private": true, "protected": true, "public": true, "return": true,
		"short": true, "static": true, "strictfp": true, "super": true,
		"switch": true, "synchronized": true, "this": true, "throw": true,
		"throws": true, "transient": true, "try": true, "void": true,
		"volatile": true, "while": true, "true": true, "false": true,
		"null": true, "var": true, "yield": true, "record": true,
		"sealed": true, "permits": true, "non-sealed": true,
	}
	return keywords[word]
}

// EscapeJavaKeyword escapes a Java keyword by adding an underscore prefix.
func EscapeJavaKeyword(name string) string {
	if IsJavaKeyword(name) {
		return "_" + name
	}
	return name
}
