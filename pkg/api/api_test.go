package api

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/source-c/go-gql2j/internal/config"
)

func TestGenerate_WithSchemaString(t *testing.T) {
	opts := Options{
		Schema:  "type User { id: ID! name: String! }",
		Package: "com.test",
	}

	result, err := Generate(opts)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result.Files, 1)
	assert.Equal(t, "User.java", result.Files[0].FileName)
	assert.Contains(t, result.Files[0].Content, "public class User")
}

func TestGenerate_WithSchemaFile(t *testing.T) {
	tmpDir := t.TempDir()
	schemaPath := filepath.Join(tmpDir, "schema.graphql")
	err := os.WriteFile(schemaPath, []byte("type User { id: ID! }"), 0644)
	require.NoError(t, err)

	opts := Options{
		SchemaPath: schemaPath,
		Package:    "com.test",
	}

	result, err := Generate(opts)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result.Files, 1)
}

func TestGenerate_WithIncludePatterns(t *testing.T) {
	tmpDir := t.TempDir()

	mainSchema := filepath.Join(tmpDir, "main.graphql")
	err := os.WriteFile(mainSchema, []byte("type Query { users: [User] }"), 0644)
	require.NoError(t, err)

	typesDir := filepath.Join(tmpDir, "types")
	err = os.Mkdir(typesDir, 0755)
	require.NoError(t, err)

	userSchema := filepath.Join(typesDir, "user.graphql")
	err = os.WriteFile(userSchema, []byte("type User { id: ID! }"), 0644)
	require.NoError(t, err)

	opts := Options{
		SchemaPath:      mainSchema,
		IncludePatterns: []string{filepath.Join(typesDir, "*.graphql")},
		Package:         "com.test",
	}

	result, err := Generate(opts)
	require.NoError(t, err)
	require.NotNil(t, result)
	// Should have Query and User
	assert.GreaterOrEqual(t, len(result.Files), 1)
}

func TestGenerate_NoSchema(t *testing.T) {
	opts := Options{
		Package: "com.test",
	}

	_, err := Generate(opts)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no schema provided")
}

func TestGenerate_WithConfigFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create schema
	schemaPath := filepath.Join(tmpDir, "schema.graphql")
	err := os.WriteFile(schemaPath, []byte("type User { id: ID! }"), 0644)
	require.NoError(t, err)

	// Create config
	configContent := `
schema:
  path: schema.graphql
output:
  package: com.custom
java:
  version: 17
`
	configPath := filepath.Join(tmpDir, "config.yaml")
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	opts := Options{
		SchemaPath: schemaPath,
		ConfigPath: configPath,
	}

	result, err := Generate(opts)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Contains(t, result.Files[0].Content, "package com.custom;")
}

func TestGenerate_WithProgrammaticConfig(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.programmatic"
	cfg.Features.Lombok.Enabled = true
	cfg.Features.Lombok.Data = true

	opts := Options{
		Schema: "type User { id: ID! name: String! }",
		Config: cfg,
	}

	result, err := Generate(opts)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Contains(t, result.Files[0].Content, "package com.programmatic;")
	assert.Contains(t, result.Files[0].Content, "@Data")
}

func TestGenerate_OptionOverrides(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Package = "com.original"

	opts := Options{
		Schema:  "type User { id: ID! }",
		Config:  cfg,
		Package: "com.overridden", // Should override config
	}

	result, err := Generate(opts)
	require.NoError(t, err)
	assert.Contains(t, result.Files[0].Content, "package com.overridden;")
}

func TestGenerate_WithLombok(t *testing.T) {
	opts := Options{
		Schema:       "type User { id: ID! name: String! }",
		Package:      "com.test",
		EnableLombok: true,
	}

	result, err := Generate(opts)
	require.NoError(t, err)
	assert.Contains(t, result.Files[0].Content, "import lombok")
}

func TestGenerate_WithValidation(t *testing.T) {
	opts := Options{
		Schema:           "type User { name: String! }",
		Package:          "com.test",
		EnableValidation: true,
	}

	result, err := Generate(opts)
	require.NoError(t, err)
	assert.Contains(t, result.Files[0].Content, "@NotNull")
}

func TestGenerate_WithJavaVersion(t *testing.T) {
	opts := Options{
		Schema:      "type User { id: ID! }",
		Package:     "com.test",
		JavaVersion: 8,
	}

	result, err := Generate(opts)
	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestGenerate_Stats(t *testing.T) {
	opts := Options{
		Schema: `
type User { id: ID! }
type Post { id: ID! }
interface Node { id: ID! }
enum Status { ACTIVE }
`,
		Package: "com.test",
	}

	result, err := Generate(opts)
	require.NoError(t, err)

	assert.Equal(t, 4, result.Stats.TotalTypes)
	assert.Equal(t, 2, result.Stats.Classes)
	assert.Equal(t, 1, result.Stats.Interfaces)
	assert.Equal(t, 1, result.Stats.Enums)
}

func TestGenerateToDir(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "generated")

	opts := Options{
		Schema:    "type User { id: ID! }",
		Package:   "com.test",
		OutputDir: outputDir,
	}

	result, err := GenerateToDir(opts)
	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify file was written
	_, err = os.Stat(filepath.Join(outputDir, "User.java"))
	assert.NoError(t, err)
}

