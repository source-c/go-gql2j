package typemap

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/source-c/go-gql2j/internal/config"
	"github.com/source-c/go-gql2j/internal/parser"
)

func TestNewTypeMapper(t *testing.T) {
	cfg := config.DefaultConfig()
	tm := NewTypeMapper(cfg)

	require.NotNil(t, tm)
	assert.NotNil(t, tm.builtinScalars)
	assert.NotNil(t, tm.customScalars)
}

func TestTypeMapper_MapType_BuiltinScalars(t *testing.T) {
	cfg := config.DefaultConfig()
	tm := NewTypeMapper(cfg)

	tests := []struct {
		name     string
		typeRef  *parser.TypeRef
		expected string
	}{
		{
			name:     "String non-null",
			typeRef:  &parser.TypeRef{Name: "String", NonNull: true},
			expected: "String",
		},
		{
			name:     "String nullable",
			typeRef:  &parser.TypeRef{Name: "String", NonNull: false},
			expected: "String",
		},
		{
			name:     "Int non-null",
			typeRef:  &parser.TypeRef{Name: "Int", NonNull: true},
			expected: "int",
		},
		{
			name:     "Int nullable",
			typeRef:  &parser.TypeRef{Name: "Int", NonNull: false},
			expected: "Integer",
		},
		{
			name:     "Float non-null",
			typeRef:  &parser.TypeRef{Name: "Float", NonNull: true},
			expected: "double",
		},
		{
			name:     "Boolean non-null",
			typeRef:  &parser.TypeRef{Name: "Boolean", NonNull: true},
			expected: "boolean",
		},
		{
			name:     "Boolean nullable",
			typeRef:  &parser.TypeRef{Name: "Boolean", NonNull: false},
			expected: "Boolean",
		},
		{
			name:     "ID non-null",
			typeRef:  &parser.TypeRef{Name: "ID", NonNull: true},
			expected: "String",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tm.MapType(tt.typeRef)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result.JavaType)
		})
	}
}

func TestTypeMapper_MapType_ListTypes(t *testing.T) {
	cfg := config.DefaultConfig()
	tm := NewTypeMapper(cfg)

	// [String!]!
	typeRef := &parser.TypeRef{
		Elem: &parser.TypeRef{
			Name:    "String",
			NonNull: true,
		},
		NonNull: true,
	}

	result, err := tm.MapType(typeRef)
	require.NoError(t, err)
	assert.Equal(t, "List<String>", result.JavaType)
	assert.True(t, result.IsCollection)
	assert.Equal(t, "String", result.ElementType)
}

func TestTypeMapper_MapType_ListOfInts(t *testing.T) {
	cfg := config.DefaultConfig()
	tm := NewTypeMapper(cfg)

	// [Int!]!
	typeRef := &parser.TypeRef{
		Elem: &parser.TypeRef{
			Name:    "Int",
			NonNull: true,
		},
		NonNull: true,
	}

	result, err := tm.MapType(typeRef)
	require.NoError(t, err)
	// Element types in collections should be boxed
	assert.Equal(t, "List<Integer>", result.JavaType)
	assert.True(t, result.IsCollection)
}

func TestTypeMapper_MapType_SetCollection(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Java.CollectionType = config.CollectionSet
	tm := NewTypeMapper(cfg)

	typeRef := &parser.TypeRef{
		Elem: &parser.TypeRef{
			Name:    "String",
			NonNull: true,
		},
		NonNull: true,
	}

	result, err := tm.MapType(typeRef)
	require.NoError(t, err)
	assert.Equal(t, "Set<String>", result.JavaType)
}

func TestTypeMapper_MapType_OptionalNullable(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Java.NullableHandling = config.NullableOptional
	tm := NewTypeMapper(cfg)

	// String (nullable)
	typeRef := &parser.TypeRef{
		Name:    "String",
		NonNull: false,
	}

	result, err := tm.MapType(typeRef)
	require.NoError(t, err)
	assert.Equal(t, "Optional<String>", result.JavaType)
	assert.True(t, result.IsOptional)
}

func TestTypeMapper_MapType_CustomScalars(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.TypeMappings.Scalars = map[string]config.ScalarMapping{
		"DateTime": {
			JavaType: "java.time.LocalDateTime",
			Imports:  []string{"java.time.LocalDateTime"},
		},
	}
	tm := NewTypeMapper(cfg)

	typeRef := &parser.TypeRef{
		Name:    "DateTime",
		NonNull: true,
	}

	result, err := tm.MapType(typeRef)
	require.NoError(t, err)
	assert.Equal(t, "java.time.LocalDateTime", result.JavaType)
	assert.Contains(t, result.Imports, "java.time.LocalDateTime")
}

func TestTypeMapper_MapType_SchemaTypes(t *testing.T) {
	cfg := config.DefaultConfig()
	tm := NewTypeMapper(cfg)

	// Set up schema types
	schemaTypes := map[string]*parser.TypeDef{
		"User": {
			Name: "User",
			Kind: parser.TypeKindObject,
		},
	}
	tm.SetSchemaTypes(schemaTypes)

	typeRef := &parser.TypeRef{
		Name:    "User",
		NonNull: true,
	}

	result, err := tm.MapType(typeRef)
	require.NoError(t, err)
	assert.Equal(t, "User", result.JavaType)
}

