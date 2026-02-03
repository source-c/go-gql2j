package generator

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/source-c/go-gql2j/internal/config"
	"github.com/source-c/go-gql2j/internal/parser"
)

func TestNewGenerator(t *testing.T) {
	cfg := config.DefaultConfig()
	gen := NewGenerator(cfg)

	require.NotNil(t, gen)
	assert.NotNil(t, gen.classGen)
	assert.NotNil(t, gen.interfaceGen)
	assert.NotNil(t, gen.enumGen)
}

func TestGenerator_Generate_BasicSchema(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	gen := NewGenerator(cfg)

	schema := &parser.Schema{
		Types: map[string]*parser.TypeDef{
			"User": {
				Name: "User",
				Kind: parser.TypeKindObject,
				Fields: []*parser.FieldDef{
					{Name: "id", Type: &parser.TypeRef{Name: "ID", NonNull: true}},
					{Name: "name", Type: &parser.TypeRef{Name: "String", NonNull: true}},
				},
			},
		},
	}

	files, err := gen.Generate(schema)
	require.NoError(t, err)
	require.Len(t, files, 1)

	userFile := files[0]
	assert.Equal(t, "User.java", userFile.FileName)
	assert.Contains(t, userFile.Content, "package com.test;")
	assert.Contains(t, userFile.Content, "public class User")
	assert.Contains(t, userFile.Content, "private String id;")
	assert.Contains(t, userFile.Content, "private String name;")
}

func TestGenerator_Generate_Interface(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	gen := NewGenerator(cfg)

	schema := &parser.Schema{
		Types: map[string]*parser.TypeDef{
			"Node": {
				Name: "Node",
				Kind: parser.TypeKindInterface,
				Fields: []*parser.FieldDef{
					{Name: "id", Type: &parser.TypeRef{Name: "ID", NonNull: true}},
				},
			},
		},
	}

	files, err := gen.Generate(schema)
	require.NoError(t, err)
	require.Len(t, files, 1)

	nodeFile := files[0]
	assert.Equal(t, "Node.java", nodeFile.FileName)
	assert.Contains(t, nodeFile.Content, "public interface Node")
	assert.Contains(t, nodeFile.Content, "String getId();")
}

func TestGenerator_Generate_Enum(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	gen := NewGenerator(cfg)

	schema := &parser.Schema{
		Types: map[string]*parser.TypeDef{
			"Status": {
				Name: "Status",
				Kind: parser.TypeKindEnum,
				EnumValues: []*parser.EnumValueDef{
					{Name: "ACTIVE"},
					{Name: "INACTIVE"},
				},
			},
		},
	}

	files, err := gen.Generate(schema)
	require.NoError(t, err)
	require.Len(t, files, 1)

	statusFile := files[0]
	assert.Equal(t, "Status.java", statusFile.FileName)
	assert.Contains(t, statusFile.Content, "public enum Status")
	assert.Contains(t, statusFile.Content, "ACTIVE")
	assert.Contains(t, statusFile.Content, "INACTIVE")
}

func TestGenerator_Generate_InputType(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	gen := NewGenerator(cfg)

	schema := &parser.Schema{
		Types: map[string]*parser.TypeDef{
			"CreateUserInput": {
				Name: "CreateUserInput",
				Kind: parser.TypeKindInputObject,
				Fields: []*parser.FieldDef{
					{Name: "name", Type: &parser.TypeRef{Name: "String", NonNull: true}},
				},
			},
		},
	}

	files, err := gen.Generate(schema)
	require.NoError(t, err)
	require.Len(t, files, 1)

	assert.Equal(t, "CreateUserInput.java", files[0].FileName)
	assert.Contains(t, files[0].Content, "public class CreateUserInput")
}

func TestGenerator_Generate_SkipsScalars(t *testing.T) {
	cfg := config.DefaultConfig()
	gen := NewGenerator(cfg)

	schema := &parser.Schema{
		Types: map[string]*parser.TypeDef{
			"CustomScalar": {
				Name: "CustomScalar",
				Kind: parser.TypeKindScalar,
			},
		},
	}

	files, err := gen.Generate(schema)
	require.NoError(t, err)
	assert.Empty(t, files)
}

func TestGenerator_Generate_UnionAsMarkerInterface(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	gen := NewGenerator(cfg)

	schema := &parser.Schema{
		Types: map[string]*parser.TypeDef{
			"SearchResult": {
				Name: "SearchResult",
				Kind: parser.TypeKindUnion,
			},
		},
	}

	files, err := gen.Generate(schema)
	require.NoError(t, err)
	require.Len(t, files, 1)
	assert.Equal(t, "SearchResult.java", files[0].FileName)
	assert.Contains(t, files[0].Content, "public interface SearchResult")
	assert.Contains(t, files[0].Content, "Marker interface")
}

