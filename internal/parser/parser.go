package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/source-c/go-gql2j/internal/errors"
)

// Parser handles GraphQL schema parsing.
type Parser struct{}

// NewParser creates a new parser instance.
func NewParser() *Parser {
	return &Parser{}
}

// ParseFile parses a GraphQL schema from a file.
func (p *Parser) ParseFile(path string) (*Schema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.NewParseError("failed to read schema file", err).
			WithLocation(&errors.Location{File: path})
	}

	return p.Parse(string(data), path)
}

// ParseFiles parses multiple GraphQL schema files.
func (p *Parser) ParseFiles(paths []string) (*Schema, error) {
	var sources []*ast.Source
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, errors.NewParseError("failed to read schema file", err).
				WithLocation(&errors.Location{File: path})
		}
		sources = append(sources, &ast.Source{
			Name:  path,
			Input: string(data),
		})
	}

	return p.parseFromSources(sources)
}

// ParseWithIncludes parses a main schema file and additional include patterns.
func (p *Parser) ParseWithIncludes(mainPath string, includePatterns []string) (*Schema, error) {
	paths := []string{mainPath}

	for _, pattern := range includePatterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return nil, errors.NewParseError(
				fmt.Sprintf("invalid include pattern: %s", pattern),
				err,
			)
		}
		paths = append(paths, matches...)
	}

	// Remove duplicates
	seen := make(map[string]bool)
	var uniquePaths []string
	for _, path := range paths {
		absPath, _ := filepath.Abs(path)
		if !seen[absPath] {
			seen[absPath] = true
			uniquePaths = append(uniquePaths, path)
		}
	}

	return p.ParseFiles(uniquePaths)
}

// Parse parses a GraphQL schema from a string.
func (p *Parser) Parse(input string, sourceName string) (*Schema, error) {
	source := &ast.Source{
		Name:  sourceName,
		Input: input,
	}

	return p.parseFromSources([]*ast.Source{source})
}

func (p *Parser) parseFromSources(sources []*ast.Source) (*Schema, error) {
	astSchema, err := gqlparser.LoadSchema(sources...)
	if err != nil {
		// Try with schema definition if not present
		if !hasSchemaDefinition(sources) {
			augmentedSources := make([]*ast.Source, len(sources))
			copy(augmentedSources, sources)
			augmentedSources[0] = &ast.Source{
				Name:  sources[0].Name,
				Input: "schema { query: Query }\n" + sources[0].Input,
			}
			astSchema, err = gqlparser.LoadSchema(augmentedSources...)
		}
		if err != nil {
			return nil, errors.NewParseError("failed to parse GraphQL schema", err)
		}
	}

	return p.convertSchema(astSchema)
}

func hasSchemaDefinition(sources []*ast.Source) bool {
	for _, s := range sources {
		if strings.Contains(s.Input, "schema {") || strings.Contains(s.Input, "schema{") {
			return true
		}
	}
	return false
}

func (p *Parser) convertSchema(astSchema *ast.Schema) (*Schema, error) {
	schema := &Schema{
		Types:      make(map[string]*TypeDef),
		Directives: make(map[string]*DirectiveDefinition),
	}

	// Convert types
	for name, def := range astSchema.Types {
		if isBuiltinType(name) {
			continue
		}

		typeDef, err := p.convertTypeDef(def)
		if err != nil {
			return nil, err
		}
		if typeDef != nil {
			schema.Types[name] = typeDef
		}
	}

	// Convert directive definitions
	for name, def := range astSchema.Directives {
		schema.Directives[name] = p.convertDirectiveDefinition(def)
	}

	return schema, nil
}

func (p *Parser) convertTypeDef(def *ast.Definition) (*TypeDef, error) {
	kind, ok := convertTypeKind(def.Kind)
	if !ok {
		return nil, nil // Skip unsupported types
	}

	sourceName := ""
	if def.Position != nil {
		sourceName = def.Position.Src.Name
	}

	typeDef := &TypeDef{
		Name:        def.Name,
		Kind:        kind,
		Description: def.Description,
		Directives:  extractDirectives(def.Directives, sourceName),
	}

	if def.Position != nil {
		typeDef.Location = &errors.Location{
			File:   sourceName,
			Line:   def.Position.Line,
			Column: def.Position.Column,
		}
	}

	// Convert interfaces
	for _, iface := range def.Interfaces {
		typeDef.Interfaces = append(typeDef.Interfaces, iface)
	}

	// Convert fields
	for _, field := range def.Fields {
		fieldDef := p.convertFieldDef(field, sourceName)
		typeDef.Fields = append(typeDef.Fields, fieldDef)
	}

	// Convert enum values
	for _, ev := range def.EnumValues {
		enumValue := p.convertEnumValue(ev, sourceName)
		typeDef.EnumValues = append(typeDef.EnumValues, enumValue)
	}

	return typeDef, nil
}

