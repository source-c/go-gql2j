package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/source-c/go-gql2j/internal/config"
	"github.com/source-c/go-gql2j/internal/parser"
)

func TestClassGenerator_Generate_BasicClass(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)

	typeDef := &parser.TypeDef{
		Name: "User",
		Kind: parser.TypeKindObject,
		Fields: []*parser.FieldDef{
			{Name: "id", Type: &parser.TypeRef{Name: "ID", NonNull: true}},
			{Name: "name", Type: &parser.TypeRef{Name: "String", NonNull: true}},
			{Name: "email", Type: &parser.TypeRef{Name: "String", NonNull: false}},
		},
	}

	gen := NewClassGenerator()
	content, err := gen.Generate(ctx, typeDef)
	require.NoError(t, err)

	assert.Contains(t, content, "package com.test;")
	assert.Contains(t, content, "public class User {")
	assert.Contains(t, content, "private String id;")
	assert.Contains(t, content, "private String name;")
	assert.Contains(t, content, "private String email;")
}

func TestClassGenerator_Generate_WithDescription(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)

	typeDef := &parser.TypeDef{
		Name:        "User",
		Kind:        parser.TypeKindObject,
		Description: "A user in the system",
		Fields: []*parser.FieldDef{
			{
				Name:        "id",
				Description: "The unique identifier",
				Type:        &parser.TypeRef{Name: "ID", NonNull: true},
			},
		},
	}

	gen := NewClassGenerator()
	content, err := gen.Generate(ctx, typeDef)
	require.NoError(t, err)

	assert.Contains(t, content, "/**")
	assert.Contains(t, content, "A user in the system")
	assert.Contains(t, content, "The unique identifier")
}

func TestClassGenerator_Generate_WithImplements(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)

	typeDef := &parser.TypeDef{
		Name:       "User",
		Kind:       parser.TypeKindObject,
		Interfaces: []string{"Node", "Timestamped"},
		Fields: []*parser.FieldDef{
			{Name: "id", Type: &parser.TypeRef{Name: "ID", NonNull: true}},
		},
	}

	gen := NewClassGenerator()
	content, err := gen.Generate(ctx, typeDef)
	require.NoError(t, err)

	assert.Contains(t, content, "implements Node, Timestamped")
}

func TestClassGenerator_Generate_WithInterfacePrefix(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	cfg.Java.Naming.InterfacePrefix = "I"
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)

	typeDef := &parser.TypeDef{
		Name:       "User",
		Kind:       parser.TypeKindObject,
		Interfaces: []string{"Node"},
		Fields: []*parser.FieldDef{
			{Name: "id", Type: &parser.TypeRef{Name: "ID", NonNull: true}},
		},
	}

	gen := NewClassGenerator()
	content, err := gen.Generate(ctx, typeDef)
	require.NoError(t, err)

	assert.Contains(t, content, "implements INode")
}

func TestClassGenerator_Generate_WithLombok(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	cfg.Features.Lombok.Enabled = true
	cfg.Features.Lombok.Data = true
	cfg.Features.Lombok.Builder = true
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)

	typeDef := &parser.TypeDef{
		Name: "User",
		Kind: parser.TypeKindObject,
		Fields: []*parser.FieldDef{
			{Name: "id", Type: &parser.TypeRef{Name: "ID", NonNull: true}},
		},
	}

	gen := NewClassGenerator()
	content, err := gen.Generate(ctx, typeDef)
	require.NoError(t, err)

	assert.Contains(t, content, "@Data")
	assert.Contains(t, content, "@Builder")
	assert.Contains(t, content, "import lombok.Data;")
	assert.Contains(t, content, "import lombok.Builder;")
}

func TestClassGenerator_Generate_WithValidation(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	cfg.Features.Validation.Enabled = true
	cfg.Features.Validation.Package = "jakarta"
	cfg.Features.Validation.NotNullOnNonNull = true
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)

	typeDef := &parser.TypeDef{
		Name: "User",
		Kind: parser.TypeKindObject,
		Fields: []*parser.FieldDef{
			{Name: "name", Type: &parser.TypeRef{Name: "String", NonNull: true}},
		},
	}

	gen := NewClassGenerator()
	content, err := gen.Generate(ctx, typeDef)
	require.NoError(t, err)

	assert.Contains(t, content, "@NotNull")
	assert.Contains(t, content, "import jakarta.validation.constraints.NotNull;")
}

