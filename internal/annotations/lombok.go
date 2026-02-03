package annotations

import (
	"github.com/source-c/go-gql2j/internal/config"
	"github.com/source-c/go-gql2j/internal/parser"
)

// LombokAnnotation represents a Lombok annotation.
type LombokAnnotation struct {
	Name    string
	Import  string
	Enabled bool
}

// LombokGenerator generates Lombok annotations.
type LombokGenerator struct {
	config *config.LombokConfig
}

// NewLombokGenerator creates a new Lombok annotation generator.
func NewLombokGenerator(cfg *config.LombokConfig) *LombokGenerator {
	return &LombokGenerator{
		config: cfg,
	}
}

// LombokAnnotations returns all available Lombok annotations.
func LombokAnnotations() map[string]LombokAnnotation {
	return map[string]LombokAnnotation{
		"data": {
			Name:   "@Data",
			Import: "lombok.Data",
		},
		"builder": {
			Name:   "@Builder",
			Import: "lombok.Builder",
		},
		"noArgsConstructor": {
			Name:   "@NoArgsConstructor",
			Import: "lombok.NoArgsConstructor",
		},
		"allArgsConstructor": {
			Name:   "@AllArgsConstructor",
			Import: "lombok.AllArgsConstructor",
		},
		"getter": {
			Name:   "@Getter",
			Import: "lombok.Getter",
		},
		"setter": {
			Name:   "@Setter",
			Import: "lombok.Setter",
		},
		"toString": {
			Name:   "@ToString",
			Import: "lombok.ToString",
		},
		"equalsAndHashCode": {
			Name:   "@EqualsAndHashCode",
			Import: "lombok.EqualsAndHashCode",
		},
		"value": {
			Name:   "@Value",
			Import: "lombok.Value",
		},
		"superBuilder": {
			Name:   "@SuperBuilder",
			Import: "lombok.experimental.SuperBuilder",
		},
	}
}

// GenerateTypeAnnotations generates Lombok annotations for a type.
func (g *LombokGenerator) GenerateTypeAnnotations(typeDef *parser.TypeDef) ([]string, []string) {
	if !g.config.Enabled {
		return nil, nil
	}

	// Check for @lombok directive override
	lombokDirective := parser.ExtractLombokDirective(typeDef.Directives)

	var annotations []string
	var imports []string

	allAnnotations := LombokAnnotations()

	// Determine which annotations to include
	enabled := g.getEnabledAnnotations(lombokDirective)

	for name, include := range enabled {
		if include {
			if ann, ok := allAnnotations[name]; ok {
				annotations = append(annotations, ann.Name)
				imports = append(imports, ann.Import)
			}
		}
	}

	return annotations, imports
}

func (g *LombokGenerator) getEnabledAnnotations(directive *parser.LombokDirectiveInfo) map[string]bool {
	// Start with config defaults
	enabled := map[string]bool{
		"data":               g.config.Data,
		"builder":            g.config.Builder,
		"noArgsConstructor":  g.config.NoArgsConstructor,
		"allArgsConstructor": g.config.AllArgsConstructor,
		"getter":             g.config.Getter,
		"setter":             g.config.Setter,
	}

	// Apply directive overrides
	if directive != nil {
		// Process excludes
		for _, exclude := range directive.Exclude {
			enabled[exclude] = false
		}

		// Process includes
		for _, include := range directive.Include {
			enabled[include] = true
		}
	}

	return enabled
}

// GenerateFieldAnnotations generates Lombok annotations for a field.
func (g *LombokGenerator) GenerateFieldAnnotations(field *parser.FieldDef) ([]string, []string) {
	// Field-level Lombok annotations could include @Getter, @Setter on individual fields
	// Currently not implemented but structure is here for extension
	return nil, nil
}

// NeedsGettersSetters returns true if getters/setters need to be generated manually.
func (g *LombokGenerator) NeedsGettersSetters() bool {
	if !g.config.Enabled {
		return true
	}
	// @Data includes getters and setters
	if g.config.Data {
		return false
	}
	// Check explicit getter/setter
	return !g.config.Getter || !g.config.Setter
}

// NeedsConstructors returns true if constructors need to be generated manually.
func (g *LombokGenerator) NeedsConstructors() bool {
	if !g.config.Enabled {
		return true
	}
	return !g.config.NoArgsConstructor && !g.config.AllArgsConstructor
}
