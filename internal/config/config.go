package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/source-c/go-gql2j/internal/errors"
)

// Load loads configuration from a YAML file.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.NewConfigError("failed to read config file", err).
			WithField("path")
	}

	return Parse(data)
}

// Parse parses configuration from YAML data.
func Parse(data []byte) (*Config, error) {
	cfg := DefaultConfig()

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, errors.NewConfigError("failed to parse config YAML", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	// Apply Java version overrides
	cfg.applyVersionOverrides()

	return cfg, nil
}

// Validate validates the configuration.
func (c *Config) Validate() error {
	errs := errors.NewErrorCollection()

	// Validate Java version
	if !isValidJavaVersion(c.Java.Version) {
		errs.Add(errors.NewConfigError(
			fmt.Sprintf("unsupported Java version: %d (supported: 8, 11, 17, 21)", c.Java.Version),
			nil,
		).WithField("java.version"))
	}

	// Validate field visibility
	if !isValidVisibility(c.Java.FieldVisibility) {
		errs.Add(errors.NewConfigError(
			fmt.Sprintf("invalid field visibility: %s (valid: private, protected, package, public)", c.Java.FieldVisibility),
			nil,
		).WithField("java.fieldVisibility"))
	}

	// Validate collection type
	if !isValidCollectionType(c.Java.CollectionType) {
		errs.Add(errors.NewConfigError(
			fmt.Sprintf("invalid collection type: %s (valid: List, Set, Collection)", c.Java.CollectionType),
			nil,
		).WithField("java.collectionType"))
	}

	// Validate nullable handling
	if !isValidNullableHandling(c.Java.NullableHandling) {
		errs.Add(errors.NewConfigError(
			fmt.Sprintf("invalid nullable handling: %s (valid: wrapper, optional, annotation)", c.Java.NullableHandling),
			nil,
		).WithField("java.nullableHandling"))
	}

	// Validate validation package
	if c.Features.Validation.Enabled && !isValidValidationPackage(c.Features.Validation.Package) {
		errs.Add(errors.NewConfigError(
			fmt.Sprintf("invalid validation package: %s (valid: jakarta, javax)", c.Features.Validation.Package),
			nil,
		).WithField("features.validation.package"))
	}

	// Validate output package format
	if c.Output.Package != "" && !isValidJavaPackage(c.Output.Package) {
		errs.Add(errors.NewConfigError(
			fmt.Sprintf("invalid Java package name: %s", c.Output.Package),
			nil,
		).WithField("output.package"))
	}

	return errs.ToError()
}

// applyVersionOverrides applies Java version-specific overrides.
func (c *Config) applyVersionOverrides() {
	if overrides, ok := c.JavaVersionOverrides[c.Java.Version]; ok {
		// Apply validation package override if not explicitly set
		if overrides.Features.Validation.Package != "" {
			c.Features.Validation.Package = overrides.Features.Validation.Package
		}
	}
}

// Merge merges another config into this one, with the other config taking precedence.
func (c *Config) Merge(other *Config) {
	if other == nil {
		return
	}

	// Schema
	if other.Schema.Path != "" {
		c.Schema.Path = other.Schema.Path
	}
	if len(other.Schema.Includes) > 0 {
		c.Schema.Includes = other.Schema.Includes
	}

	// Output
	if other.Output.Directory != "" {
		c.Output.Directory = other.Output.Directory
	}
	if other.Output.Package != "" {
		c.Output.Package = other.Output.Package
	}

	// Java
	if other.Java.Version != 0 {
		c.Java.Version = other.Java.Version
	}
	if other.Java.FieldVisibility != "" {
		c.Java.FieldVisibility = other.Java.FieldVisibility
	}
	if other.Java.CollectionType != "" {
		c.Java.CollectionType = other.Java.CollectionType
	}
	if other.Java.NullableHandling != "" {
		c.Java.NullableHandling = other.Java.NullableHandling
	}

	// Type mappings
	for k, v := range other.TypeMappings.Scalars {
		if c.TypeMappings.Scalars == nil {
			c.TypeMappings.Scalars = make(map[string]ScalarMapping)
		}
		c.TypeMappings.Scalars[k] = v
	}
}

// ResolvePaths resolves relative paths in the configuration.
func (c *Config) ResolvePaths(basePath string) error {
	if c.Schema.Path != "" && !filepath.IsAbs(c.Schema.Path) {
		c.Schema.Path = filepath.Join(basePath, c.Schema.Path)
	}

	for i, include := range c.Schema.Includes {
		if !filepath.IsAbs(include) {
			c.Schema.Includes[i] = filepath.Join(basePath, include)
		}
	}

	if c.Output.Directory != "" && !filepath.IsAbs(c.Output.Directory) {
		c.Output.Directory = filepath.Join(basePath, c.Output.Directory)
	}

	return nil
}

// Validation helpers

func isValidJavaVersion(version int) bool {
	for _, v := range SupportedJavaVersions() {
		if v == version {
			return true
		}
	}
	return false
}

func isValidVisibility(v string) bool {
	switch v {
	case VisibilityPrivate, VisibilityProtected, VisibilityPackage, VisibilityPublic:
		return true
	}
	return false
}

func isValidCollectionType(ct string) bool {
	switch ct {
	case CollectionList, CollectionSet, CollectionCollection:
		return true
	}
	return false
}

func isValidNullableHandling(nh string) bool {
	switch nh {
	case NullableWrapper, NullableOptional, NullableAnnotation:
		return true
	}
	return false
}

func isValidValidationPackage(pkg string) bool {
	switch pkg {
	case ValidationJakarta, ValidationJavax:
		return true
	}
	return false
}

func isValidJavaPackage(pkg string) bool {
	if pkg == "" {
		return false
	}
	// Simple validation: check for valid Java identifier characters
	for i, r := range pkg {
		if r == '.' {
			continue
		}
		if i == 0 && !isJavaIdentifierStart(r) {
			return false
		}
		if !isJavaIdentifierPart(r) {
			return false
		}
	}
	return true
}

func isJavaIdentifierStart(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_' || r == '$'
}

func isJavaIdentifierPart(r rune) bool {
	return isJavaIdentifierStart(r) || (r >= '0' && r <= '9')
}