func TestGenerator_Generate_WithSkipDirective(t *testing.T) {
	cfg := config.DefaultConfig()
	gen := NewGenerator(cfg)

	schema := &parser.Schema{
		Types: map[string]*parser.TypeDef{
			"SkippedType": {
				Name: "SkippedType",
				Kind: parser.TypeKindObject,
				Directives: []*parser.DirectiveDef{
					{Name: "skip"},
				},
			},
		},
	}

	files, err := gen.Generate(schema)
	require.NoError(t, err)
	assert.Empty(t, files)
}

func TestGenerator_GenerateType_NotFound(t *testing.T) {
	cfg := config.DefaultConfig()
	gen := NewGenerator(cfg)

	schema := &parser.Schema{
		Types: map[string]*parser.TypeDef{},
	}

	_, err := gen.GenerateType(schema, "NonexistentType")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "type not found")
}

func TestGenerator_GenerateWithResult(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	gen := NewGenerator(cfg)

	schema := &parser.Schema{
		Types: map[string]*parser.TypeDef{
			"User": {
				Name: "User",
				Kind: parser.TypeKindObject,
				Fields: []*parser.FieldDef{
					{Name: "id", Type: &parser.TypeRef{Name: "ID", NonNull: true}},
				},
			},
			"Status": {
				Name: "Status",
				Kind: parser.TypeKindEnum,
				EnumValues: []*parser.EnumValueDef{
					{Name: "ACTIVE"},
				},
			},
		},
	}

	result := gen.GenerateWithResult(schema)
	require.NotNil(t, result)
	assert.Len(t, result.Files, 2)
	assert.Empty(t, result.Errors)
}

func TestGetStats(t *testing.T) {
	files := []*GeneratedFile{
		{FileName: "User.java", TypeDef: &parser.TypeDef{Kind: parser.TypeKindObject}},
		{FileName: "Post.java", TypeDef: &parser.TypeDef{Kind: parser.TypeKindObject}},
		{FileName: "Node.java", TypeDef: &parser.TypeDef{Kind: parser.TypeKindInterface}},
		{FileName: "Status.java", TypeDef: &parser.TypeDef{Kind: parser.TypeKindEnum}},
		{FileName: "Input.java", TypeDef: &parser.TypeDef{Kind: parser.TypeKindInputObject}},
	}
	errors := []error{nil}

	stats := GetStats(files, errors)

	assert.Equal(t, 5, stats.TotalTypes)
	assert.Equal(t, 3, stats.Classes)
	assert.Equal(t, 1, stats.Interfaces)
	assert.Equal(t, 1, stats.Enums)
	assert.Equal(t, 1, stats.ErrorCount)
}

func TestNewContext(t *testing.T) {
	cfg := config.DefaultConfig()
	schema := &parser.Schema{
		Types: map[string]*parser.TypeDef{
			"User": {Name: "User", Kind: parser.TypeKindObject},
		},
	}

	ctx := NewContext(cfg, schema)

	require.NotNil(t, ctx)
	assert.NotNil(t, ctx.TypeMapper)
	assert.NotNil(t, ctx.NamingHelper)
	assert.NotNil(t, ctx.LombokGen)
	assert.NotNil(t, ctx.ValidationGen)
	assert.NotNil(t, ctx.CustomAnnotation)
}

func TestTypeContext_ShouldSkip(t *testing.T) {
	cfg := config.DefaultConfig()
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)

	// Type without skip directive
	typeDef := &parser.TypeDef{Name: "User", Kind: parser.TypeKindObject}
	tc := NewTypeContext(ctx, typeDef)
	assert.False(t, tc.ShouldSkip())

	// Type with skip directive
	skippedType := &parser.TypeDef{
		Name:       "Skipped",
		Kind:       parser.TypeKindObject,
		Directives: []*parser.DirectiveDef{{Name: "skip"}},
	}
	tc2 := NewTypeContext(ctx, skippedType)
	assert.True(t, tc2.ShouldSkip())
}