func TestClassGenerator_Generate_SkipDirective(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)

	typeDef := &parser.TypeDef{
		Name:       "SkippedType",
		Kind:       parser.TypeKindObject,
		Directives: []*parser.DirectiveDef{{Name: "skip"}},
		Fields: []*parser.FieldDef{
			{Name: "id", Type: &parser.TypeRef{Name: "ID", NonNull: true}},
		},
	}

	gen := NewClassGenerator()
	content, err := gen.Generate(ctx, typeDef)
	require.NoError(t, err)
	assert.Empty(t, content)
}

func TestClassGenerator_Generate_SkipField(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)

	typeDef := &parser.TypeDef{
		Name: "User",
		Kind: parser.TypeKindObject,
		Fields: []*parser.FieldDef{
			{Name: "id", Type: &parser.TypeRef{Name: "ID", NonNull: true}},
			{
				Name:       "internal",
				Type:       &parser.TypeRef{Name: "String", NonNull: true},
				Directives: []*parser.DirectiveDef{{Name: "skip"}},
			},
		},
	}

	gen := NewClassGenerator()
	content, err := gen.Generate(ctx, typeDef)
	require.NoError(t, err)

	assert.Contains(t, content, "private String id;")
	assert.NotContains(t, content, "internal")
}

func TestClassGenerator_Generate_WithJavaNameDirective(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)

	typeDef := &parser.TypeDef{
		Name: "User",
		Kind: parser.TypeKindObject,
		Directives: []*parser.DirectiveDef{
			{
				Name:      "javaName",
				Arguments: map[string]interface{}{"name": "UserEntity"},
			},
		},
		Fields: []*parser.FieldDef{
			{Name: "id", Type: &parser.TypeRef{Name: "ID", NonNull: true}},
		},
	}

	gen := NewClassGenerator()
	content, err := gen.Generate(ctx, typeDef)
	require.NoError(t, err)

	assert.Contains(t, content, "public class UserEntity")
}

func TestClassGenerator_Generate_WithDeprecated(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)

	typeDef := &parser.TypeDef{
		Name: "OldUser",
		Kind: parser.TypeKindObject,
		Directives: []*parser.DirectiveDef{
			{Name: "deprecated", Arguments: map[string]interface{}{"reason": "Use User instead"}},
		},
		Fields: []*parser.FieldDef{
			{Name: "id", Type: &parser.TypeRef{Name: "ID", NonNull: true}},
		},
	}

	gen := NewClassGenerator()
	content, err := gen.Generate(ctx, typeDef)
	require.NoError(t, err)

	assert.Contains(t, content, "@Deprecated")
}

func TestClassGenerator_Generate_ListFields(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)

	typeDef := &parser.TypeDef{
		Name: "Post",
		Kind: parser.TypeKindObject,
		Fields: []*parser.FieldDef{
			{
				Name: "tags",
				Type: &parser.TypeRef{
					Elem:    &parser.TypeRef{Name: "String", NonNull: true},
					NonNull: true,
				},
			},
		},
	}

	gen := NewClassGenerator()
	content, err := gen.Generate(ctx, typeDef)
	require.NoError(t, err)

	assert.Contains(t, content, "List<String> tags")
	assert.Contains(t, content, "import java.util.List;")
}

func TestClassGenerator_Generate_PublicVisibility(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	cfg.Java.FieldVisibility = config.VisibilityPublic
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)

	typeDef := &parser.TypeDef{
		Name: "User",
		Kind: parser.TypeKindObject,
		Fields: []*parser.FieldDef{
			{Name: "name", Type: &parser.TypeRef{Name: "String", NonNull: true}},
		},
	}

	gen := NewClassGenerator()
	content, err := gen.Generate(ctx, typeDef)
	require.NoError(t, err)

	assert.Contains(t, content, "public String name;")
}

func TestClassGenerator_Generate_InputObject(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)

	typeDef := &parser.TypeDef{
		Name: "CreateUserInput",
		Kind: parser.TypeKindInputObject,
		Fields: []*parser.FieldDef{
			{Name: "name", Type: &parser.TypeRef{Name: "String", NonNull: true}},
			{Name: "email", Type: &parser.TypeRef{Name: "String", NonNull: false}},
		},
	}

	gen := NewClassGenerator()
	content, err := gen.Generate(ctx, typeDef)
	require.NoError(t, err)

	assert.Contains(t, content, "public class CreateUserInput")
}
