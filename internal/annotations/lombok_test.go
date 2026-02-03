package annotations

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/source-c/go-gql2j/internal/config"
	"github.com/source-c/go-gql2j/internal/parser"
)

func TestLombokAnnotations(t *testing.T) {
	anns := LombokAnnotations()

	assert.Contains(t, anns, "data")
	assert.Contains(t, anns, "builder")
	assert.Contains(t, anns, "noArgsConstructor")
	assert.Contains(t, anns, "allArgsConstructor")
	assert.Contains(t, anns, "getter")
	assert.Contains(t, anns, "setter")

	assert.Equal(t, "@Data", anns["data"].Name)
	assert.Equal(t, "lombok.Data", anns["data"].Import)
}

func TestNewLombokGenerator(t *testing.T) {
	cfg := &config.LombokConfig{Enabled: true}
	gen := NewLombokGenerator(cfg)

	require.NotNil(t, gen)
	assert.Equal(t, cfg, gen.config)
}

func TestLombokGenerator_GenerateTypeAnnotations_Disabled(t *testing.T) {
	cfg := &config.LombokConfig{Enabled: false}
	gen := NewLombokGenerator(cfg)

	typeDef := &parser.TypeDef{Name: "User"}
	annotations, imports := gen.GenerateTypeAnnotations(typeDef)

	assert.Empty(t, annotations)
	assert.Empty(t, imports)
}

func TestLombokGenerator_GenerateTypeAnnotations_DataEnabled(t *testing.T) {
	cfg := &config.LombokConfig{
		Enabled: true,
		Data:    true,
	}
	gen := NewLombokGenerator(cfg)

	typeDef := &parser.TypeDef{Name: "User"}
	annotations, imports := gen.GenerateTypeAnnotations(typeDef)

	assert.Contains(t, annotations, "@Data")
	assert.Contains(t, imports, "lombok.Data")
}

func TestLombokGenerator_GenerateTypeAnnotations_Multiple(t *testing.T) {
	cfg := &config.LombokConfig{
		Enabled:           true,
		Data:              true,
		Builder:           true,
		NoArgsConstructor: true,
	}
	gen := NewLombokGenerator(cfg)

	typeDef := &parser.TypeDef{Name: "User"}
	annotations, imports := gen.GenerateTypeAnnotations(typeDef)

	assert.Contains(t, annotations, "@Data")
	assert.Contains(t, annotations, "@Builder")
	assert.Contains(t, annotations, "@NoArgsConstructor")
	assert.Contains(t, imports, "lombok.Data")
	assert.Contains(t, imports, "lombok.Builder")
	assert.Contains(t, imports, "lombok.NoArgsConstructor")
}

func TestLombokGenerator_GenerateTypeAnnotations_WithDirectiveOverride(t *testing.T) {
	cfg := &config.LombokConfig{
		Enabled: true,
		Data:    true,
		Builder: false,
	}
	gen := NewLombokGenerator(cfg)

	// Type with @lombok directive that excludes data and includes builder
	typeDef := &parser.TypeDef{
		Name: "User",
		Directives: []*parser.DirectiveDef{
			{
				Name: "lombok",
				Arguments: map[string]interface{}{
					"exclude": []interface{}{"data"},
					"include": []interface{}{"builder"},
				},
			},
		},
	}
	annotations, imports := gen.GenerateTypeAnnotations(typeDef)

	assert.NotContains(t, annotations, "@Data")
	assert.Contains(t, annotations, "@Builder")
	assert.Contains(t, imports, "lombok.Builder")
}

func TestLombokGenerator_GenerateFieldAnnotations(t *testing.T) {
	cfg := &config.LombokConfig{Enabled: true}
	gen := NewLombokGenerator(cfg)

	field := &parser.FieldDef{Name: "name"}
	annotations, imports := gen.GenerateFieldAnnotations(field)

	// Currently no field-level Lombok annotations
	assert.Empty(t, annotations)
	assert.Empty(t, imports)
}

func TestLombokGenerator_NeedsGettersSetters_Disabled(t *testing.T) {
	cfg := &config.LombokConfig{Enabled: false}
	gen := NewLombokGenerator(cfg)

	assert.True(t, gen.NeedsGettersSetters())
}

func TestLombokGenerator_NeedsGettersSetters_DataEnabled(t *testing.T) {
	cfg := &config.LombokConfig{
		Enabled: true,
		Data:    true,
	}
	gen := NewLombokGenerator(cfg)

	// @Data includes getters/setters
	assert.False(t, gen.NeedsGettersSetters())
}

func TestLombokGenerator_NeedsGettersSetters_GetterSetterEnabled(t *testing.T) {
	cfg := &config.LombokConfig{
		Enabled: true,
		Getter:  true,
		Setter:  true,
	}
	gen := NewLombokGenerator(cfg)

	assert.False(t, gen.NeedsGettersSetters())
}

func TestLombokGenerator_NeedsGettersSetters_OnlyGetter(t *testing.T) {
	cfg := &config.LombokConfig{
		Enabled: true,
		Getter:  true,
		Setter:  false,
	}
	gen := NewLombokGenerator(cfg)

	// Need to generate setters manually
	assert.True(t, gen.NeedsGettersSetters())
}

func TestLombokGenerator_NeedsConstructors_Disabled(t *testing.T) {
	cfg := &config.LombokConfig{Enabled: false}
	gen := NewLombokGenerator(cfg)

	assert.True(t, gen.NeedsConstructors())
}

func TestLombokGenerator_NeedsConstructors_NoArgsEnabled(t *testing.T) {
	cfg := &config.LombokConfig{
		Enabled:           true,
		NoArgsConstructor: true,
	}
	gen := NewLombokGenerator(cfg)

	assert.False(t, gen.NeedsConstructors())
}

func TestLombokGenerator_NeedsConstructors_AllArgsEnabled(t *testing.T) {
	cfg := &config.LombokConfig{
		Enabled:            true,
		AllArgsConstructor: true,
	}
	gen := NewLombokGenerator(cfg)

	assert.False(t, gen.NeedsConstructors())
}

func TestLombokGenerator_NeedsConstructors_NoneEnabled(t *testing.T) {
	cfg := &config.LombokConfig{
		Enabled:            true,
		NoArgsConstructor:  false,
		AllArgsConstructor: false,
	}
	gen := NewLombokGenerator(cfg)

	assert.True(t, gen.NeedsConstructors())
}