func TestTypeContext_GetVisibility(t *testing.T) {
	tests := []struct {
		visibility string
		expected   string
	}{
		{config.VisibilityPrivate, "private"},
		{config.VisibilityProtected, "protected"},
		{config.VisibilityPackage, ""},
		{config.VisibilityPublic, "public"},
	}

	for _, tt := range tests {
		t.Run(tt.visibility, func(t *testing.T) {
			cfg := config.DefaultConfig()
			cfg.Java.FieldVisibility = tt.visibility
			schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
			ctx := NewContext(cfg, schema)
			typeDef := &parser.TypeDef{Name: "Test", Kind: parser.TypeKindObject}
			tc := NewTypeContext(ctx, typeDef)

			assert.Equal(t, tt.expected, tc.GetVisibility())
		})
	}
}

func TestFieldContext_IsBooleanType(t *testing.T) {
	cfg := config.DefaultConfig()
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)
	typeDef := &parser.TypeDef{Name: "Test", Kind: parser.TypeKindObject}
	tc := NewTypeContext(ctx, typeDef)

	// Boolean field
	boolField := &parser.FieldDef{
		Name: "active",
		Type: &parser.TypeRef{Name: "Boolean", NonNull: true},
	}
	fc, err := NewFieldContext(tc, boolField)
	require.NoError(t, err)
	assert.True(t, fc.IsBooleanType())

	// String field
	stringField := &parser.FieldDef{
		Name: "name",
		Type: &parser.TypeRef{Name: "String", NonNull: true},
	}
	fc2, err := NewFieldContext(tc, stringField)
	require.NoError(t, err)
	assert.False(t, fc2.IsBooleanType())
}

func TestFieldContext_ShouldSkip(t *testing.T) {
	cfg := config.DefaultConfig()
	schema := &parser.Schema{Types: map[string]*parser.TypeDef{}}
	ctx := NewContext(cfg, schema)
	typeDef := &parser.TypeDef{Name: "Test", Kind: parser.TypeKindObject}
	tc := NewTypeContext(ctx, typeDef)

	// Field without skip
	field := &parser.FieldDef{
		Name: "name",
		Type: &parser.TypeRef{Name: "String", NonNull: true},
	}
	fc, err := NewFieldContext(tc, field)
	require.NoError(t, err)
	assert.False(t, fc.ShouldSkip())

	// Field with skip
	skippedField := &parser.FieldDef{
		Name:       "internal",
		Type:       &parser.TypeRef{Name: "String", NonNull: true},
		Directives: []*parser.DirectiveDef{{Name: "skip"}},
	}
	fc2, err := NewFieldContext(tc, skippedField)
	require.NoError(t, err)
	assert.True(t, fc2.ShouldSkip())
}

func TestEscapeJavaKeyword(t *testing.T) {
	// Java keywords should be escaped
	assert.Equal(t, "_class", EscapeJavaKeyword("class"))
	assert.Equal(t, "_public", EscapeJavaKeyword("public"))
	assert.Equal(t, "_final", EscapeJavaKeyword("final"))

	// Non-keywords should remain unchanged
	assert.Equal(t, "name", EscapeJavaKeyword("name"))
	assert.Equal(t, "email", EscapeJavaKeyword("email"))
}

func TestGenerator_GeneratesGettersAndSetters(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	// Lombok disabled, so getters/setters should be generated
	cfg.Features.Lombok.Enabled = false
	gen := NewGenerator(cfg)

	schema := &parser.Schema{
		Types: map[string]*parser.TypeDef{
			"User": {
				Name: "User",
				Kind: parser.TypeKindObject,
				Fields: []*parser.FieldDef{
					{Name: "name", Type: &parser.TypeRef{Name: "String", NonNull: true}},
				},
			},
		},
	}

	files, err := gen.Generate(schema)
	require.NoError(t, err)
	require.Len(t, files, 1)

	content := files[0].Content
	assert.Contains(t, content, "public String getName()")
	assert.Contains(t, content, "public void setName(String name)")
}

func TestGenerator_WithLombok(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.test"
	cfg.Features.Lombok.Enabled = true
	cfg.Features.Lombok.Data = true
	gen := NewGenerator(cfg)

	schema := &parser.Schema{
		Types: map[string]*parser.TypeDef{
			"User": {
				Name: "User",
				Kind: parser.TypeKindObject,
				Fields: []*parser.FieldDef{
					{Name: "name", Type: &parser.TypeRef{Name: "String", NonNull: true}},
				},
			},
		},
	}

	files, err := gen.Generate(schema)
	require.NoError(t, err)
	require.Len(t, files, 1)

	content := files[0].Content
	assert.Contains(t, content, "@Data")
	assert.Contains(t, content, "import lombok.Data;")
	// With @Data, manual getters/setters shouldn't be generated
	assert.False(t, strings.Contains(content, "public String getName()"))
}
