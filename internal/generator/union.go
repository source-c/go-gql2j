package generator

import (
	"strings"

	"github.com/source-c/go-gql2j/internal/parser"
)

// UnionGenerator generates Java marker interfaces for GraphQL unions.
type UnionGenerator struct{}

// NewUnionGenerator creates a new union generator.
func NewUnionGenerator() *UnionGenerator {
	return &UnionGenerator{}
}

// Generate generates a Java marker interface for a union type.
func (g *UnionGenerator) Generate(ctx *Context, typeDef *parser.TypeDef) (string, error) {
	tc := NewTypeContext(ctx, typeDef)

	if tc.ShouldSkip() {
		return "", nil
	}

	var sb strings.Builder

	// Generate package declaration
	sb.WriteString("package ")
	sb.WriteString(ctx.Config.Output.Package)
	sb.WriteString(";\n\n")

	// Generate Javadoc
	if typeDef.Description != "" {
		sb.WriteString(g.generateJavadoc(typeDef.Description))
	} else {
		// Add a default description for unions
		sb.WriteString("/**\n")
		sb.WriteString(" * Union type marker interface.\n")
		sb.WriteString(" * Possible types: ")
		// Note: GraphQL union members aren't directly available in TypeDef
		// but the interface still serves as a marker
		sb.WriteString("\n */\n")
	}

	// Generate marker interface (empty interface)
	sb.WriteString("public interface ")
	sb.WriteString(tc.TypeName)
	sb.WriteString(" {\n")
	sb.WriteString("    // Marker interface for GraphQL union type\n")
	sb.WriteString("}\n")

	return sb.String(), nil
}

func (g *UnionGenerator) generateJavadoc(description string) string {
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
