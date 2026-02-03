package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocation_String(t *testing.T) {
	tests := []struct {
		name     string
		location *Location
		expected string
	}{
		{
			name:     "nil location",
			location: nil,
			expected: "",
		},
		{
			name:     "file only",
			location: &Location{File: "test.graphql"},
			expected: "test.graphql",
		},
		{
			name:     "file and line",
			location: &Location{File: "test.graphql", Line: 10},
			expected: "test.graphql:10",
		},
		{
			name:     "file, line, and column",
			location: &Location{File: "test.graphql", Line: 10, Column: 5},
			expected: "test.graphql:10:5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.location.String())
		})
	}
}

func TestGeneratorError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *GeneratorError
		expected string
	}{
		{
			name: "basic error",
			err: &GeneratorError{
				Code:    CodeConfig,
				Message: "test error",
			},
			expected: "[CONFIG] test error",
		},
		{
			name: "error with location",
			err: &GeneratorError{
				Code:     CodeParse,
				Message:  "parse failed",
				Location: &Location{File: "schema.graphql", Line: 10},
			},
			expected: "[PARSE] parse failed at schema.graphql:10",
		},
		{
			name: "error with cause",
			err: &GeneratorError{
				Code:    CodeGenerate,
				Message: "generation failed",
				Cause:   errors.New("underlying error"),
			},
			expected: "[GENERATE] generation failed: underlying error",
		},
		{
			name: "error with location and cause",
			err: &GeneratorError{
				Code:     CodeTypemap,
				Message:  "type mapping failed",
				Location: &Location{File: "schema.graphql", Line: 5},
				Cause:    errors.New("unknown type"),
			},
			expected: "[TYPEMAP] type mapping failed at schema.graphql:5: unknown type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.err.Error())
		})
	}
}

func TestGeneratorError_Unwrap(t *testing.T) {
	cause := errors.New("underlying cause")
	err := &GeneratorError{
		Code:    CodeConfig,
		Message: "wrapper error",
		Cause:   cause,
	}

	assert.Equal(t, cause, err.Unwrap())
	assert.True(t, errors.Is(err, cause))
}

func TestGeneratorError_WithContext(t *testing.T) {
	err := &GeneratorError{
		Code:    CodeConfig,
		Message: "test",
	}

	err.WithContext("key1", "value1").WithContext("key2", 42)

	assert.Equal(t, "value1", err.Context["key1"])
	assert.Equal(t, 42, err.Context["key2"])
}

func TestGeneratorError_WithLocation(t *testing.T) {
	err := &GeneratorError{
		Code:    CodeConfig,
		Message: "test",
	}

	loc := &Location{File: "test.graphql", Line: 10}
	err.WithLocation(loc)

	assert.Equal(t, loc, err.Location)
}

func TestNewConfigError(t *testing.T) {
	cause := errors.New("file not found")
	err := NewConfigError("config load failed", cause)

	assert.Equal(t, CodeConfig, err.Code)
	assert.Equal(t, "config load failed", err.Message)
	assert.Equal(t, cause, err.Cause)
}

func TestConfigError_WithField(t *testing.T) {
	err := NewConfigError("invalid value", nil)
	err.WithField("java.version")

	assert.Equal(t, "java.version", err.Field)
	assert.Equal(t, "java.version", err.Context["field"])
}

func TestNewParseError(t *testing.T) {
	err := NewParseError("syntax error", nil)

	assert.Equal(t, CodeParse, err.Code)
	assert.Equal(t, "syntax error", err.Message)
}

func TestParseError_WithTypeName(t *testing.T) {
	err := NewParseError("invalid type", nil)
	err.WithTypeName("User")

	assert.Equal(t, "User", err.TypeName)
	assert.Equal(t, "User", err.Context["type"])
}

func TestNewGenerateError(t *testing.T) {
	err := NewGenerateError("generation failed", nil)

	assert.Equal(t, CodeGenerate, err.Code)
	assert.Equal(t, "generation failed", err.Message)
}

func TestGenerateError_WithTypeAndField(t *testing.T) {
	err := NewGenerateError("field error", nil)
	err.WithTypeName("User").WithFieldName("email")

	assert.Equal(t, "User", err.TypeName)
	assert.Equal(t, "email", err.FieldName)
	assert.Equal(t, "User", err.Context["type"])
	assert.Equal(t, "email", err.Context["field"])
}

func TestNewDirectiveError(t *testing.T) {
	err := NewDirectiveError("invalid directive", nil)

	assert.Equal(t, CodeDirective, err.Code)
}

func TestDirectiveError_WithDetails(t *testing.T) {
	err := NewDirectiveError("directive error", nil)
	err.WithDirective("@javaName").WithTypeName("User").WithFieldName("id")

	assert.Equal(t, "@javaName", err.DirectiveName)
	assert.Equal(t, "User", err.TypeName)
	assert.Equal(t, "id", err.FieldName)
}

func TestNewTypeMappingError(t *testing.T) {
	err := NewTypeMappingError("unknown type", nil)

	assert.Equal(t, CodeTypemap, err.Code)
}

func TestTypeMappingError_WithTypes(t *testing.T) {
	err := NewTypeMappingError("mapping failed", nil)
	err.WithSourceType("CustomScalar").WithTargetType("String")

	assert.Equal(t, "CustomScalar", err.SourceType)
	assert.Equal(t, "String", err.TargetType)
}

func TestNewOutputError(t *testing.T) {
	err := NewOutputError("write failed", nil)

	assert.Equal(t, CodeOutput, err.Code)
}

func TestOutputError_WithFilePath(t *testing.T) {
	err := NewOutputError("write failed", nil)
	err.WithFilePath("/path/to/file.java")

	assert.Equal(t, "/path/to/file.java", err.FilePath)
	assert.Equal(t, "/path/to/file.java", err.Context["file"])
}

func TestErrorCollection_Empty(t *testing.T) {
	ec := NewErrorCollection()

	assert.False(t, ec.HasErrors())
	assert.Empty(t, ec.Errors())
	assert.Equal(t, 0, ec.Count())
	assert.Empty(t, ec.Error())
	assert.Nil(t, ec.ToError())
}

func TestErrorCollection_SingleError(t *testing.T) {
	ec := NewErrorCollection()
	err := errors.New("single error")
	ec.Add(err)

	assert.True(t, ec.HasErrors())
	assert.Len(t, ec.Errors(), 1)
	assert.Equal(t, 1, ec.Count())
	assert.Equal(t, err, ec.ToError())
}

func TestErrorCollection_MultipleErrors(t *testing.T) {
	ec := NewErrorCollection()
	err1 := errors.New("error 1")
	err2 := errors.New("error 2")
	ec.Add(err1)
	ec.Add(err2)

	assert.True(t, ec.HasErrors())
	assert.Len(t, ec.Errors(), 2)
	assert.Equal(t, 2, ec.Count())

	collectionErr := ec.ToError()
	require.NotNil(t, collectionErr)
	assert.Contains(t, collectionErr.Error(), "2 error(s) occurred")
	assert.Contains(t, collectionErr.Error(), "error 1")
	assert.Contains(t, collectionErr.Error(), "error 2")
}

func TestErrorCollection_AddNil(t *testing.T) {
	ec := NewErrorCollection()
	ec.Add(nil)

	assert.False(t, ec.HasErrors())
	assert.Equal(t, 0, ec.Count())
}
