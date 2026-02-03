package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/source-c/go-gql2j/internal/config"
	"github.com/source-c/go-gql2j/internal/parser"
)

func TestInterfaceGenerator_Generate_BasicInterface(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)

	typeDef := &parser.TypeDef{
		Name: "Node",
		Kind: parser.TypeKindInterface,
		Fields: []*parser.FieldDef{
			{Name: "id", Type: &parser.TypeRef{Name: "ID", NonNull: true}},
		},
	}

	gen := NewInterfaceGenerator()
	content, err := gen.Generate(ctx, typeDef)
	require.NoError(t, err)

	assert.Contains(t, content, "package com.test;")
	assert.Contains(t, content, "public interface Node {")
	assert.Contains(t, content, "String getId();")
}

func TestInterfaceGenerator_Generate_WithDescription(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)

	typeDef := &parser.TypeDef{
		Name:        "Node",
		Kind:        parser.TypeKindInterface,
		Description: "A node in the graph",
		Fields: []*parser.FieldDef{
			{
				Name:        "id",
				Description: "The unique identifier",
				Type:        &parser.TypeRef{Name: "ID", NonNull: true},
			},
		},
	}

	gen := NewInterfaceGenerator()
	content, err := gen.Generate(ctx, typeDef)
	require.NoError(t, err)

	assert.Contains(t, content, "/**")
	assert.Contains(t, content, "A node in the graph")
	assert.Contains(t, content, "The unique identifier")
}

func TestInterfaceGenerator_Generate_WithExtends(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)

	typeDef := &parser.TypeDef{
		Name:       "Timestamped",
		Kind:       parser.TypeKindInterface,
		Interfaces: []string{"Node", "Auditable"},
		Fields: []*parser.FieldDef{
			{Name: "createdAt", Type: &parser.TypeRef{Name: "String", NonNull: true}},
		},
	}

	gen := NewInterfaceGenerator()
	content, err := gen.Generate(ctx, typeDef)
	require.NoError(t, err)

	assert.Contains(t, content, "extends Node, Auditable")
}

func TestInterfaceGenerator_Generate_WithInterfacePrefix(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	cfg.Java.Naming.InterfacePrefix = "I"
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)

	typeDef := &parser.TypeDef{
		Name:       "Timestamped",
		Kind:       parser.TypeKindInterface,
		Interfaces: []string{"Node"},
		Fields: []*parser.FieldDef{
			{Name: "createdAt", Type: &parser.TypeRef{Name: "String", NonNull: true}},
		},
	}

	gen := NewInterfaceGenerator()
	content, err := gen.Generate(ctx, typeDef)
	require.NoError(t, err)

	assert.Contains(t, content, "public interface ITimestamped")
	assert.Contains(t, content, "extends INode")
}

func TestInterfaceGenerator_Generate_SkipDirective(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)

	typeDef := &parser.TypeDef{
		Name:       "SkippedInterface",
		Kind:       parser.TypeKindInterface,
		Directives: []*parser.DirectiveDef{{Name: "skip"}},
		Fields: []*parser.FieldDef{
			{Name: "id", Type: &parser.TypeRef{Name: "ID", NonNull: true}},
		},
	}

	gen := NewInterfaceGenerator()
	content, err := gen.Generate(ctx, typeDef)
	require.NoError(t, err)
	assert.Empty(t, content)
}

func TestInterfaceGenerator_Generate_SkipField(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)

	typeDef := &parser.TypeDef{
		Name: "Node",
		Kind: parser.TypeKindInterface,
		Fields: []*parser.FieldDef{
			{Name: "id", Type: &parser.TypeRef{Name: "ID", NonNull: true}},
			{
				Name:       "internal",
				Type:       &parser.TypeRef{Name: "String", NonNull: true},
				Directives: []*parser.DirectiveDef{{Name: "skip"}},
			},
		},
	}

	gen := NewInterfaceGenerator()
	content, err := gen.Generate(ctx, typeDef)
	require.NoError(t, err)

	assert.Contains(t, content, "getId()")
	assert.NotContains(t, content, "getInternal()")
}

