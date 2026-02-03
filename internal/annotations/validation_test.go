package annotations

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/source-c/go-gql2j/internal/config"
	"github.com/source-c/go-gql2j/internal/parser"
)

func TestNewValidationGenerator(t *testing.T) {
	cfg := &config.ValidationConfig{Enabled: true}
	gen := NewValidationGenerator(cfg)

	require.NotNil(t, gen)
	assert.Equal(t, cfg, gen.config)
}

func TestValidationGenerator_GetValidationPackage_Jakarta(t *testing.T) {
	cfg := &config.ValidationConfig{
		Enabled: true,
		Package: config.ValidationJakarta,
	}
	gen := NewValidationGenerator(cfg)

	assert.Equal(t, "jakarta.validation.constraints", gen.GetValidationPackage())
}

func TestValidationGenerator_GetValidationPackage_Javax(t *testing.T) {
	cfg := &config.ValidationConfig{
		Enabled: true,
		Package: config.ValidationJavax,
	}
	gen := NewValidationGenerator(cfg)

	assert.Equal(t, "javax.validation.constraints", gen.GetValidationPackage())
}

func TestValidationGenerator_GenerateFieldAnnotations_Disabled(t *testing.T) {
	cfg := &config.ValidationConfig{Enabled: false}
	gen := NewValidationGenerator(cfg)

	field := &parser.FieldDef{Name: "name"}
	annotations, imports := gen.GenerateFieldAnnotations(field, true)

	assert.Empty(t, annotations)
	assert.Empty(t, imports)
}

func TestValidationGenerator_GenerateFieldAnnotations_NotNull(t *testing.T) {
	cfg := &config.ValidationConfig{
		Enabled:          true,
		Package:          config.ValidationJakarta,
		NotNullOnNonNull: true,
	}
	gen := NewValidationGenerator(cfg)

	field := &parser.FieldDef{Name: "name"}
	annotations, imports := gen.GenerateFieldAnnotations(field, true)

	assert.Contains(t, annotations, "@NotNull")
	assert.Contains(t, imports, "jakarta.validation.constraints.NotNull")
}

func TestValidationGenerator_GenerateFieldAnnotations_NullableField(t *testing.T) {
	cfg := &config.ValidationConfig{
		Enabled:          true,
		Package:          config.ValidationJakarta,
		NotNullOnNonNull: true,
	}
	gen := NewValidationGenerator(cfg)

	field := &parser.FieldDef{Name: "email"}
	annotations, imports := gen.GenerateFieldAnnotations(field, false) // nullable

	// Should not add @NotNull for nullable fields
	assert.NotContains(t, annotations, "@NotNull")
	assert.Empty(t, imports)
}

func TestValidationGenerator_GenerateFieldAnnotations_WithConstraintDirective_Size(t *testing.T) {
	cfg := &config.ValidationConfig{
		Enabled: true,
		Package: config.ValidationJakarta,
	}
	gen := NewValidationGenerator(cfg)

	minLen := 1
	maxLen := 255
	field := &parser.FieldDef{
		Name: "email",
		Directives: []*parser.DirectiveDef{
			{
				Name: "constraint",
				Arguments: map[string]interface{}{
					"minLength": int64(minLen),
					"maxLength": int64(maxLen),
				},
			},
		},
	}
	annotations, imports := gen.GenerateFieldAnnotations(field, false)

	assert.Contains(t, annotations, "@Size(min = 1, max = 255)")
	assert.Contains(t, imports, "jakarta.validation.constraints.Size")
}

func TestValidationGenerator_GenerateFieldAnnotations_WithConstraintDirective_MinMax(t *testing.T) {
	cfg := &config.ValidationConfig{
		Enabled: true,
		Package: config.ValidationJakarta,
	}
	gen := NewValidationGenerator(cfg)

	field := &parser.FieldDef{
		Name: "age",
		Directives: []*parser.DirectiveDef{
			{
				Name: "constraint",
				Arguments: map[string]interface{}{
					"min": int64(0),
					"max": int64(150),
				},
			},
		},
	}
	annotations, imports := gen.GenerateFieldAnnotations(field, false)

	assert.Contains(t, annotations, "@Min(0)")
	assert.Contains(t, annotations, "@Max(150)")
	assert.Contains(t, imports, "jakarta.validation.constraints.Min")
	assert.Contains(t, imports, "jakarta.validation.constraints.Max")
}

