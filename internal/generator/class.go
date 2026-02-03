package generator

import (
	"strings"

	"github.com/source-c/go-gql2j/internal/errors"
	"github.com/source-c/go-gql2j/internal/parser"
)

// ClassGenerator generates Java classes.
type ClassGenerator struct {
	fieldGen *FieldGenerator
}

// NewClassGenerator creates a new class generator.
func NewClassGenerator() *ClassGenerator {
	return &ClassGenerator{
		fieldGen: NewFieldGenerator(),
	}
}

// Generate generates a Java class from a type definition.
func (g *ClassGenerator) Generate(ctx *Context, typeDef *parser.TypeDef) (string, error) {
	tc := NewTypeContext(ctx, typeDef)

	if tc.ShouldSkip() {
		return "", nil
	}

	var sb strings.Builder

	// Generate package declaration
	sb.WriteString("package ")
	sb.WriteString(ctx.Config.Output.Package)
	sb.WriteString(";\n\n")

	// We'll insert imports later after we know what we need
	importPlaceholder := sb.Len()

	// Generate class annotations
	classAnnotations := g.generateClassAnnotations(tc)

	// Generate Javadoc
	if typeDef.Description != "" {
		sb.WriteString(g.generateJavadoc(typeDef.Description))
	}

	// Write annotations
	for _, ann := range classAnnotations {
		sb.WriteString(ann)
		sb.WriteString("\n")
	}

	// Generate class declaration
	sb.WriteString("public class ")
	sb.WriteString(tc.TypeName)

	// Add implements clause for interfaces
	if len(typeDef.Interfaces) > 0 {
		sb.WriteString(" implements ")
		for i, iface := range typeDef.Interfaces {
			if i > 0 {
				sb.WriteString(", ")
			}
			// Apply interface naming convention
			ifaceName := iface
			if ctx.Config.Java.Naming.InterfacePrefix != "" {
				ifaceName = ctx.Config.Java.Naming.InterfacePrefix + iface
			}
			sb.WriteString(ifaceName)
		}
	}

	sb.WriteString(" {\n\n")

	// Generate fields
	var fieldContexts []*FieldContext
	for _, field := range typeDef.Fields {
		fc, err := NewFieldContext(tc, field)
		if err != nil {
			return "", errors.NewGenerateError("failed to create field context", err).
				WithTypeName(typeDef.Name).
				WithFieldName(field.Name)
		}

		if fc.ShouldSkip() {
			continue
		}

		fieldContexts = append(fieldContexts, fc)
		fieldCode := g.fieldGen.GenerateField(fc)
		sb.WriteString(fieldCode)
		sb.WriteString("\n")
	}

	// Generate getters and setters if needed
	if tc.LombokGen.NeedsGettersSetters() && len(fieldContexts) > 0 {
		for _, fc := range fieldContexts {
			sb.WriteString(g.fieldGen.GenerateGetter(fc))
			sb.WriteString("\n")
			sb.WriteString(g.fieldGen.GenerateSetter(fc))
			sb.WriteString("\n")
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

func (g *ClassGenerator) generateClassAnnotations(tc *TypeContext) []string {
	var annotations []string

	// Deprecated annotation
	if deprecated, _ := tc.CustomAnnotation.GenerateDeprecatedAnnotation(tc.TypeDef.Directives); deprecated != "" {
		annotations = append(annotations, deprecated)
	}

	// Lombok annotations
	lombokAnns, lombokImports := tc.LombokGen.GenerateTypeAnnotations(tc.TypeDef)
	annotations = append(annotations, lombokAnns...)
	tc.Imports.AddAll(lombokImports)

	// Validation annotations
	validationAnns, validationImports := tc.ValidationGen.GenerateTypeAnnotations(tc.TypeDef)
	annotations = append(annotations, validationAnns...)
	tc.Imports.AddAll(validationImports)

	// Custom annotations
	customAnns, customImports := tc.CustomAnnotation.GenerateTypeAnnotations(tc.TypeDef)
	annotations = append(annotations, customAnns...)
	tc.Imports.AddAll(customImports)

	return annotations
}

func (g *ClassGenerator) generateJavadoc(description string) string {
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
