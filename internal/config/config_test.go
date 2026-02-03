package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	assert.NotNil(t, cfg)
	assert.Equal(t, "./generated", cfg.Output.Directory)
	assert.Equal(t, "com.example.model", cfg.Output.Package)
	assert.Equal(t, 17, cfg.Java.Version)
	assert.Equal(t, VisibilityPrivate, cfg.Java.FieldVisibility)
	assert.Equal(t, CollectionList, cfg.Java.CollectionType)
	assert.Equal(t, NullableWrapper, cfg.Java.NullableHandling)
	assert.Equal(t, FieldCaseCamel, cfg.Java.Naming.FieldCase)
	assert.False(t, cfg.Features.Lombok.Enabled)
	assert.False(t, cfg.Features.Validation.Enabled)
	assert.Equal(t, ValidationJakarta, cfg.Features.Validation.Package)
}

func TestLoad_ValidConfig(t *testing.T) {
	// Create a temporary config file
	content := `
schema:
  path: "schema.graphql"
output:
  directory: "./gen"
  package: "com.test.model"
java:
  version: 17
  fieldVisibility: "private"
`
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	err := os.WriteFile(configPath, []byte(content), 0644)
	require.NoError(t, err)

	cfg, err := Load(configPath)
	require.NoError(t, err)
	assert.Equal(t, "schema.graphql", cfg.Schema.Path)
	assert.Equal(t, "./gen", cfg.Output.Directory)
	assert.Equal(t, "com.test.model", cfg.Output.Package)
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("/nonexistent/config.yaml")
	assert.Error(t, err)
}

func TestParse_ValidYAML(t *testing.T) {
	data := []byte(`
java:
  version: 21
  fieldVisibility: "protected"
  collectionType: "Set"
features:
  lombok:
    enabled: true
    data: true
  validation:
    enabled: true
    package: "jakarta"
`)

	cfg, err := Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 21, cfg.Java.Version)
	assert.Equal(t, VisibilityProtected, cfg.Java.FieldVisibility)
	assert.Equal(t, CollectionSet, cfg.Java.CollectionType)
	assert.True(t, cfg.Features.Lombok.Enabled)
	assert.True(t, cfg.Features.Validation.Enabled)
}

func TestParse_InvalidYAML(t *testing.T) {
	data := []byte(`invalid: yaml: content: [`)
	_, err := Parse(data)
	assert.Error(t, err)
}

func TestConfig_Validate_ValidConfig(t *testing.T) {
	cfg := DefaultConfig()
	err := cfg.Validate()
	assert.NoError(t, err)
}

func TestConfig_Validate_InvalidJavaVersion(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Java.Version = 5

	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported Java version")
}

func TestConfig_Validate_InvalidVisibility(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Java.FieldVisibility = "invalid"

	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid field visibility")
}

func TestConfig_Validate_InvalidCollectionType(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Java.CollectionType = "Array"

	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid collection type")
}

func TestConfig_Validate_InvalidNullableHandling(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Java.NullableHandling = "null"

	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid nullable handling")
}

func TestConfig_Validate_InvalidValidationPackage(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Features.Validation.Enabled = true
	cfg.Features.Validation.Package = "invalid"

	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid validation package")
}

func TestConfig_Validate_InvalidPackageName(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Output.Package = "123invalid"

	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid Java package name")
}

func TestConfig_Validate_MultipleErrors(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Java.Version = 5
	cfg.Java.FieldVisibility = "invalid"

	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported Java version")
	assert.Contains(t, err.Error(), "invalid field visibility")
}

func TestConfig_Merge(t *testing.T) {
	base := DefaultConfig()
	other := &Config{
		Schema: SchemaConfig{
			Path: "custom.graphql",
		},
		Output: OutputConfig{
			Directory: "./custom",
			Package:   "com.custom",
		},
		Java: JavaConfig{
			Version: 21,
		},
	}

	base.Merge(other)

	assert.Equal(t, "custom.graphql", base.Schema.Path)
	assert.Equal(t, "./custom", base.Output.Directory)
	assert.Equal(t, "com.custom", base.Output.Package)
	assert.Equal(t, 21, base.Java.Version)
}

func TestConfig_Merge_Nil(t *testing.T) {
	base := DefaultConfig()
	original := base.Output.Package

	base.Merge(nil)

	assert.Equal(t, original, base.Output.Package)
}

func TestConfig_ResolvePaths(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Schema.Path = "schema.graphql"
	cfg.Schema.Includes = []string{"types/*.graphql"}
	cfg.Output.Directory = "generated"

	err := cfg.ResolvePaths("/base/path")
	require.NoError(t, err)

	assert.Equal(t, "/base/path/schema.graphql", cfg.Schema.Path)
	assert.Equal(t, "/base/path/types/*.graphql", cfg.Schema.Includes[0])
	assert.Equal(t, "/base/path/generated", cfg.Output.Directory)
}

func TestConfig_ResolvePaths_AbsolutePaths(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Schema.Path = "/absolute/schema.graphql"
	cfg.Output.Directory = "/absolute/generated"

	err := cfg.ResolvePaths("/base/path")
	require.NoError(t, err)

	assert.Equal(t, "/absolute/schema.graphql", cfg.Schema.Path)
	assert.Equal(t, "/absolute/generated", cfg.Output.Directory)
}

func TestSupportedJavaVersions(t *testing.T) {
	versions := SupportedJavaVersions()

	assert.Contains(t, versions, 8)
	assert.Contains(t, versions, 11)
	assert.Contains(t, versions, 17)
	assert.Contains(t, versions, 21)
}

func TestDefaultScalarMappings(t *testing.T) {
	mappings := DefaultScalarMappings()

	assert.Contains(t, mappings, "DateTime")
	assert.Contains(t, mappings, "UUID")
	assert.Contains(t, mappings, "BigDecimal")
	assert.Equal(t, "java.time.LocalDateTime", mappings["DateTime"].JavaType)
}

func TestConfig_JavaVersionOverrides(t *testing.T) {
	data := []byte(`
java:
  version: 8
features:
  validation:
    enabled: true
`)

	cfg, err := Parse(data)
	require.NoError(t, err)

	// Java 8 should use javax instead of jakarta
	assert.Equal(t, ValidationJavax, cfg.Features.Validation.Package)
}