func TestValidationGenerator_GenerateFieldAnnotations_WithConstraintDirective_Pattern(t *testing.T) {
	cfg := &config.ValidationConfig{
		Enabled: true,
		Package: config.ValidationJakarta,
	}
	gen := NewValidationGenerator(cfg)

	field := &parser.FieldDef{
		Name: "code",
		Directives: []*parser.DirectiveDef{
			{
				Name: "constraint",
				Arguments: map[string]interface{}{
					"pattern": "^[A-Z0-9]+$",
				},
			},
		},
	}
	annotations, imports := gen.GenerateFieldAnnotations(field, false)

	assert.Contains(t, annotations, "@Pattern(regexp = \"^[A-Z0-9]+$\")")
	assert.Contains(t, imports, "jakarta.validation.constraints.Pattern")
}

func TestValidationGenerator_GenerateFieldAnnotations_WithConstraintDirective_Email(t *testing.T) {
	cfg := &config.ValidationConfig{
		Enabled: true,
		Package: config.ValidationJakarta,
	}
	gen := NewValidationGenerator(cfg)

	field := &parser.FieldDef{
		Name: "email",
		Directives: []*parser.DirectiveDef{
			{
				Name: "constraint",
				Arguments: map[string]interface{}{
					"email": true,
				},
			},
		},
	}
	annotations, imports := gen.GenerateFieldAnnotations(field, false)

	assert.Contains(t, annotations, "@Email")
	assert.Contains(t, imports, "jakarta.validation.constraints.Email")
}

func TestValidationGenerator_GenerateFieldAnnotations_WithConstraintDirective_NotBlank(t *testing.T) {
	cfg := &config.ValidationConfig{
		Enabled: true,
		Package: config.ValidationJakarta,
	}
	gen := NewValidationGenerator(cfg)

	field := &parser.FieldDef{
		Name: "name",
		Directives: []*parser.DirectiveDef{
			{
				Name: "constraint",
				Arguments: map[string]interface{}{
					"notBlank": true,
				},
			},
		},
	}
	annotations, imports := gen.GenerateFieldAnnotations(field, false)

	assert.Contains(t, annotations, "@NotBlank")
	assert.Contains(t, imports, "jakarta.validation.constraints.NotBlank")
}

func TestValidationGenerator_GenerateFieldAnnotations_WithConstraintDirective_NotNull(t *testing.T) {
	cfg := &config.ValidationConfig{
		Enabled: true,
		Package: config.ValidationJakarta,
	}
	gen := NewValidationGenerator(cfg)

	field := &parser.FieldDef{
		Name: "required",
		Directives: []*parser.DirectiveDef{
			{
				Name: "constraint",
				Arguments: map[string]interface{}{
					"notNull": true,
				},
			},
		},
	}
	annotations, imports := gen.GenerateFieldAnnotations(field, false)

	assert.Contains(t, annotations, "@NotNull")
	assert.Contains(t, imports, "jakarta.validation.constraints.NotNull")
}

func TestValidationGenerator_GenerateFieldAnnotations_MultipleConstraints(t *testing.T) {
	cfg := &config.ValidationConfig{
		Enabled:          true,
		Package:          config.ValidationJakarta,
		NotNullOnNonNull: true,
	}
	gen := NewValidationGenerator(cfg)

	field := &parser.FieldDef{
		Name: "email",
		Directives: []*parser.DirectiveDef{
			{
				Name: "constraint",
				Arguments: map[string]interface{}{
					"email":     true,
					"maxLength": int64(255),
				},
			},
		},
	}
	annotations, _ := gen.GenerateFieldAnnotations(field, true)

	// Should have @NotNull (from nonNull), @Email, and @Size
	assert.Contains(t, annotations, "@NotNull")
	assert.Contains(t, annotations, "@Email")
	assert.Contains(t, annotations, "@Size(max = 255)")
}

func TestValidationGenerator_GenerateFieldAnnotations_Javax(t *testing.T) {
	cfg := &config.ValidationConfig{
		Enabled:          true,
		Package:          config.ValidationJavax,
		NotNullOnNonNull: true,
	}
	gen := NewValidationGenerator(cfg)

	field := &parser.FieldDef{Name: "name"}
	annotations, imports := gen.GenerateFieldAnnotations(field, true)

	assert.Contains(t, annotations, "@NotNull")
	assert.Contains(t, imports, "javax.validation.constraints.NotNull")
}

func TestValidationGenerator_GenerateTypeAnnotations(t *testing.T) {
	cfg := &config.ValidationConfig{Enabled: true}
	gen := NewValidationGenerator(cfg)

	typeDef := &parser.TypeDef{Name: "User"}
	annotations, _ := gen.GenerateTypeAnnotations(typeDef)

	// Currently no type-level validation annotations
	assert.Empty(t, annotations)
}
