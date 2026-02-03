package typemap

import (
	"github.com/source-c/go-gql2j/internal/config"
	"github.com/source-c/go-gql2j/internal/errors"
	"github.com/source-c/go-gql2j/internal/parser"
)

// TypeMapper handles GraphQL to Java type mapping.
type TypeMapper struct {
	config         *config.Config
	customScalars  map[string]ScalarInfo
	builtinScalars map[string]ScalarInfo
	schemaTypes    map[string]*parser.TypeDef
}

// NewTypeMapper creates a new TypeMapper with the given configuration.
func NewTypeMapper(cfg *config.Config) *TypeMapper {
	tm := &TypeMapper{
		config:         cfg,
		customScalars:  make(map[string]ScalarInfo),
		builtinScalars: BuiltinScalars(),
		schemaTypes:    make(map[string]*parser.TypeDef),
	}

	// Load custom scalar mappings from config
	for name, mapping := range cfg.TypeMappings.Scalars {
		tm.customScalars[name] = ScalarInfo{
			JavaType: mapping.JavaType,
			Imports:  mapping.Imports,
		}
	}

	return tm
}

// SetSchemaTypes sets the schema types for reference type resolution.
func (tm *TypeMapper) SetSchemaTypes(types map[string]*parser.TypeDef) {
	tm.schemaTypes = types
}

// MapResult contains the result of a type mapping.
type MapResult struct {
	JavaType     string
	Imports      []string
	IsPrimitive  bool
	IsCollection bool
	IsOptional   bool
	ElementType  string // For collections/optionals
}

// MapType maps a GraphQL type reference to Java.
func (tm *TypeMapper) MapType(typeRef *parser.TypeRef) (*MapResult, error) {
	if typeRef == nil {
		return &MapResult{JavaType: "Object"}, nil
	}

	return tm.mapTypeInternal(typeRef, true)
}

func (tm *TypeMapper) mapTypeInternal(typeRef *parser.TypeRef, topLevel bool) (*MapResult, error) {
	// Handle list types
	if typeRef.IsList() {
		elemResult, err := tm.mapTypeInternal(typeRef.Elem, false)
		if err != nil {
			return nil, errors.NewTypeMappingError(
				"failed to map list element type",
				err,
			).WithSourceType(typeRef.Elem.Name)
		}

		// Box primitive element types
		elemType := elemResult.JavaType
		if elemResult.IsPrimitive {
			elemType = BoxType(elemType)
		}

		collectionType := tm.config.Java.CollectionType
		collectionInfo := GetCollectionInfo(collectionType)

		result := &MapResult{
			JavaType:     FormatCollectionType(collectionType, elemType),
			Imports:      collectionInfo.Imports,
			IsCollection: true,
			ElementType:  elemType,
		}

		// Add element type imports
		result.Imports = append(result.Imports, elemResult.Imports...)

		return result, nil
	}

	// Map the named type
	namedType := typeRef.Name
	result, err := tm.mapNamedType(namedType)
	if err != nil {
		return nil, errors.NewTypeMappingError(
			"failed to map named type",
			err,
		).WithSourceType(namedType)
	}

	// Handle nullability
	if topLevel && !typeRef.NonNull {
		result = tm.applyNullability(result)
	}

	// Use primitive for non-null built-in scalars at top level
	if topLevel && typeRef.NonNull {
		if scalar, ok := tm.builtinScalars[namedType]; ok {
			if scalar.PrimitiveType != "" {
				result.JavaType = scalar.PrimitiveType
				result.IsPrimitive = true
			}
		}
	}

	return result, nil
}

func (tm *TypeMapper) mapNamedType(name string) (*MapResult, error) {
	// Check custom scalars first (from config)
	if scalar, ok := tm.customScalars[name]; ok {
		return &MapResult{
			JavaType:    scalar.JavaType,
			Imports:     scalar.Imports,
			IsPrimitive: scalar.PrimitiveType != "",
		}, nil
	}

	// Check built-in scalars
	if scalar, ok := tm.builtinScalars[name]; ok {
		return &MapResult{
			JavaType:    scalar.JavaType,
			Imports:     scalar.Imports,
			IsPrimitive: false, // Use wrapper by default
		}, nil
	}

	// Check common scalars
	commonScalars := CommonScalars()
	if scalar, ok := commonScalars[name]; ok {
		return &MapResult{
			JavaType:    scalar.JavaType,
			Imports:     scalar.Imports,
			IsPrimitive: false,
		}, nil
	}

	// Check if it's a schema type (object, interface, enum, input)
	if typeDef, ok := tm.schemaTypes[name]; ok {
		return &MapResult{
			JavaType: tm.getJavaTypeName(typeDef),
		}, nil
	}

	// Unknown type - use as-is (assume it's a custom type in the schema)
	return &MapResult{
		JavaType: name,
	}, nil
}

