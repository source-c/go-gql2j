package generator

import (
	"strings"

	"github.com/source-c/go-gql2j/internal/parser"
)

// FieldGenerator generates Java field declarations.
type FieldGenerator struct{}

// NewFieldGenerator creates a new field generator.
func NewFieldGenerator() *FieldGenerator {
	return &FieldGenerator{}
}

// GenerateField generates a field declaration.
func (g *FieldGenerator) GenerateField(fc *FieldContext) string {
	var sb strings.Builder

	// Add field imports
	fc.TypeContext.Imports.AddAll(fc.Imports)

	// Generate Javadoc if description exists
	if fc.Field.Description != "" {
		sb.WriteString(g.generateJavadoc(fc.Field.Description, "    "))
	}

	// Generate annotations
	annotations := g.generateFieldAnnotations(fc)
	for _, ann := range annotations {
		sb.WriteString("    ")
		sb.WriteString(ann)
		sb.WriteString("\n")
	}

	// Generate field declaration
	visibility := fc.GetVisibility()
	if visibility != "" {
		sb.WriteString("    ")
		sb.WriteString(visibility)
		sb.WriteString(" ")
	} else {
		sb.WriteString("    ")
	}
	sb.WriteString(fc.JavaType)
	sb.WriteString(" ")
	sb.WriteString(fc.FieldName)
	sb.WriteString(";\n")

	return sb.String()
}

func (g *FieldGenerator) generateFieldAnnotations(fc *FieldContext) []string {
	var annotations []string

	// Deprecated annotation
	if deprecated, _ := fc.CustomAnnotation.GenerateDeprecatedAnnotation(fc.Field.Directives); deprecated != "" {
		annotations = append(annotations, deprecated)
	}

	// Validation annotations
	validationAnns, validationImports := fc.ValidationGen.GenerateFieldAnnotations(fc.Field, fc.IsNonNull)
	annotations = append(annotations, validationAnns...)
	fc.TypeContext.Imports.AddAll(validationImports)

	// Custom annotations
	customAnns, customImports := fc.CustomAnnotation.GenerateFieldAnnotations(fc.Field)
	annotations = append(annotations, customAnns...)
	fc.TypeContext.Imports.AddAll(customImports)

	return annotations
}

func (g *FieldGenerator) generateJavadoc(description string, indent string) string {
	var sb strings.Builder
	sb.WriteString(indent)
	sb.WriteString("/**\n")

	// Handle multi-line descriptions
	lines := strings.Split(description, "\n")
	for _, line := range lines {
		sb.WriteString(indent)
		sb.WriteString(" * ")
		sb.WriteString(strings.TrimSpace(line))
		sb.WriteString("\n")
	}

	sb.WriteString(indent)
	sb.WriteString(" */\n")
	return sb.String()
}

// GenerateGetter generates a getter method.
func (g *FieldGenerator) GenerateGetter(fc *FieldContext) string {
	var sb strings.Builder

	methodName := fc.NamingHelper.GetGetterName(fc.FieldName, fc.IsBooleanType())

	sb.WriteString("    public ")
	sb.WriteString(fc.JavaType)
	sb.WriteString(" ")
	sb.WriteString(methodName)
	sb.WriteString("() {\n")
	sb.WriteString("        return this.")
	sb.WriteString(fc.FieldName)
	sb.WriteString(";\n")
	sb.WriteString("    }\n")

	return sb.String()
}

// GenerateSetter generates a setter method.
func (g *FieldGenerator) GenerateSetter(fc *FieldContext) string {
	var sb strings.Builder

	methodName := fc.NamingHelper.GetSetterName(fc.FieldName)

	sb.WriteString("    public void ")
	sb.WriteString(methodName)
	sb.WriteString("(")
	sb.WriteString(fc.JavaType)
	sb.WriteString(" ")
	sb.WriteString(fc.FieldName)
	sb.WriteString(") {\n")
	sb.WriteString("        this.")
	sb.WriteString(fc.FieldName)
	sb.WriteString(" = ")
	sb.WriteString(fc.FieldName)
	sb.WriteString(";\n")
	sb.WriteString("    }\n")

	return sb.String()
}

// GenerateInterfaceMethod generates an interface method declaration.
func (g *FieldGenerator) GenerateInterfaceMethod(tc *TypeContext, field *parser.FieldDef) (string, error) {
	fc, err := NewFieldContext(tc, field)
	if err != nil {
		return "", err
	}

	if fc.ShouldSkip() {
		return "", nil
	}

	var sb strings.Builder

	// Add field imports
	tc.Imports.AddAll(fc.Imports)

	// Generate Javadoc if description exists
	if field.Description != "" {
		sb.WriteString(g.generateJavadoc(field.Description, "    "))
	}

	// Generate annotations
	if deprecated, _ := fc.CustomAnnotation.GenerateDeprecatedAnnotation(field.Directives); deprecated != "" {
		sb.WriteString("    ")
		sb.WriteString(deprecated)
		sb.WriteString("\n")
	}

	// Generate method signature
	methodName := fc.NamingHelper.GetGetterName(fc.FieldName, fc.IsBooleanType())

	sb.WriteString("    ")
	sb.WriteString(fc.JavaType)
	sb.WriteString(" ")
	sb.WriteString(methodName)
	sb.WriteString("();\n")

	return sb.String(), nil
}
