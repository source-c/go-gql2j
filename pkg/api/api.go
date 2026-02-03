// Package api provides a public API for the gql2j library.
package api

import (
	"github.com/source-c/go-gql2j/internal/config"
	"github.com/source-c/go-gql2j/internal/errors"
	"github.com/source-c/go-gql2j/internal/generator"
	"github.com/source-c/go-gql2j/internal/output"
	"github.com/source-c/go-gql2j/internal/parser"
)

// Options configures the code generation.
type Options struct {
	// Schema is the GraphQL schema content.
	Schema string

	// SchemaPath is the path to the GraphQL schema file.
	SchemaPath string

	// IncludePatterns are glob patterns for additional schema files.
	IncludePatterns []string

	// OutputDir is the directory for generated files.
	OutputDir string

	// Package is the Java package name.
	Package string

	// ConfigPath is the path to a YAML config file.
	ConfigPath string

	// Config is a programmatic configuration (takes precedence over ConfigPath).
	Config *config.Config

	// JavaVersion is the target Java version (8, 11, 17, 21).
	JavaVersion int

	// EnableLombok enables Lombok annotations.
	EnableLombok bool

	// EnableValidation enables JSR-303 validation annotations.
	EnableValidation bool

	// ValidationPackage is "jakarta" or "javax".
	ValidationPackage string
}

// Result contains the generation results.
type Result struct {
	// Files are the generated Java files.
	Files []*GeneratedFile

	// Errors are any errors that occurred during generation.
	Errors []error

	// Stats contains generation statistics.
	Stats Stats
}

// GeneratedFile represents a generated Java file.
type GeneratedFile struct {
	FileName string
	Content  string
}

// Stats contains generation statistics.
type Stats struct {
	TotalTypes int
	Classes    int
	Interfaces int
	Enums      int
	Errors     int
}

// Generate generates Java code from a GraphQL schema.
func Generate(opts Options) (*Result, error) {
	// Load or create configuration
	cfg, err := loadConfig(opts)
	if err != nil {
		return nil, err
	}

	// Parse the schema
	p := parser.NewParser()
	var schema *parser.Schema

	if opts.Schema != "" {
		schema, err = p.Parse(opts.Schema, "schema.graphql")
	} else if opts.SchemaPath != "" {
		if len(opts.IncludePatterns) > 0 {
			schema, err = p.ParseWithIncludes(opts.SchemaPath, opts.IncludePatterns)
		} else {
			schema, err = p.ParseFile(opts.SchemaPath)
		}
	} else if cfg.Schema.Path != "" {
		schema, err = p.ParseWithIncludes(cfg.Schema.Path, cfg.Schema.Includes)
	} else {
		return nil, errors.NewConfigError("no schema provided", nil)
	}

	if err != nil {
		return nil, err
	}

	// Generate code
	gen := generator.NewGenerator(cfg)
	genResult := gen.GenerateWithResult(schema)

	// Convert to public types
	result := &Result{
		Errors: genResult.Errors,
	}

	for _, f := range genResult.Files {
		result.Files = append(result.Files, &GeneratedFile{
			FileName: f.FileName,
			Content:  f.Content,
		})
	}

	// Calculate stats
	stats := generator.GetStats(genResult.Files, genResult.Errors)
	result.Stats = Stats{
		TotalTypes: stats.TotalTypes,
		Classes:    stats.Classes,
		Interfaces: stats.Interfaces,
		Enums:      stats.Enums,
		Errors:     stats.ErrorCount,
	}

	return result, nil
}

// GenerateToDir generates Java code and writes it to a directory.
func GenerateToDir(opts Options) (*Result, error) {
	result, err := Generate(opts)
	if err != nil {
		return nil, err
	}

	// Determine output directory
	outputDir := opts.OutputDir
	if outputDir == "" {
		cfg, _ := loadConfig(opts)
		if cfg != nil {
			outputDir = cfg.Output.Directory
		}
	}
	if outputDir == "" {
		outputDir = "./generated"
	}

	// Write files
	w := output.NewWriter(outputDir)
	internalFiles := make([]*generator.GeneratedFile, len(result.Files))
	for i, f := range result.Files {
		internalFiles[i] = &generator.GeneratedFile{
			FileName: f.FileName,
			Content:  f.Content,
		}
	}

	if err := w.WriteAll(internalFiles); err != nil {
		result.Errors = append(result.Errors, err)
	}

	return result, nil
}

func loadConfig(opts Options) (*config.Config, error) {
	var cfg *config.Config

	// Use provided config
	if opts.Config != nil {
		cfg = opts.Config
	} else if opts.ConfigPath != "" {
		// Load from file
		var err error
		cfg, err = config.Load(opts.ConfigPath)
		if err != nil {
			return nil, err
		}
	} else {
		// Use defaults
		cfg = config.DefaultConfig()
	}

	// Apply option overrides
	if opts.OutputDir != "" {
		cfg.Output.Directory = opts.OutputDir
	}
	if opts.Package != "" {
		cfg.Output.Package = opts.Package
	}
	if opts.JavaVersion != 0 {
		cfg.Java.Version = opts.JavaVersion
	}
	if opts.EnableLombok {
		cfg.Features.Lombok.Enabled = true
	}
	if opts.EnableValidation {
		cfg.Features.Validation.Enabled = true
	}
	if opts.ValidationPackage != "" {
		cfg.Features.Validation.Package = opts.ValidationPackage
	}

	return cfg, nil
}

// ParseSchema parses a GraphQL schema from a string.
func ParseSchema(content string) (*parser.Schema, error) {
	p := parser.NewParser()
	return p.Parse(content, "schema.graphql")
}

// ParseSchemaFile parses a GraphQL schema from a file.
func ParseSchemaFile(path string) (*parser.Schema, error) {
	p := parser.NewParser()
	return p.ParseFile(path)
}

// DefaultConfig returns the default configuration.
func DefaultConfig() *config.Config {
	return config.DefaultConfig()
}

// LoadConfig loads configuration from a YAML file.
func LoadConfig(path string) (*config.Config, error) {
	return config.Load(path)
}
