package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/source-c/go-gql2j/internal/config"
	"github.com/source-c/go-gql2j/internal/parser"
)

func TestEnumGenerator_Generate_BasicEnum(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)

	typeDef := &parser.TypeDef{
		Name: "Status",
		Kind: parser.TypeKindEnum,
		EnumValues: []*parser.EnumValueDef{
			{Name: "ACTIVE"},
			{Name: "INACTIVE"},
			{Name: "PENDING"},
		},
	}

	gen := NewEnumGenerator()
	content, err := gen.Generate(ctx, typeDef)
	require.NoError(t, err)

	assert.Contains(t, content, "package com.test;")
	assert.Contains(t, content, "public enum Status {")
	assert.Contains(t, content, "ACTIVE,")
	assert.Contains(t, content, "INACTIVE,")
	assert.Contains(t, content, "PENDING;")
}

func TestEnumGenerator_Generate_WithDescription(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)

	typeDef := &parser.TypeDef{
		Name:        "Status",
		Kind:        parser.TypeKindEnum,
		Description: "The status of an entity",
		EnumValues: []*parser.EnumValueDef{
			{Name: "ACTIVE", Description: "Entity is active"},
			{Name: "INACTIVE"},
		},
	}

	gen := NewEnumGenerator()
	content, err := gen.Generate(ctx, typeDef)
	require.NoError(t, err)

	assert.Contains(t, content, "/**")
	assert.Contains(t, content, "The status of an entity")
	assert.Contains(t, content, "Entity is active")
}

func TestEnumGenerator_Generate_WithDeprecated(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)

	typeDef := &parser.TypeDef{
		Name: "Status",
		Kind: parser.TypeKindEnum,
		EnumValues: []*parser.EnumValueDef{
			{Name: "ACTIVE"},
			{
				Name: "PENDING",
				Directives: []*parser.DirectiveDef{
					{Name: "deprecated", Arguments: map[string]interface{}{"reason": "Use ACTIVE"}},
				},
			},
		},
	}

	gen := NewEnumGenerator()
	content, err := gen.Generate(ctx, typeDef)
	require.NoError(t, err)

	assert.Contains(t, content, "@Deprecated")
	assert.Contains(t, content, "PENDING")
}

func TestEnumGenerator_Generate_SkipDirective(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)

	typeDef := &parser.TypeDef{
		Name:       "SkippedEnum",
		Kind:       parser.TypeKindEnum,
		Directives: []*parser.DirectiveDef{{Name: "skip"}},
		EnumValues: []*parser.EnumValueDef{
			{Name: "VALUE"},
		},
	}

	gen := NewEnumGenerator()
	content, err := gen.Generate(ctx, typeDef)
	require.NoError(t, err)
	assert.Empty(t, content)
}

func TestEnumGenerator_Generate_SkipValue(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)

	typeDef := &parser.TypeDef{
		Name: "Status",
		Kind: parser.TypeKindEnum,
		EnumValues: []*parser.EnumValueDef{
			{Name: "ACTIVE"},
			{
				Name:       "INTERNAL",
				Directives: []*parser.DirectiveDef{{Name: "skip"}},
			},
			{Name: "INACTIVE"},
		},
	}

	gen := NewEnumGenerator()
	content, err := gen.Generate(ctx, typeDef)
	require.NoError(t, err)

	assert.Contains(t, content, "ACTIVE")
	assert.Contains(t, content, "INACTIVE")
	assert.NotContains(t, content, "INTERNAL")
}

func TestEnumGenerator_Generate_WithJavaName(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)

	typeDef := &parser.TypeDef{
		Name: "Status",
		Kind: parser.TypeKindEnum,
		Directives: []*parser.DirectiveDef{
			{Name: "javaName", Arguments: map[string]interface{}{"name": "EntityStatus"}},
		},
		EnumValues: []*parser.EnumValueDef{
			{Name: "ACTIVE"},
		},
	}

	gen := NewEnumGenerator()
	content, err := gen.Generate(ctx, typeDef)
	require.NoError(t, err)

	assert.Contains(t, content, "public enum EntityStatus")
}

func TestEnumGenerator_Generate_DeprecatedEnum(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)

	typeDef := &parser.TypeDef{
		Name: "OldStatus",
		Kind: parser.TypeKindEnum,
		Directives: []*parser.DirectiveDef{
			{Name: "deprecated"},
		},
		EnumValues: []*parser.EnumValueDef{
			{Name: "VALUE"},
		},
	}

	gen := NewEnumGenerator()
	content, err := gen.Generate(ctx, typeDef)
	require.NoError(t, err)

	// @Deprecated should appear before the enum declaration
	deprecatedIdx := indexOf(content, "@Deprecated")
	enumIdx := indexOf(content, "public enum")
	assert.True(t, deprecatedIdx < enumIdx, "@Deprecated should come before enum declaration")
}

func TestEnumGenerator_Generate_SingleValue(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)

	typeDef := &parser.TypeDef{
		Name: "Singleton",
		Kind: parser.TypeKindEnum,
		EnumValues: []*parser.EnumValueDef{
			{Name: "INSTANCE"},
		},
	}

	gen := NewEnumGenerator()
	content, err := gen.Generate(ctx, typeDef)
	require.NoError(t, err)

	// Single value should end with semicolon, not comma
	assert.Contains(t, content, "INSTANCE;")
	assert.NotContains(t, content, "INSTANCE,")
}

// Helper function
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
