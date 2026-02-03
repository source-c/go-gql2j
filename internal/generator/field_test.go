package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/source-c/go-gql2j/internal/config"
	"github.com/source-c/go-gql2j/internal/parser"
)

func TestFieldGenerator_GenerateField_Basic(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)
	typeDef := &parser.TypeDef{Name: "User", Kind: parser.TypeKindObject}
	tc := NewTypeContext(ctx, typeDef)

	field := &parser.FieldDef{
		Name: "name",
		Type: &parser.TypeRef{Name: "String", NonNull: true},
	}
	fc, err := NewFieldContext(tc, field)
	require.NoError(t, err)

	gen := NewFieldGenerator()
	content := gen.GenerateField(fc)

	assert.Contains(t, content, "private String name;")
}

func TestFieldGenerator_GenerateField_WithDescription(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)
	typeDef := &parser.TypeDef{Name: "User", Kind: parser.TypeKindObject}
	tc := NewTypeContext(ctx, typeDef)

	field := &parser.FieldDef{
		Name:        "email",
		Description: "User's email address",
		Type:        &parser.TypeRef{Name: "String", NonNull: true},
	}
	fc, err := NewFieldContext(tc, field)
	require.NoError(t, err)

	gen := NewFieldGenerator()
	content := gen.GenerateField(fc)

	assert.Contains(t, content, "/**")
	assert.Contains(t, content, "User's email address")
	assert.Contains(t, content, "*/")
}

func TestFieldGenerator_GenerateGetter(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)
	typeDef := &parser.TypeDef{Name: "User", Kind: parser.TypeKindObject}
	tc := NewTypeContext(ctx, typeDef)

	field := &parser.FieldDef{
		Name: "name",
		Type: &parser.TypeRef{Name: "String", NonNull: true},
	}
	fc, err := NewFieldContext(tc, field)
	require.NoError(t, err)

	gen := NewFieldGenerator()
	content := gen.GenerateGetter(fc)

	assert.Contains(t, content, "public String getName()")
	assert.Contains(t, content, "return this.name;")
}

func TestFieldGenerator_GenerateSetter(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)
	typeDef := &parser.TypeDef{Name: "User", Kind: parser.TypeKindObject}
	tc := NewTypeContext(ctx, typeDef)

	field := &parser.FieldDef{
		Name: "name",
		Type: &parser.TypeRef{Name: "String", NonNull: true},
	}
	fc, err := NewFieldContext(tc, field)
	require.NoError(t, err)

	gen := NewFieldGenerator()
	content := gen.GenerateSetter(fc)

	assert.Contains(t, content, "public void setName(String name)")
	assert.Contains(t, content, "this.name = name;")
}

func TestFieldGenerator_GenerateGetter_Boolean(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)
	typeDef := &parser.TypeDef{Name: "User", Kind: parser.TypeKindObject}
	tc := NewTypeContext(ctx, typeDef)

	field := &parser.FieldDef{
		Name: "active",
		Type: &parser.TypeRef{Name: "Boolean", NonNull: true},
	}
	fc, err := NewFieldContext(tc, field)
	require.NoError(t, err)

	gen := NewFieldGenerator()
	content := gen.GenerateGetter(fc)

	// Boolean getters should use "is" prefix
	assert.Contains(t, content, "public boolean isActive()")
}

func TestFieldGenerator_GenerateInterfaceMethod(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)
	typeDef := &parser.TypeDef{Name: "Node", Kind: parser.TypeKindInterface}
	tc := NewTypeContext(ctx, typeDef)

	field := &parser.FieldDef{
		Name: "id",
		Type: &parser.TypeRef{Name: "ID", NonNull: true},
	}

	gen := NewFieldGenerator()
	content, err := gen.GenerateInterfaceMethod(tc, field)
	require.NoError(t, err)

	assert.Contains(t, content, "String getId();")
}

func TestFieldGenerator_GenerateInterfaceMethod_WithDescription(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)
	typeDef := &parser.TypeDef{Name: "Node", Kind: parser.TypeKindInterface}
	tc := NewTypeContext(ctx, typeDef)

	field := &parser.FieldDef{
		Name:        "id",
		Description: "Unique identifier",
		Type:        &parser.TypeRef{Name: "ID", NonNull: true},
	}

	gen := NewFieldGenerator()
	content, err := gen.GenerateInterfaceMethod(tc, field)
	require.NoError(t, err)

	assert.Contains(t, content, "/**")
	assert.Contains(t, content, "Unique identifier")
}

