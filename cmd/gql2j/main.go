package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/source-c/go-gql2j/internal/config"
	"github.com/source-c/go-gql2j/internal/generator"
	"github.com/source-c/go-gql2j/internal/output"
	"github.com/source-c/go-gql2j/internal/parser"
)

func main() {
	// Define flags
	configPath := flag.String("config", "", "Path to YAML config file")
	schemaPath := flag.String("schema", "", "GraphQL schema path (overrides config)")
	outputDir := flag.String("output", "", "Output directory (overrides config)")
	packageName := flag.String("package", "", "Java package name (overrides config)")
	javaVersion := flag.Int("java-version", 0, "Target Java version: 8, 11, 17, 21")
	lombok := flag.Bool("lombok", false, "Enable Lombok annotations")
	lombokDisable := flag.Bool("lombok-disable", false, "Disable Lombok annotations")
	validation := flag.Bool("validation", false, "Enable JSR-303 validation")
	validationDisable := flag.Bool("validation-disable", false, "Disable JSR-303 validation")
	validationPkg := flag.String("validation-package", "", "Validation package: jakarta or javax")
	clean := flag.Bool("clean", false, "Clean output directory before generating")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	version := flag.Bool("version", false, "Print version information")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "gql2j - GraphQL to Java code generator\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  gql2j -schema schema.graphql -output ./generated -package com.example.model\n")
		fmt.Fprintf(os.Stderr, "  gql2j -config gql2j.yaml\n")
		fmt.Fprintf(os.Stderr, "  gql2j -config gql2j.yaml -java-version 8 -lombok=false\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	// Handle version flag
	if *version {
		fmt.Println("gql2j version 1.0.1")
		os.Exit(0)
	}

	// Load configuration
	cfg, err := loadConfiguration(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	// Apply flag overrides
	applyOverrides(cfg, *schemaPath, *outputDir, *packageName, *javaVersion,
		*lombok, *lombokDisable, *validation, *validationDisable, *validationPkg)

	// Validate we have required settings
	if cfg.Schema.Path == "" && *schemaPath == "" {
		fmt.Fprintf(os.Stderr, "Error: schema path is required\n")
		fmt.Fprintf(os.Stderr, "Use -schema flag or specify in config file\n")
		os.Exit(1)
	}

	// Use schema from flag if provided
	if *schemaPath != "" {
		cfg.Schema.Path = *schemaPath
	}

	// Resolve paths relative to config file if using config
	if *configPath != "" {
		configDir := filepath.Dir(*configPath)
		if err := cfg.ResolvePaths(configDir); err != nil {
			fmt.Fprintf(os.Stderr, "Error resolving paths: %v\n", err)
			os.Exit(1)
		}
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Configuration error: %v\n", err)
		os.Exit(1)
	}

	// Parse the schema
	if *verbose {
		fmt.Printf("Parsing schema: %s\n", cfg.Schema.Path)
	}

	p := parser.NewParser()
	var schema *parser.Schema

	if len(cfg.Schema.Includes) > 0 {
		schema, err = p.ParseWithIncludes(cfg.Schema.Path, cfg.Schema.Includes)
	} else {
		schema, err = p.ParseFile(cfg.Schema.Path)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing schema: %v\n", err)
		os.Exit(1)
	}

	// Create output directory and optionally clean
	writer := output.NewWriter(cfg.Output.Directory)

	if *clean {
		if *verbose {
			fmt.Printf("Cleaning output directory: %s\n", cfg.Output.Directory)
		}
		if err := writer.Clean(); err != nil {
			fmt.Fprintf(os.Stderr, "Error cleaning output directory: %v\n", err)
			os.Exit(1)
		}
	}

	if err := writer.EnsureDir(); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	// Generate code
	if *verbose {
		fmt.Printf("Generating Java code to: %s\n", cfg.Output.Directory)
		fmt.Printf("Package: %s\n", cfg.Output.Package)
		fmt.Printf("Java version: %d\n", cfg.Java.Version)
		if cfg.Features.Lombok.Enabled {
			fmt.Println("Lombok: enabled")
		}
		if cfg.Features.Validation.Enabled {
			fmt.Printf("Validation: enabled (%s)\n", cfg.Features.Validation.Package)
		}
	}

	gen := generator.NewGenerator(cfg)
	files, err := gen.Generate(schema)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating code: %v\n", err)
		// Continue to write what we can
	}

	// Write files
	result := writer.WriteAllWithResult(files)

	// Report results
	for _, path := range result.Written {
		fmt.Printf("Generated: %s\n", path)
	}

	for _, path := range result.Skipped {
		fmt.Printf("Skipped: %s\n", path)
	}

	for _, err := range result.Errors {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}

	// Print summary
	stats := generator.GetStats(files, result.Errors)
	fmt.Printf("\nGeneration complete: %d classes, %d interfaces, %d enums\n",
		stats.Classes, stats.Interfaces, stats.Enums)

	if stats.ErrorCount > 0 {
		fmt.Printf("%d error(s) occurred\n", stats.ErrorCount)
		os.Exit(1)
	}
}

func loadConfiguration(configPath string) (*config.Config, error) {
	if configPath != "" {
		return config.Load(configPath)
	}

	// Try to find a config file in the current directory
	defaultPaths := []string{"gql2j.yaml", "gql2j.yml", ".gql2j.yaml", ".gql2j.yml"}
	for _, path := range defaultPaths {
		if _, err := os.Stat(path); err == nil {
			return config.Load(path)
		}
	}

	// Return default configuration
	return config.DefaultConfig(), nil
}

func applyOverrides(cfg *config.Config, schemaPath, outputDir, packageName string,
	javaVersion int, lombok, lombokDisable, validation, validationDisable bool, validationPkg string) {

	if schemaPath != "" {
		cfg.Schema.Path = schemaPath
	}
	if outputDir != "" {
		cfg.Output.Directory = outputDir
	}
	if packageName != "" {
		cfg.Output.Package = packageName
	}
	if javaVersion != 0 {
		cfg.Java.Version = javaVersion
	}
	if lombok {
		cfg.Features.Lombok.Enabled = true
	}
	if lombokDisable {
		cfg.Features.Lombok.Enabled = false
	}
	if validation {
		cfg.Features.Validation.Enabled = true
	}
	if validationDisable {
		cfg.Features.Validation.Enabled = false
	}
	if validationPkg != "" {
		cfg.Features.Validation.Package = validationPkg
	}
}
