package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser_Parse_BasicSchema(t *testing.T) {
	schema := `
type User {
  id: ID!
  name: String!
  email: String
}
`
	p := NewParser()
	result, err := p.Parse(schema, "test.graphql")
	require.NoError(t, err)
	require.NotNil(t, result)

	userType := result.GetType("User")
	require.NotNil(t, userType)
	assert.Equal(t, "User", userType.Name)
	assert.Equal(t, TypeKindObject, userType.Kind)
	assert.Len(t, userType.Fields, 3)

	// Check field types
	idField := findField(userType.Fields, "id")
	require.NotNil(t, idField)
	assert.True(t, idField.Type.NonNull)
	assert.Equal(t, "ID", idField.Type.Name)

	emailField := findField(userType.Fields, "email")
	require.NotNil(t, emailField)
	assert.False(t, emailField.Type.NonNull)
}

func TestParser_Parse_ListTypes(t *testing.T) {
	schema := `
type Post {
  tags: [String!]!
  authors: [User]
}

type User {
  id: ID!
}
`
	p := NewParser()
	result, err := p.Parse(schema, "test.graphql")
	require.NoError(t, err)

	postType := result.GetType("Post")
	require.NotNil(t, postType)

	tagsField := findField(postType.Fields, "tags")
	require.NotNil(t, tagsField)
	assert.True(t, tagsField.Type.IsList())
	assert.True(t, tagsField.Type.NonNull)
	assert.True(t, tagsField.Type.Elem.NonNull)

	authorsField := findField(postType.Fields, "authors")
	require.NotNil(t, authorsField)
	assert.True(t, authorsField.Type.IsList())
	assert.False(t, authorsField.Type.NonNull)
}

func TestParser_Parse_Interface(t *testing.T) {
	schema := `
interface Node {
  id: ID!
}

type User implements Node {
  id: ID!
  name: String!
}
`
	p := NewParser()
	result, err := p.Parse(schema, "test.graphql")
	require.NoError(t, err)

	nodeInterface := result.GetType("Node")
	require.NotNil(t, nodeInterface)
	assert.Equal(t, TypeKindInterface, nodeInterface.Kind)

	userType := result.GetType("User")
	require.NotNil(t, userType)
	assert.Contains(t, userType.Interfaces, "Node")
}

func TestParser_Parse_Enum(t *testing.T) {
	schema := `
enum Status {
  ACTIVE
  INACTIVE
  PENDING
}
`
	p := NewParser()
	result, err := p.Parse(schema, "test.graphql")
	require.NoError(t, err)

	statusEnum := result.GetType("Status")
	require.NotNil(t, statusEnum)
	assert.Equal(t, TypeKindEnum, statusEnum.Kind)
	assert.Len(t, statusEnum.EnumValues, 3)

	values := []string{}
	for _, v := range statusEnum.EnumValues {
		values = append(values, v.Name)
	}
	assert.Contains(t, values, "ACTIVE")
	assert.Contains(t, values, "INACTIVE")
	assert.Contains(t, values, "PENDING")
}

func TestParser_Parse_InputType(t *testing.T) {
	schema := `
input CreateUserInput {
  name: String!
  email: String
}
`
	p := NewParser()
	result, err := p.Parse(schema, "test.graphql")
	require.NoError(t, err)

	inputType := result.GetType("CreateUserInput")
	require.NotNil(t, inputType)
	assert.Equal(t, TypeKindInputObject, inputType.Kind)
	assert.Len(t, inputType.Fields, 2)
}

func TestParser_Parse_WithDirectives(t *testing.T) {
	// Define custom directives in schema to make them valid
	schema := `
directive @javaName(name: String!) on OBJECT | FIELD_DEFINITION
directive @constraint(maxLength: Int) on FIELD_DEFINITION

type Query { user: User }
type User @javaName(name: "UserEntity") {
  id: ID!
  email: String! @constraint(maxLength: 255)
}
`
	p := NewParser()
	result, err := p.Parse(schema, "test.graphql")
	require.NoError(t, err)

	userType := result.GetType("User")
	require.NotNil(t, userType)
	assert.True(t, userType.HasDirective("javaName"))

	javaName := userType.GetDirective("javaName")
	require.NotNil(t, javaName)
	assert.Equal(t, "UserEntity", javaName.GetArgumentString("name"))

	emailField := findField(userType.Fields, "email")
	require.NotNil(t, emailField)
	constraint := emailField.GetDirective("constraint")
	require.NotNil(t, constraint)
	maxLen, ok := constraint.GetArgumentInt("maxLength")
	assert.True(t, ok)
	assert.Equal(t, 255, maxLen)
}

func TestParser_Parse_Descriptions(t *testing.T) {
	schema := `
"""
A user in the system
"""
type User {
  "The unique identifier"
  id: ID!
}
`
	p := NewParser()
	result, err := p.Parse(schema, "test.graphql")
	require.NoError(t, err)

	userType := result.GetType("User")
	require.NotNil(t, userType)
	assert.Contains(t, userType.Description, "A user in the system")

	idField := findField(userType.Fields, "id")
	require.NotNil(t, idField)
	assert.Contains(t, idField.Description, "unique identifier")
}

