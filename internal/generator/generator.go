package generator

import (
	"github.com/source-c/go-gql2j/internal/config"
	"github.com/source-c/go-gql2j/internal/errors"
	"github.com/source-c/go-gql2j/internal/parser"
)

// GeneratedFile represents a generated Java file.
type GeneratedFile struct {
	FileName string
	Content  string
	TypeDef  *parser.TypeDef
}

// Generator orchestrates Java code generation.
type Generator struct {
	config       *config.Config
	classGen     *ClassGenerator
	interfaceGen *InterfaceGenerator
	enumGen      *EnumGenerator
	unionGen     *UnionGenerator
}

// NewGenerator creates a new generator.
func NewGenerator(cfg *config.Config) *Generator {
	return &Generator{
		config:       cfg,
		classGen:     NewClassGenerator(),
		interfaceGen: NewInterfaceGenerator(),
		enumGen:      NewEnumGenerator(),
		unionGen:     NewUnionGenerator(),
	}
}

// Generate generates Java files for all types in the schema.
func (g *Generator) Generate(schema *parser.Schema) ([]*GeneratedFile, error) {
	ctx := NewContext(g.config, schema)
	errs := errors.NewErrorCollection()

	var files []*GeneratedFile

	for _, typeDef := range schema.Types {
		file, err := g.generateType(ctx, typeDef)
		if err != nil {
			errs.Add(err)
			continue
		}
		if file != nil {
			files = append(files, file)
		}
	}

	if errs.HasErrors() {
		return files, errs.ToError()
	}

	return files, nil
}

// GenerateType generates a Java file for a single type.
func (g *Generator) GenerateType(schema *parser.Schema, typeName string) (*GeneratedFile, error) {
	typeDef := schema.GetType(typeName)
	if typeDef == nil {
		return nil, errors.NewGenerateError("type not found: "+typeName, nil)
	}

	ctx := NewContext(g.config, schema)
	return g.generateType(ctx, typeDef)
}

func (g *Generator) generateType(ctx *Context, typeDef *parser.TypeDef) (*GeneratedFile, error) {
	var content string
	var err error

	switch typeDef.Kind {
	case parser.TypeKindObject, parser.TypeKindInputObject:
		content, err = g.classGen.Generate(ctx, typeDef)

	case parser.TypeKindInterface:
		content, err = g.interfaceGen.Generate(ctx, typeDef)

	case parser.TypeKindEnum:
		content, err = g.enumGen.Generate(ctx, typeDef)

	case parser.TypeKindUnion:
		// Generate unions as marker interfaces
		content, err = g.unionGen.Generate(ctx, typeDef)

	case parser.TypeKindScalar:
		// Scalars are mapped to existing Java types
		return nil, nil

	default:
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	// Empty content means the type was skipped
	if content == "" {
		return nil, nil
	}

	// Determine the file name
	typeName := ctx.NamingHelper.GetTypeName(typeDef)
	fileName := typeName + ".java"

	return &GeneratedFile{
		FileName: fileName,
		Content:  content,
		TypeDef:  typeDef,
	}, nil
}

// Result represents the complete generation result.
type Result struct {
	Files    []*GeneratedFile
	Errors   []error
	Warnings []string
}

// GenerateWithResult generates Java files and returns detailed results.
func (g *Generator) GenerateWithResult(schema *parser.Schema) *Result {
	result := &Result{}
	ctx := NewContext(g.config, schema)

	for _, typeDef := range schema.Types {
		file, err := g.generateType(ctx, typeDef)
		if err != nil {
			result.Errors = append(result.Errors, err)
			continue
		}
		if file != nil {
			result.Files = append(result.Files, file)
		}
	}

	return result
}

// Stats returns generation statistics.
type Stats struct {
	TotalTypes   int
	Classes      int
	Interfaces   int
	Enums        int
	Skipped      int
	ErrorCount   int
}

// GetStats returns statistics about the generated files.
func GetStats(files []*GeneratedFile, errors []error) Stats {
	stats := Stats{
		TotalTypes: len(files),
		ErrorCount: len(errors),
	}

	for _, file := range files {
		switch file.TypeDef.Kind {
		case parser.TypeKindObject, parser.TypeKindInputObject:
			stats.Classes++
		case parser.TypeKindInterface, parser.TypeKindUnion:
			// Unions are generated as marker interfaces
			stats.Interfaces++
		case parser.TypeKindEnum:
			stats.Enums++
		}
	}

	return stats
}
