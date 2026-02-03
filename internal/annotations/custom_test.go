package annotations

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/source-c/go-gql2j/internal/parser"
)

func TestNewCustomAnnotationGenerator(t *testing.T) {
	gen := NewCustomAnnotationGenerator()
	require.NotNil(t, gen)
}

func TestCustomAnnotationGenerator_GenerateTypeAnnotations_NoDirectives(t *testing.T) {
	gen := NewCustomAnnotationGenerator()

	typeDef := &parser.TypeDef{Name: "User"}
	annotations, imports := gen.GenerateTypeAnnotations(typeDef)

	assert.Empty(t, annotations)
	assert.Empty(t, imports)
}

func TestCustomAnnotationGenerator_GenerateTypeAnnotations_WithAnnotation(t *testing.T) {
	gen := NewCustomAnnotationGenerator()

	typeDef := &parser.TypeDef{
		Name: "User",
		Directives: []*parser.DirectiveDef{
			{
				Name: "annotation",
				Arguments: map[string]interface{}{
					"value":   "@Entity",
					"imports": []interface{}{"javax.persistence.Entity"},
				},
			},
		},
	}
	annotations, imports := gen.GenerateTypeAnnotations(typeDef)

	assert.Contains(t, annotations, "@Entity")
	assert.Contains(t, imports, "javax.persistence.Entity")
}

func TestCustomAnnotationGenerator_GenerateTypeAnnotations_MultipleAnnotations(t *testing.T) {
	gen := NewCustomAnnotationGenerator()

	typeDef := &parser.TypeDef{
		Name: "User",
		Directives: []*parser.DirectiveDef{
			{
				Name: "annotation",
				Arguments: map[string]interface{}{
					"value":   "@Entity",
					"imports": []interface{}{"javax.persistence.Entity"},
				},
			},
			{
				Name: "annotation",
				Arguments: map[string]interface{}{
					"value":   "@Table(name = \"users\")",
					"imports": []interface{}{"javax.persistence.Table"},
				},
			},
		},
	}
	annotations, imports := gen.GenerateTypeAnnotations(typeDef)

	assert.Len(t, annotations, 2)
	assert.Contains(t, annotations, "@Entity")
	assert.Contains(t, annotations, "@Table(name = \"users\")")
	assert.Contains(t, imports, "javax.persistence.Entity")
	assert.Contains(t, imports, "javax.persistence.Table")
}

func TestCustomAnnotationGenerator_GenerateFieldAnnotations_NoDirectives(t *testing.T) {
	gen := NewCustomAnnotationGenerator()

	field := &parser.FieldDef{Name: "name"}
	annotations, imports := gen.GenerateFieldAnnotations(field)

	assert.Empty(t, annotations)
	assert.Empty(t, imports)
}

func TestCustomAnnotationGenerator_GenerateFieldAnnotations_WithAnnotation(t *testing.T) {
	gen := NewCustomAnnotationGenerator()

	field := &parser.FieldDef{
		Name: "id",
		Directives: []*parser.DirectiveDef{
			{
				Name: "annotation",
				Arguments: map[string]interface{}{
					"value":   "@Id",
					"imports": []interface{}{"javax.persistence.Id"},
				},
			},
		},
	}
	annotations, imports := gen.GenerateFieldAnnotations(field)

	assert.Contains(t, annotations, "@Id")
	assert.Contains(t, imports, "javax.persistence.Id")
}

func TestCustomAnnotationGenerator_GenerateEnumValueAnnotations_NoDirectives(t *testing.T) {
	gen := NewCustomAnnotationGenerator()

	enumValue := &parser.EnumValueDef{Name: "ACTIVE"}
	annotations, imports := gen.GenerateEnumValueAnnotations(enumValue)

	assert.Empty(t, annotations)
	assert.Empty(t, imports)
}

func TestCustomAnnotationGenerator_GenerateEnumValueAnnotations_WithAnnotation(t *testing.T) {
	gen := NewCustomAnnotationGenerator()

	enumValue := &parser.EnumValueDef{
		Name: "ACTIVE",
		Directives: []*parser.DirectiveDef{
			{
				Name: "annotation",
				Arguments: map[string]interface{}{
					"value":   "@JsonProperty(\"active\")",
					"imports": []interface{}{"com.fasterxml.jackson.annotation.JsonProperty"},
				},
			},
		},
	}
	annotations, imports := gen.GenerateEnumValueAnnotations(enumValue)

	assert.Contains(t, annotations, "@JsonProperty(\"active\")")
	assert.Contains(t, imports, "com.fasterxml.jackson.annotation.JsonProperty")
}

func TestCustomAnnotationGenerator_GenerateDeprecatedAnnotation_NoDirective(t *testing.T) {
	gen := NewCustomAnnotationGenerator()

	directives := []*parser.DirectiveDef{}
	annotation, _ := gen.GenerateDeprecatedAnnotation(directives)

	assert.Empty(t, annotation)
}

func TestCustomAnnotationGenerator_GenerateDeprecatedAnnotation_WithDirective(t *testing.T) {
	gen := NewCustomAnnotationGenerator()

	directives := []*parser.DirectiveDef{
		{
			Name: "deprecated",
			Arguments: map[string]interface{}{
				"reason": "Use newField instead",
			},
		},
	}
	annotation, _ := gen.GenerateDeprecatedAnnotation(directives)

	assert.Equal(t, "@Deprecated", annotation)
}

func TestCustomAnnotationGenerator_GenerateDeprecatedAnnotation_WithoutReason(t *testing.T) {
	gen := NewCustomAnnotationGenerator()

	directives := []*parser.DirectiveDef{
		{Name: "deprecated"},
	}
	annotation, _ := gen.GenerateDeprecatedAnnotation(directives)

	assert.Equal(t, "@Deprecated", annotation)
}