func TestTypeMapper_MapType_NilTypeRef(t *testing.T) {
	cfg := config.DefaultConfig()
	tm := NewTypeMapper(cfg)

	result, err := tm.MapType(nil)
	require.NoError(t, err)
	assert.Equal(t, "Object", result.JavaType)
}

func TestTypeMapper_MapFieldType_WithJavaTypeDirective(t *testing.T) {
	cfg := config.DefaultConfig()
	tm := NewTypeMapper(cfg)

	field := &parser.FieldDef{
		Name: "customField",
		Type: &parser.TypeRef{Name: "String", NonNull: true},
		Directives: []*parser.DirectiveDef{
			{
				Name: "javaType",
				Arguments: map[string]interface{}{
					"type":    "java.time.Instant",
					"imports": []interface{}{"java.time.Instant"},
				},
			},
		},
	}

	result, err := tm.MapFieldType(field)
	require.NoError(t, err)
	assert.Equal(t, "java.time.Instant", result.JavaType)
}

func TestTypeMapper_MapFieldType_WithCollectionDirective(t *testing.T) {
	cfg := config.DefaultConfig()
	tm := NewTypeMapper(cfg)

	field := &parser.FieldDef{
		Name: "tags",
		Type: &parser.TypeRef{
			Elem:    &parser.TypeRef{Name: "String", NonNull: true},
			NonNull: true,
		},
		Directives: []*parser.DirectiveDef{
			{
				Name: "collection",
				Arguments: map[string]interface{}{
					"type": "Set",
				},
			},
		},
	}

	result, err := tm.MapFieldType(field)
	require.NoError(t, err)
	assert.Equal(t, "Set<String>", result.JavaType)
}

func TestTypeMapper_ValidateMapping_KnownType(t *testing.T) {
	cfg := config.DefaultConfig()
	tm := NewTypeMapper(cfg)

	typeRef := &parser.TypeRef{Name: "String", NonNull: true}
	err := tm.ValidateMapping(typeRef)
	assert.NoError(t, err)
}

func TestTypeMapper_ValidateMapping_UnknownType(t *testing.T) {
	cfg := config.DefaultConfig()
	tm := NewTypeMapper(cfg)

	typeRef := &parser.TypeRef{Name: "UnknownType", NonNull: true}
	err := tm.ValidateMapping(typeRef)
	// Unknown types return a warning error, not nil
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown type")
}

func TestTypeMapper_ValidateMapping_NilTypeRef(t *testing.T) {
	cfg := config.DefaultConfig()
	tm := NewTypeMapper(cfg)

	err := tm.ValidateMapping(nil)
	assert.NoError(t, err)
}

func TestBoxType(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"int", "Integer"},
		{"long", "Long"},
		{"double", "Double"},
		{"float", "Float"},
		{"boolean", "Boolean"},
		{"byte", "Byte"},
		{"short", "Short"},
		{"char", "Character"},
		{"String", "String"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := BoxType(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuiltinScalars(t *testing.T) {
	scalars := BuiltinScalars()

	assert.Contains(t, scalars, "String")
	assert.Contains(t, scalars, "Int")
	assert.Contains(t, scalars, "Float")
	assert.Contains(t, scalars, "Boolean")
	assert.Contains(t, scalars, "ID")

	intScalar := scalars["Int"]
	assert.Equal(t, "Integer", intScalar.JavaType)
	assert.Equal(t, "int", intScalar.PrimitiveType)
}

func TestCommonScalars(t *testing.T) {
	scalars := CommonScalars()

	// Check common scalars exist
	assert.Contains(t, scalars, "DateTime")
	assert.Contains(t, scalars, "Date")
	assert.Contains(t, scalars, "UUID")
	assert.Contains(t, scalars, "BigDecimal")
}

func TestGetCollectionInfo(t *testing.T) {
	listInfo := GetCollectionInfo("List")
	assert.Equal(t, "java.util.List", listInfo.Imports[0])

	setInfo := GetCollectionInfo("Set")
	assert.Equal(t, "java.util.Set", setInfo.Imports[0])
}

func TestFormatCollectionType(t *testing.T) {
	assert.Equal(t, "List<String>", FormatCollectionType("List", "String"))
	assert.Equal(t, "Set<Integer>", FormatCollectionType("Set", "Integer"))
	assert.Equal(t, "Collection<User>", FormatCollectionType("Collection", "User"))
}

func TestFormatOptionalType(t *testing.T) {
	assert.Equal(t, "Optional<String>", FormatOptionalType("String"))
	assert.Equal(t, "Optional<Integer>", FormatOptionalType("Integer"))
}

func TestGetOptionalInfo(t *testing.T) {
	info := GetOptionalInfo()
	assert.Contains(t, info.Imports, "java.util.Optional")
}