func TestParser_Parse_InvalidSchema(t *testing.T) {
	schema := `
type User {
  id: InvalidType!
}
`
	p := NewParser()
	_, err := p.Parse(schema, "test.graphql")
	assert.Error(t, err)
}

func TestParser_ParseFile(t *testing.T) {
	tmpDir := t.TempDir()
	schemaPath := filepath.Join(tmpDir, "schema.graphql")
	content := `type User { id: ID! }`
	err := os.WriteFile(schemaPath, []byte(content), 0644)
	require.NoError(t, err)

	p := NewParser()
	result, err := p.ParseFile(schemaPath)
	require.NoError(t, err)
	require.NotNil(t, result)

	userType := result.GetType("User")
	require.NotNil(t, userType)
}

func TestParser_ParseFile_NotFound(t *testing.T) {
	p := NewParser()
	_, err := p.ParseFile("/nonexistent/schema.graphql")
	assert.Error(t, err)
}

func TestParser_ParseFiles_Multiple(t *testing.T) {
	tmpDir := t.TempDir()

	schema1 := filepath.Join(tmpDir, "types.graphql")
	err := os.WriteFile(schema1, []byte(`type User { id: ID! }`), 0644)
	require.NoError(t, err)

	schema2 := filepath.Join(tmpDir, "inputs.graphql")
	err = os.WriteFile(schema2, []byte(`input CreateUserInput { name: String! }`), 0644)
	require.NoError(t, err)

	p := NewParser()
	result, err := p.ParseFiles([]string{schema1, schema2})
	require.NoError(t, err)

	assert.NotNil(t, result.GetType("User"))
	assert.NotNil(t, result.GetType("CreateUserInput"))
}

func TestParser_ParseWithIncludes(t *testing.T) {
	tmpDir := t.TempDir()

	mainSchema := filepath.Join(tmpDir, "main.graphql")
	err := os.WriteFile(mainSchema, []byte(`type Query { users: [User] }`), 0644)
	require.NoError(t, err)

	typesDir := filepath.Join(tmpDir, "types")
	err = os.Mkdir(typesDir, 0755)
	require.NoError(t, err)

	userSchema := filepath.Join(typesDir, "user.graphql")
	err = os.WriteFile(userSchema, []byte(`type User { id: ID! name: String! }`), 0644)
	require.NoError(t, err)

	p := NewParser()
	result, err := p.ParseWithIncludes(mainSchema, []string{filepath.Join(typesDir, "*.graphql")})
	require.NoError(t, err)

	assert.NotNil(t, result.GetType("User"))
	assert.NotNil(t, result.GetType("Query"))
}

func TestSchema_TypesByKind(t *testing.T) {
	schema := `
type User { id: ID! }
type Post { id: ID! }
interface Node { id: ID! }
enum Status { ACTIVE }
input CreateUserInput { name: String! }
`
	p := NewParser()
	result, err := p.Parse(schema, "test.graphql")
	require.NoError(t, err)

	objects := result.ObjectTypes()
	assert.Len(t, objects, 2)

	interfaces := result.InterfaceTypes()
	assert.Len(t, interfaces, 1)

	enums := result.EnumTypes()
	assert.Len(t, enums, 1)

	inputs := result.InputTypes()
	assert.Len(t, inputs, 1)
}

func TestTypeRef_NamedType(t *testing.T) {
	// Simple type
	simple := &TypeRef{Name: "String", NonNull: true}
	assert.Equal(t, "String", simple.NamedType())
	assert.True(t, simple.IsNamed())
	assert.False(t, simple.IsList())

	// List type
	list := &TypeRef{
		Elem:    &TypeRef{Name: "String", NonNull: true},
		NonNull: true,
	}
	assert.Equal(t, "String", list.NamedType())
	assert.False(t, list.IsNamed())
	assert.True(t, list.IsList())
}

func TestDirectiveDef_GetArguments(t *testing.T) {
	directive := &DirectiveDef{
		Name: "constraint",
		Arguments: map[string]interface{}{
			"min":     int64(1),
			"max":     int64(100),
			"pattern": "^[a-z]+$",
			"email":   true,
		},
	}

	// String argument
	assert.Equal(t, "^[a-z]+$", directive.GetArgumentString("pattern"))
	assert.Equal(t, "", directive.GetArgumentString("nonexistent"))

	// Int argument
	min, ok := directive.GetArgumentInt("min")
	assert.True(t, ok)
	assert.Equal(t, 1, min)

	_, ok = directive.GetArgumentInt("nonexistent")
	assert.False(t, ok)

	// Bool argument
	email, ok := directive.GetArgumentBool("email")
	assert.True(t, ok)
	assert.True(t, email)

	_, ok = directive.GetArgumentBool("nonexistent")
	assert.False(t, ok)
}

func TestDirectiveDef_GetArgumentStringSlice(t *testing.T) {
	directive := &DirectiveDef{
		Name: "test",
		Arguments: map[string]interface{}{
			"values": []interface{}{"a", "b", "c"},
		},
	}

	values := directive.GetArgumentStringSlice("values")
	assert.Equal(t, []string{"a", "b", "c"}, values)

	assert.Nil(t, directive.GetArgumentStringSlice("nonexistent"))
}

// Helper function
func findField(fields []*FieldDef, name string) *FieldDef {
	for _, f := range fields {
		if f.Name == name {
			return f
		}
	}
	return nil
}
