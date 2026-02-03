package annotations

import (
	"github.com/source-c/go-gql2j/internal/parser"
)

// CustomAnnotation represents a custom annotation added via directive.
type CustomAnnotation struct {
	Value   string
	Imports []string
}

// CustomAnnotationGenerator generates custom annotations from directives.
type CustomAnnotationGenerator struct{}

// NewCustomAnnotationGenerator creates a new custom annotation generator.
func NewCustomAnnotationGenerator() *CustomAnnotationGenerator {
	return &CustomAnnotationGenerator{}
}

// GenerateTypeAnnotations extracts custom annotations for a type.
func (g *CustomAnnotationGenerator) GenerateTypeAnnotations(typeDef *parser.TypeDef) ([]string, []string) {
	return g.extractAnnotations(typeDef.Directives)
}

// GenerateFieldAnnotations extracts custom annotations for a field.
func (g *CustomAnnotationGenerator) GenerateFieldAnnotations(field *parser.FieldDef) ([]string, []string) {
	return g.extractAnnotations(field.Directives)
}

// GenerateEnumValueAnnotations extracts custom annotations for an enum value.
func (g *CustomAnnotationGenerator) GenerateEnumValueAnnotations(enumValue *parser.EnumValueDef) ([]string, []string) {
	return g.extractAnnotations(enumValue.Directives)
}

func (g *CustomAnnotationGenerator) extractAnnotations(directives []*parser.DirectiveDef) ([]string, []string) {
	annotationInfos := parser.ExtractAnnotationDirectives(directives)
	if len(annotationInfos) == 0 {
		return nil, nil
	}

	var annotations []string
	var imports []string

	for _, info := range annotationInfos {
		annotations = append(annotations, info.Value)
		imports = append(imports, info.Imports...)
	}

	return annotations, imports
}

// GenerateDeprecatedAnnotation generates @Deprecated annotation if applicable.
func (g *CustomAnnotationGenerator) GenerateDeprecatedAnnotation(directives []*parser.DirectiveDef) (string, string) {
	deprecated := parser.ExtractDeprecatedDirective(directives)
	if deprecated == nil {
		return "", ""
	}

	annotation := "@Deprecated"
	// Note: Java 9+ supports @Deprecated(forRemoval = ..., since = "...")
	// For simplicity, we just use the basic annotation

	return annotation, ""
}