func TestFieldGenerator_GenerateInterfaceMethod_Skip(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)
	typeDef := &parser.TypeDef{Name: "Node", Kind: parser.TypeKindInterface}
	tc := NewTypeContext(ctx, typeDef)

	field := &parser.FieldDef{
		Name:       "internal",
		Type:       &parser.TypeRef{Name: "String", NonNull: true},
		Directives: []*parser.DirectiveDef{{Name: "skip"}},
	}

	gen := NewFieldGenerator()
	content, err := gen.GenerateInterfaceMethod(tc, field)
	require.NoError(t, err)
	assert.Empty(t, content)
}

func TestFieldGenerator_GenerateInterfaceMethod_Deprecated(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)
	typeDef := &parser.TypeDef{Name: "Node", Kind: parser.TypeKindInterface}
	tc := NewTypeContext(ctx, typeDef)

	field := &parser.FieldDef{
		Name: "oldField",
		Type: &parser.TypeRef{Name: "String", NonNull: true},
		Directives: []*parser.DirectiveDef{
			{Name: "deprecated"},
		},
	}

	gen := NewFieldGenerator()
	content, err := gen.GenerateInterfaceMethod(tc, field)
	require.NoError(t, err)

	assert.Contains(t, content, "@Deprecated")
}

func TestFieldGenerator_GenerateField_WithValidation(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	cfg.Features.Validation.Enabled = true
	cfg.Features.Validation.NotNullOnNonNull = true
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)
	typeDef := &parser.TypeDef{Name: "User", Kind: parser.TypeKindObject}
	tc := NewTypeContext(ctx, typeDef)

	field := &parser.FieldDef{
		Name: "name",
		Type: &parser.TypeRef{Name: "String", NonNull: true},
	}
	fc, err := NewFieldContext(tc, field)
	require.NoError(t, err)

	gen := NewFieldGenerator()
	content := gen.GenerateField(fc)

	assert.Contains(t, content, "@NotNull")
}

func TestFieldGenerator_GenerateField_WithDeprecated(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)
	typeDef := &parser.TypeDef{Name: "User", Kind: parser.TypeKindObject}
	tc := NewTypeContext(ctx, typeDef)

	field := &parser.FieldDef{
		Name: "oldField",
		Type: &parser.TypeRef{Name: "String", NonNull: true},
		Directives: []*parser.DirectiveDef{
			{Name: "deprecated"},
		},
	}
	fc, err := NewFieldContext(tc, field)
	require.NoError(t, err)

	gen := NewFieldGenerator()
	content := gen.GenerateField(fc)

	assert.Contains(t, content, "@Deprecated")
}

func TestFieldGenerator_GenerateField_ListType(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)
	typeDef := &parser.TypeDef{Name: "Post", Kind: parser.TypeKindObject}
	tc := NewTypeContext(ctx, typeDef)

	field := &parser.FieldDef{
		Name: "tags",
		Type: &parser.TypeRef{
			Elem:    &parser.TypeRef{Name: "String", NonNull: true},
			NonNull: true,
		},
	}
	fc, err := NewFieldContext(tc, field)
	require.NoError(t, err)

	gen := NewFieldGenerator()
	content := gen.GenerateField(fc)

	assert.Contains(t, content, "List<String> tags;")
}

func TestFieldGenerator_GenerateField_PublicVisibility(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	cfg.Java.FieldVisibility = config.VisibilityPublic
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)
	typeDef := &parser.TypeDef{Name: "User", Kind: parser.TypeKindObject}
	tc := NewTypeContext(ctx, typeDef)

	field := &parser.FieldDef{
		Name: "name",
		Type: &parser.TypeRef{Name: "String", NonNull: true},
	}
	fc, err := NewFieldContext(tc, field)
	require.NoError(t, err)

	gen := NewFieldGenerator()
	content := gen.GenerateField(fc)

	assert.Contains(t, content, "public String name;")
}

func TestFieldGenerator_GenerateField_PackageVisibility(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	cfg.Java.FieldVisibility = config.VisibilityPackage
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)
	typeDef := &parser.TypeDef{Name: "User", Kind: parser.TypeKindObject}
	tc := NewTypeContext(ctx, typeDef)

	field := &parser.FieldDef{
		Name: "name",
		Type: &parser.TypeRef{Name: "String", NonNull: true},
	}
	fc, err := NewFieldContext(tc, field)
	require.NoError(t, err)

	gen := NewFieldGenerator()
	content := gen.GenerateField(fc)

	// Package-private has no modifier keyword
	assert.Contains(t, content, "    String name;")
	assert.NotContains(t, content, "private")
	assert.NotContains(t, content, "public")
}
