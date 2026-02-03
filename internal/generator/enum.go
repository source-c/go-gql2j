package generator

import (
	"strings"

	"github.com/source-c/go-gql2j/internal/parser"
)

// EnumGenerator generates Java enums.
type EnumGenerator struct{}

// NewEnumGenerator creates a new enum generator.
func NewEnumGenerator() *EnumGenerator {
	return &EnumGenerator{}
}

// Generate generates a Java enum from a type definition.
func (g *EnumGenerator) Generate(ctx *Context, typeDef *parser.TypeDef) (string, error) {
	tc := NewTypeContext(ctx, typeDef)

	if tc.ShouldSkip() {
		return "", nil
	}

	var sb strings.Builder

	// Generate package declaration
	sb.WriteString("package ")
	sb.WriteString(ctx.Config.Output.Package)
	sb.WriteString(";\n\n")

	// We'll insert imports later
	importPlaceholder := sb.Len()

	// Generate enum annotations
	enumAnnotations := g.generateEnumAnnotations(tc)

	// Generate Javadoc
	if typeDef.Description != "" {
		sb.WriteString(g.generateJavadoc(typeDef.Description))
	}

	// Write annotations
	for _, ann := range enumAnnotations {
		sb.WriteString(ann)
		sb.WriteString("\n")
	}

	// Generate enum declaration
	sb.WriteString("public enum ")
	sb.WriteString(tc.TypeName)
	sb.WriteString(" {\n")

	// Generate enum values
	for i, enumValue := range typeDef.EnumValues {
		// Check if value should be skipped
		if parser.ExtractSkipDirective(enumValue.Directives) != nil {
			continue
		}

		// Generate Javadoc for value
		if enumValue.Description != "" {
			sb.WriteString(g.generateValueJavadoc(enumValue.Description))
		}

		// Generate annotations for value
		valueAnnotations := g.generateValueAnnotations(tc, enumValue)
		for _, ann := range valueAnnotations {
			sb.WriteString("    ")
			sb.WriteString(ann)
			sb.WriteString("\n")
		}

		// Get the Java name for the enum value
		valueName := tc.NamingHelper.GetEnumValueName(enumValue)

		sb.WriteString("    ")
		sb.WriteString(valueName)

		// Add comma or semicolon
		if i < len(typeDef.EnumValues)-1 {
			sb.WriteString(",\n")
		} else {
			sb.WriteString(";\n")
		}
	}

	sb.WriteString("}\n")

	// Build final output with imports
	var result strings.Builder
	result.WriteString(sb.String()[:importPlaceholder])

	imports := tc.Imports.GenerateImportBlock()
	if imports != "" {
		result.WriteString(imports)
		result.WriteString("\n")
	}

	result.WriteString(sb.String()[importPlaceholder:])

	return result.String(), nil
}

func (g *EnumGenerator) generateEnumAnnotations(tc *TypeContext) []string {
	var annotations []string

	// Deprecated annotation
	if deprecated, _ := tc.CustomAnnotation.GenerateDeprecatedAnnotation(tc.TypeDef.Directives); deprecated != "" {
		annotations = append(annotations, deprecated)
	}

	// Custom annotations
	customAnns, customImports := tc.CustomAnnotation.GenerateTypeAnnotations(tc.TypeDef)
	annotations = append(annotations, customAnns...)
	tc.Imports.AddAll(customImports)

	return annotations
}

func (g *EnumGenerator) generateValueAnnotations(tc *TypeContext, enumValue *parser.EnumValueDef) []string {
	var annotations []string

	// Deprecated annotation
	if deprecated, _ := tc.CustomAnnotation.GenerateDeprecatedAnnotation(enumValue.Directives); deprecated != "" {
		annotations = append(annotations, deprecated)
	}

	// Custom annotations
	customAnns, customImports := tc.CustomAnnotation.GenerateEnumValueAnnotations(enumValue)
	annotations = append(annotations, customAnns...)
	tc.Imports.AddAll(customImports)

	return annotations
}

func (g *EnumGenerator) generateJavadoc(description string) string {
	var sb strings.Builder
	sb.WriteString("/**\n")

	lines := strings.Split(description, "\n")
	for _, line := range lines {
		sb.WriteString(" * ")
		sb.WriteString(strings.TrimSpace(line))
		sb.WriteString("\n")
	}

	sb.WriteString(" */\n")
	return sb.String()
}

func (g *EnumGenerator) generateValueJavadoc(description string) string {
	var sb strings.Builder
	sb.WriteString("    /**\n")

	lines := strings.Split(description, "\n")
	for _, line := range lines {
		sb.WriteString("     * ")
		sb.WriteString(strings.TrimSpace(line))
		sb.WriteString("\n")
	}

	sb.WriteString("     */\n")
	return sb.String()
}
