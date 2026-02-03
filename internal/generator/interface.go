package generator

import (
	"strings"

	"github.com/source-c/go-gql2j/internal/errors"
	"github.com/source-c/go-gql2j/internal/parser"
)

// InterfaceGenerator generates Java interfaces.
type InterfaceGenerator struct {
	fieldGen *FieldGenerator
}

// NewInterfaceGenerator creates a new interface generator.
func NewInterfaceGenerator() *InterfaceGenerator {
	return &InterfaceGenerator{
		fieldGen: NewFieldGenerator(),
	}
}

// Generate generates a Java interface from a type definition.
func (g *InterfaceGenerator) Generate(ctx *Context, typeDef *parser.TypeDef) (string, error) {
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

	// Generate interface annotations
	interfaceAnnotations := g.generateInterfaceAnnotations(tc)

	// Generate Javadoc
	if typeDef.Description != "" {
		sb.WriteString(g.generateJavadoc(typeDef.Description))
	}

	// Write annotations
	for _, ann := range interfaceAnnotations {
		sb.WriteString(ann)
		sb.WriteString("\n")
	}

	// Generate interface declaration
	sb.WriteString("public interface ")
	sb.WriteString(tc.TypeName)

	// Add extends clause for parent interfaces
	if len(typeDef.Interfaces) > 0 {
		sb.WriteString(" extends ")
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

	// Generate method declarations
	for _, field := range typeDef.Fields {
		methodCode, err := g.fieldGen.GenerateInterfaceMethod(tc, field)
		if err != nil {
			return "", errors.NewGenerateError(
				"failed to generate interface method",
				err,
			).WithTypeName(typeDef.Name).WithFieldName(field.Name)
		}
		if methodCode != "" {
			sb.WriteString(methodCode)
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

func (g *InterfaceGenerator) generateInterfaceAnnotations(tc *TypeContext) []string {
	var annotations []string

	// Deprecated annotation
	if deprecated, _ := tc.CustomAnnotation.GenerateDeprecatedAnnotation(tc.TypeDef.Directives); deprecated != "" {
		annotations = append(annotations, deprecated)
	}

	// Custom annotations (no Lombok for interfaces typically)
	customAnns, customImports := tc.CustomAnnotation.GenerateTypeAnnotations(tc.TypeDef)
	annotations = append(annotations, customAnns...)
	tc.Imports.AddAll(customImports)

	return annotations
}

func (g *InterfaceGenerator) generateJavadoc(description string) string {
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