func (p *Parser) convertFieldDef(field *ast.FieldDefinition, sourceName string) *FieldDef {
	fieldDef := &FieldDef{
		Name:        field.Name,
		Description: field.Description,
		Type:        convertTypeRef(field.Type),
		Directives:  extractDirectives(field.Directives, sourceName),
	}

	if field.Position != nil {
		fieldDef.Location = &errors.Location{
			File:   sourceName,
			Line:   field.Position.Line,
			Column: field.Position.Column,
		}
	}

	// Convert arguments
	for _, arg := range field.Arguments {
		argDef := &ArgumentDef{
			Name:        arg.Name,
			Description: arg.Description,
			Type:        convertTypeRef(arg.Type),
		}
		if arg.DefaultValue != nil {
			argDef.DefaultValue = valueToInterface(arg.DefaultValue)
		}
		fieldDef.Arguments = append(fieldDef.Arguments, argDef)
	}

	if field.DefaultValue != nil {
		fieldDef.DefaultValue = valueToInterface(field.DefaultValue)
	}

	return fieldDef
}

func (p *Parser) convertEnumValue(ev *ast.EnumValueDefinition, sourceName string) *EnumValueDef {
	enumValue := &EnumValueDef{
		Name:        ev.Name,
		Description: ev.Description,
		Directives:  extractDirectives(ev.Directives, sourceName),
	}

	if ev.Position != nil {
		enumValue.Location = &errors.Location{
			File:   sourceName,
			Line:   ev.Position.Line,
			Column: ev.Position.Column,
		}
	}

	return enumValue
}

func (p *Parser) convertDirectiveDefinition(def *ast.DirectiveDefinition) *DirectiveDefinition {
	directive := &DirectiveDefinition{
		Name:        def.Name,
		Description: def.Description,
	}

	for _, arg := range def.Arguments {
		argDef := &ArgumentDef{
			Name:        arg.Name,
			Description: arg.Description,
			Type:        convertTypeRef(arg.Type),
		}
		if arg.DefaultValue != nil {
			argDef.DefaultValue = valueToInterface(arg.DefaultValue)
		}
		directive.Arguments = append(directive.Arguments, argDef)
	}

	for _, loc := range def.Locations {
		directive.Locations = append(directive.Locations, string(loc))
	}

	return directive
}

func convertTypeRef(astType *ast.Type) *TypeRef {
	if astType == nil {
		return nil
	}

	if astType.Elem != nil {
		// List type
		return &TypeRef{
			Elem:    convertTypeRef(astType.Elem),
			NonNull: astType.NonNull,
		}
	}

	// Named type
	return &TypeRef{
		Name:    astType.NamedType,
		NonNull: astType.NonNull,
	}
}

func convertTypeKind(kind ast.DefinitionKind) (TypeKind, bool) {
	switch kind {
	case ast.Object:
		return TypeKindObject, true
	case ast.Interface:
		return TypeKindInterface, true
	case ast.InputObject:
		return TypeKindInputObject, true
	case ast.Enum:
		return TypeKindEnum, true
	case ast.Union:
		return TypeKindUnion, true
	case ast.Scalar:
		return TypeKindScalar, true
	default:
		return "", false
	}
}

func isBuiltinType(name string) bool {
	builtinTypes := map[string]bool{
		"String":              true,
		"Int":                 true,
		"Float":               true,
		"Boolean":             true,
		"ID":                  true,
		"__Schema":            true,
		"__Type":              true,
		"__Field":             true,
		"__Directive":         true,
		"__EnumValue":         true,
		"__InputValue":        true,
		"__TypeKind":          true,
		"__DirectiveLocation": true,
	}
	return builtinTypes[name] || strings.HasPrefix(name, "__")
}
