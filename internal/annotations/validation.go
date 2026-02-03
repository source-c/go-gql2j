package annotations

import (
	"fmt"
	"strings"

	"github.com/source-c/go-gql2j/internal/config"
	"github.com/source-c/go-gql2j/internal/parser"
)

// ValidationAnnotation represents a JSR-303/JSR-380 validation annotation.
type ValidationAnnotation struct {
	Name    string
	Params  string
	Imports []string
}

// ValidationGenerator generates validation annotations.
type ValidationGenerator struct {
	config *config.ValidationConfig
}

// NewValidationGenerator creates a new validation annotation generator.
func NewValidationGenerator(cfg *config.ValidationConfig) *ValidationGenerator {
	return &ValidationGenerator{
		config: cfg,
	}
}

// GetValidationPackage returns the base validation package.
func (g *ValidationGenerator) GetValidationPackage() string {
	switch g.config.Package {
	case config.ValidationJakarta:
		return "jakarta.validation.constraints"
	case config.ValidationJavax:
		return "javax.validation.constraints"
	default:
		return "jakarta.validation.constraints"
	}
}

// GenerateFieldAnnotations generates validation annotations for a field.
func (g *ValidationGenerator) GenerateFieldAnnotations(field *parser.FieldDef, nonNull bool) ([]string, []string) {
	if !g.config.Enabled {
		return nil, nil
	}

	var annotations []string
	var imports []string
	basePkg := g.GetValidationPackage()

	// Handle @NotNull for non-null fields
	if nonNull && g.config.NotNullOnNonNull {
		annotations = append(annotations, "@NotNull")
		imports = append(imports, basePkg+".NotNull")
	}

	// Handle @constraint directive
	constraint := parser.ExtractConstraintDirective(field.Directives)
	if constraint != nil {
		constraintAnns, constraintImports := g.generateConstraintAnnotations(constraint, basePkg)
		annotations = append(annotations, constraintAnns...)
		imports = append(imports, constraintImports...)
	}

	return annotations, imports
}

func (g *ValidationGenerator) generateConstraintAnnotations(c *parser.ConstraintDirectiveInfo, basePkg string) ([]string, []string) {
	var annotations []string
	var imports []string

	// @Size for string/collection length constraints
	if c.MinLength != nil || c.MaxLength != nil {
		var params []string
		if c.MinLength != nil {
			params = append(params, fmt.Sprintf("min = %d", *c.MinLength))
		}
		if c.MaxLength != nil {
			params = append(params, fmt.Sprintf("max = %d", *c.MaxLength))
		}
		annotations = append(annotations, fmt.Sprintf("@Size(%s)", strings.Join(params, ", ")))
		imports = append(imports, basePkg+".Size")
	}

	// @Min and @Max for numeric constraints
	if c.Min != nil {
		annotations = append(annotations, fmt.Sprintf("@Min(%d)", *c.Min))
		imports = append(imports, basePkg+".Min")
	}
	if c.Max != nil {
		annotations = append(annotations, fmt.Sprintf("@Max(%d)", *c.Max))
		imports = append(imports, basePkg+".Max")
	}

	// @Pattern for regex constraints
	if c.Pattern != "" {
		// Escape quotes in pattern
		escapedPattern := strings.ReplaceAll(c.Pattern, "\"", "\\\"")
		annotations = append(annotations, fmt.Sprintf("@Pattern(regexp = \"%s\")", escapedPattern))
		imports = append(imports, basePkg+".Pattern")
	}

	// @NotNull from constraint
	if c.NotNull {
		annotations = append(annotations, "@NotNull")
		imports = append(imports, basePkg+".NotNull")
	}

	// @NotBlank for non-blank strings
	if c.NotBlank {
		annotations = append(annotations, "@NotBlank")
		imports = append(imports, basePkg+".NotBlank")
	}

	// @Email for email validation
	if c.Email {
		annotations = append(annotations, "@Email")
		imports = append(imports, basePkg+".Email")
	}

	return annotations, imports
}

// GenerateTypeAnnotations generates validation annotations for a type (if any).
func (g *ValidationGenerator) GenerateTypeAnnotations(typeDef *parser.TypeDef) ([]string, []string) {
	// Type-level validation annotations could include @Valid for nested objects
	// Currently not implemented but structure is here for extension
	return nil, nil
}