func TestInterfaceGenerator_Generate_WithDeprecated(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)

	typeDef := &parser.TypeDef{
		Name: "OldInterface",
		Kind: parser.TypeKindInterface,
		Directives: []*parser.DirectiveDef{
			{Name: "deprecated"},
		},
		Fields: []*parser.FieldDef{
			{Name: "id", Type: &parser.TypeRef{Name: "ID", NonNull: true}},
		},
	}

	gen := NewInterfaceGenerator()
	content, err := gen.Generate(ctx, typeDef)
	require.NoError(t, err)

	assert.Contains(t, content, "@Deprecated")
}

func TestInterfaceGenerator_Generate_DeprecatedField(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)

	typeDef := &parser.TypeDef{
		Name: "Node",
		Kind: parser.TypeKindInterface,
		Fields: []*parser.FieldDef{
			{Name: "id", Type: &parser.TypeRef{Name: "ID", NonNull: true}},
			{
				Name: "oldField",
				Type: &parser.TypeRef{Name: "String", NonNull: true},
				Directives: []*parser.DirectiveDef{
					{Name: "deprecated", Arguments: map[string]interface{}{"reason": "Use newField"}},
				},
			},
		},
	}

	gen := NewInterfaceGenerator()
	content, err := gen.Generate(ctx, typeDef)
	require.NoError(t, err)

	assert.Contains(t, content, "@Deprecated")
	assert.Contains(t, content, "getOldField()")
}

func TestInterfaceGenerator_Generate_ListReturnType(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)

	typeDef := &parser.TypeDef{
		Name: "Container",
		Kind: parser.TypeKindInterface,
		Fields: []*parser.FieldDef{
			{
				Name: "items",
				Type: &parser.TypeRef{
					Elem:    &parser.TypeRef{Name: "String", NonNull: true},
					NonNull: true,
				},
			},
		},
	}

	gen := NewInterfaceGenerator()
	content, err := gen.Generate(ctx, typeDef)
	require.NoError(t, err)

	assert.Contains(t, content, "List<String> getItems();")
	assert.Contains(t, content, "import java.util.List;")
}

func TestInterfaceGenerator_Generate_BooleanGetter(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)

	typeDef := &parser.TypeDef{
		Name: "Activatable",
		Kind: parser.TypeKindInterface,
		Fields: []*parser.FieldDef{
			{Name: "active", Type: &parser.TypeRef{Name: "Boolean", NonNull: true}},
		},
	}

	gen := NewInterfaceGenerator()
	content, err := gen.Generate(ctx, typeDef)
	require.NoError(t, err)

	// Boolean getters should use "is" prefix
	assert.Contains(t, content, "isActive()")
}

func TestInterfaceGenerator_Generate_MultipleFields(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)

	typeDef := &parser.TypeDef{
		Name: "Entity",
		Kind: parser.TypeKindInterface,
		Fields: []*parser.FieldDef{
			{Name: "id", Type: &parser.TypeRef{Name: "ID", NonNull: true}},
			{Name: "name", Type: &parser.TypeRef{Name: "String", NonNull: true}},
			{Name: "description", Type: &parser.TypeRef{Name: "String", NonNull: false}},
		},
	}

	gen := NewInterfaceGenerator()
	content, err := gen.Generate(ctx, typeDef)
	require.NoError(t, err)

	assert.Contains(t, content, "String getId();")
	assert.Contains(t, content, "String getName();")
	assert.Contains(t, content, "String getDescription();")
}

func TestInterfaceGenerator_Generate_WithJavaName(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)

	typeDef := &parser.TypeDef{
		Name: "Node",
		Kind: parser.TypeKindInterface,
		Directives: []*parser.DirectiveDef{
			{Name: "javaName", Arguments: map[string]interface{}{"name": "INode"}},
		},
		Fields: []*parser.FieldDef{
			{Name: "id", Type: &parser.TypeRef{Name: "ID", NonNull: true}},
		},
	}

	gen := NewInterfaceGenerator()
	content, err := gen.Generate(ctx, typeDef)
	require.NoError(t, err)

	assert.Contains(t, content, "public interface INode")
}