func TestGenerateToDir_DefaultOutput(t *testing.T) {
	// Use a specific temp directory that we can clean up
	tmpDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tmpDir)

	opts := Options{
		Schema:  "type User { id: ID! }",
		Package: "com.test",
		// OutputDir not specified, should use default ./generated
	}

	_, err := GenerateToDir(opts)
	require.NoError(t, err)

	// Verify file was written to default location
	_, err = os.Stat(filepath.Join(tmpDir, "generated", "User.java"))
	assert.NoError(t, err)
}

func TestParseSchema(t *testing.T) {
	schema, err := ParseSchema("type User { id: ID! }")
	require.NoError(t, err)
	require.NotNil(t, schema)

	userType := schema.GetType("User")
	require.NotNil(t, userType)
}

func TestParseSchema_Invalid(t *testing.T) {
	_, err := ParseSchema("invalid { schema }")
	assert.Error(t, err)
}

func TestParseSchemaFile(t *testing.T) {
	tmpDir := t.TempDir()
	schemaPath := filepath.Join(tmpDir, "schema.graphql")
	err := os.WriteFile(schemaPath, []byte("type User { id: ID! }"), 0644)
	require.NoError(t, err)

	schema, err := ParseSchemaFile(schemaPath)
	require.NoError(t, err)
	require.NotNil(t, schema)
}

func TestParseSchemaFile_NotFound(t *testing.T) {
	_, err := ParseSchemaFile("/nonexistent/schema.graphql")
	assert.Error(t, err)
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	require.NotNil(t, cfg)
	assert.Equal(t, 17, cfg.Java.Version)
}

func TestLoadConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configContent := `
java:
  version: 21
output:
  package: com.test
`
	configPath := filepath.Join(tmpDir, "config.yaml")
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	cfg, err := LoadConfig(configPath)
	require.NoError(t, err)
	assert.Equal(t, 21, cfg.Java.Version)
}

func TestLoadConfig_NotFound(t *testing.T) {
	_, err := LoadConfig("/nonexistent/config.yaml")
	assert.Error(t, err)
}

func TestGenerate_ComplexSchema(t *testing.T) {
	schema := `
interface Node {
  id: ID!
}

type User implements Node {
  id: ID!
  name: String!
  email: String
  posts: [Post!]!
  status: Status!
}

type Post implements Node {
  id: ID!
  title: String!
  author: User!
}

enum Status {
  ACTIVE
  INACTIVE
}

input CreateUserInput {
  name: String!
  email: String
}
`

	opts := Options{
		Schema:  schema,
		Package: "com.test.model",
	}

	result, err := Generate(opts)
	require.NoError(t, err)
	require.NotNil(t, result)

	// Should generate Node (interface), User (class), Post (class), Status (enum), CreateUserInput (input class)
	assert.Len(t, result.Files, 5)

	// Verify all types are generated
	fileNames := make([]string, len(result.Files))
	for i, f := range result.Files {
		fileNames[i] = f.FileName
	}
	assert.Contains(t, fileNames, "Node.java")
	assert.Contains(t, fileNames, "User.java")
	assert.Contains(t, fileNames, "Post.java")
	assert.Contains(t, fileNames, "Status.java")
	assert.Contains(t, fileNames, "CreateUserInput.java")
}

func TestGenerate_WithDirectives(t *testing.T) {
	// Define custom directives in schema
	schema := `
directive @javaName(name: String!) on OBJECT | FIELD_DEFINITION
directive @gql2jSkip on OBJECT | FIELD_DEFINITION

type Query { user: User }
type User @javaName(name: "UserEntity") {
  id: ID!
  name: String!
}

type SkippedType @gql2jSkip {
  field: String
}
`

	opts := Options{
		Schema:  schema,
		Package: "com.test",
	}

	result, err := Generate(opts)
	require.NoError(t, err)

	// Find UserEntity file
	var userEntityFile *GeneratedFile
	for _, f := range result.Files {
		if f.FileName == "UserEntity.java" {
			userEntityFile = f
			break
		}
	}
	require.NotNil(t, userEntityFile, "UserEntity.java should be generated")
	assert.Contains(t, userEntityFile.Content, "public class UserEntity")
}