func (tm *TypeMapper) getJavaTypeName(typeDef *parser.TypeDef) string {
	// Check for @javaName directive
	if javaName := parser.ExtractJavaNameDirective(typeDef.Directives); javaName != nil {
		return javaName.Name
	}

	name := typeDef.Name

	// Apply naming conventions from config
	switch typeDef.Kind {
	case parser.TypeKindInterface:
		if tm.config.Java.Naming.InterfacePrefix != "" {
			name = tm.config.Java.Naming.InterfacePrefix + name
		}
	case parser.TypeKindObject, parser.TypeKindInputObject:
		if tm.config.Java.Naming.ClassSuffix != "" {
			name = name + tm.config.Java.Naming.ClassSuffix
		}
	}

	return name
}

func (tm *TypeMapper) applyNullability(result *MapResult) *MapResult {
	switch tm.config.Java.NullableHandling {
	case config.NullableOptional:
		// Box primitives for Optional
		javaType := result.JavaType
		if result.IsPrimitive {
			javaType = BoxType(javaType)
		}
		result.JavaType = FormatOptionalType(javaType)
		result.Imports = append(result.Imports, GetOptionalInfo().Imports...)
		result.IsOptional = true
		result.IsPrimitive = false

	case config.NullableAnnotation:
		// Keep the type as-is but use annotation (handled elsewhere)
		// Box primitives since they can't be null
		if result.IsPrimitive {
			result.JavaType = BoxType(result.JavaType)
			result.IsPrimitive = false
		}

	case config.NullableWrapper:
		fallthrough
	default:
		// Use wrapper types (default)
		if result.IsPrimitive {
			result.JavaType = BoxType(result.JavaType)
			result.IsPrimitive = false
		}
	}

	return result
}

// MapFieldType maps a field type, considering field-level directives.
func (tm *TypeMapper) MapFieldType(field *parser.FieldDef) (*MapResult, error) {
	// Check for @javaType directive
	if javaType := parser.ExtractJavaTypeDirective(field.Directives); javaType != nil {
		return &MapResult{
			JavaType: javaType.Type,
			Imports:  javaType.Imports,
		}, nil
	}

	// Check for @collection directive override
	collectionOverride := parser.ExtractCollectionDirective(field.Directives)

	result, err := tm.MapType(field.Type)
	if err != nil {
		return nil, err
	}

	// Apply collection override
	if collectionOverride != nil && result.IsCollection {
		collectionInfo := GetCollectionInfo(collectionOverride.Type)
		result.JavaType = FormatCollectionType(collectionOverride.Type, result.ElementType)
		result.Imports = collectionInfo.Imports
		result.Imports = append(result.Imports, tm.getElementImports(result.ElementType)...)
	}

	return result, nil
}

func (tm *TypeMapper) getElementImports(elementType string) []string {
	// Check if element type needs imports
	if scalar, ok := tm.customScalars[elementType]; ok {
		return scalar.Imports
	}
	commonScalars := CommonScalars()
	if scalar, ok := commonScalars[elementType]; ok {
		return scalar.Imports
	}
	return nil
}

// ValidateMapping validates that a type can be mapped.
func (tm *TypeMapper) ValidateMapping(typeRef *parser.TypeRef) error {
	if typeRef == nil {
		return nil
	}

	if typeRef.IsList() {
		return tm.ValidateMapping(typeRef.Elem)
	}

	name := typeRef.Name

	// Check known types
	if _, ok := tm.builtinScalars[name]; ok {
		return nil
	}
	if _, ok := tm.customScalars[name]; ok {
		return nil
	}
	if _, ok := CommonScalars()[name]; ok {
		return nil
	}
	if _, ok := tm.schemaTypes[name]; ok {
		return nil
	}

	// Unknown type - this is a warning, not an error
	// as it might be defined elsewhere or be valid
	return errors.NewTypeMappingError(
		"unknown type, will use as-is",
		nil,
	).WithSourceType(name)
}
