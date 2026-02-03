package errors

import (
	"fmt"
	"strings"
)

// Location represents a position in a source file.
type Location struct {
	File   string
	Line   int
	Column int
}

// String returns a formatted location string.
func (l *Location) String() string {
	if l == nil {
		return ""
	}
	if l.Column > 0 {
		return fmt.Sprintf("%s:%d:%d", l.File, l.Line, l.Column)
	}
	if l.Line > 0 {
		return fmt.Sprintf("%s:%d", l.File, l.Line)
	}
	return l.File
}

// GeneratorError is the base error type for all generator errors.
type GeneratorError struct {
	Code     ErrorCode
	Message  string
	Cause    error
	Context  map[string]interface{}
	Location *Location
}

// Error implements the error interface.
func (e *GeneratorError) Error() string {
	var sb strings.Builder

	sb.WriteString("[")
	sb.WriteString(string(e.Code))
	sb.WriteString("] ")
	sb.WriteString(e.Message)

	if e.Location != nil {
		sb.WriteString(" at ")
		sb.WriteString(e.Location.String())
	}

	if e.Cause != nil {
		sb.WriteString(": ")
		sb.WriteString(e.Cause.Error())
	}

	return sb.String()
}

// Unwrap returns the underlying cause.
func (e *GeneratorError) Unwrap() error {
	return e.Cause
}

// WithContext adds context to the error and returns it.
func (e *GeneratorError) WithContext(key string, value interface{}) *GeneratorError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// WithLocation sets the location on the error and returns it.
func (e *GeneratorError) WithLocation(loc *Location) *GeneratorError {
	e.Location = loc
	return e
}

// ConfigError represents a configuration error.
type ConfigError struct {
	GeneratorError
	Field string
}

// NewConfigError creates a new configuration error.
func NewConfigError(message string, cause error) *ConfigError {
	return &ConfigError{
		GeneratorError: GeneratorError{
			Code:    CodeConfig,
			Message: message,
			Cause:   cause,
		},
	}
}

// WithField sets the field that caused the error.
func (e *ConfigError) WithField(field string) *ConfigError {
	e.Field = field
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context["field"] = field
	return e
}

// ParseError represents a schema parsing error.
type ParseError struct {
	GeneratorError
	TypeName string
}

// NewParseError creates a new parse error.
func NewParseError(message string, cause error) *ParseError {
	return &ParseError{
		GeneratorError: GeneratorError{
			Code:    CodeParse,
			Message: message,
			Cause:   cause,
		},
	}
}

// WithTypeName sets the type name that caused the error.
func (e *ParseError) WithTypeName(typeName string) *ParseError {
	e.TypeName = typeName
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context["type"] = typeName
	return e
}

// GenerateError represents a code generation error.
type GenerateError struct {
	GeneratorError
	TypeName  string
	FieldName string
}

// NewGenerateError creates a new generation error.
func NewGenerateError(message string, cause error) *GenerateError {
	return &GenerateError{
		GeneratorError: GeneratorError{
			Code:    CodeGenerate,
			Message: message,
			Cause:   cause,
		},
	}
}

// WithTypeName sets the type name.
func (e *GenerateError) WithTypeName(typeName string) *GenerateError {
	e.TypeName = typeName
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context["type"] = typeName
	return e
}

// WithFieldName sets the field name.
func (e *GenerateError) WithFieldName(fieldName string) *GenerateError {
	e.FieldName = fieldName
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context["field"] = fieldName
	return e
}

// DirectiveError represents a directive processing error.
type DirectiveError struct {
	GeneratorError
	DirectiveName string
	TypeName      string
	FieldName     string
}

// NewDirectiveError creates a new directive error.
func NewDirectiveError(message string, cause error) *DirectiveError {
	return &DirectiveError{
		GeneratorError: GeneratorError{
			Code:    CodeDirective,
			Message: message,
			Cause:   cause,
		},
	}
}

// WithDirective sets the directive name.
func (e *DirectiveError) WithDirective(name string) *DirectiveError {
	e.DirectiveName = name
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context["directive"] = name
	return e
}

// WithTypeName sets the type name.
func (e *DirectiveError) WithTypeName(typeName string) *DirectiveError {
	e.TypeName = typeName
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context["type"] = typeName
	return e
}

// WithFieldName sets the field name.
func (e *DirectiveError) WithFieldName(fieldName string) *DirectiveError {
	e.FieldName = fieldName
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context["field"] = fieldName
	return e
}

// TypeMappingError represents a type mapping error.
type TypeMappingError struct {
	GeneratorError
	SourceType string
	TargetType string
}

// NewTypeMappingError creates a new type mapping error.
func NewTypeMappingError(message string, cause error) *TypeMappingError {
	return &TypeMappingError{
		GeneratorError: GeneratorError{
			Code:    CodeTypemap,
			Message: message,
			Cause:   cause,
		},
	}
}

// WithSourceType sets the source type.
func (e *TypeMappingError) WithSourceType(sourceType string) *TypeMappingError {
	e.SourceType = sourceType
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context["sourceType"] = sourceType
	return e
}

// WithTargetType sets the target type.
func (e *TypeMappingError) WithTargetType(targetType string) *TypeMappingError {
	e.TargetType = targetType
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context["targetType"] = targetType
	return e
}

// OutputError represents a file output error.
type OutputError struct {
	GeneratorError
	FilePath string
}

// NewOutputError creates a new output error.
func NewOutputError(message string, cause error) *OutputError {
	return &OutputError{
		GeneratorError: GeneratorError{
			Code:    CodeOutput,
			Message: message,
			Cause:   cause,
		},
	}
}

// WithFilePath sets the file path.
func (e *OutputError) WithFilePath(path string) *OutputError {
	e.FilePath = path
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context["file"] = path
	return e
}

// ErrorCollection collects multiple errors for batch reporting.
type ErrorCollection struct {
	errors []error
}

// NewErrorCollection creates a new error collection.
func NewErrorCollection() *ErrorCollection {
	return &ErrorCollection{
		errors: make([]error, 0),
	}
}

// Add adds an error to the collection.
func (c *ErrorCollection) Add(err error) {
	if err != nil {
		c.errors = append(c.errors, err)
	}
}

// HasErrors returns true if there are any errors.
func (c *ErrorCollection) HasErrors() bool {
	return len(c.errors) > 0
}

// Errors returns all collected errors.
func (c *ErrorCollection) Errors() []error {
	return c.errors
}

// Count returns the number of errors.
func (c *ErrorCollection) Count() int {
	return len(c.errors)
}

// Error implements the error interface.
func (c *ErrorCollection) Error() string {
	if len(c.errors) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d error(s) occurred:\n", len(c.errors)))
	for i, err := range c.errors {
		sb.WriteString(fmt.Sprintf("  %d. %s\n", i+1, err.Error()))
	}
	return sb.String()
}

// ToError returns nil if no errors, the single error if one, or the collection if multiple.
func (c *ErrorCollection) ToError() error {
	switch len(c.errors) {
	case 0:
		return nil
	case 1:
		return c.errors[0]
	default:
		return c
	}
}
