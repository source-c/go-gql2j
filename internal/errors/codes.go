package errors

// ErrorCode represents the category of an error.
type ErrorCode string

const (
	// CodeConfig indicates a configuration-related error.
	CodeConfig ErrorCode = "CONFIG"

	// CodeParse indicates a GraphQL schema parsing error.
	CodeParse ErrorCode = "PARSE"

	// CodeGenerate indicates a code generation error.
	CodeGenerate ErrorCode = "GENERATE"

	// CodeOutput indicates a file output error.
	CodeOutput ErrorCode = "OUTPUT"

	// CodeDirective indicates a directive processing error.
	CodeDirective ErrorCode = "DIRECTIVE"

	// CodeTypemap indicates a type mapping error.
	CodeTypemap ErrorCode = "TYPEMAP"
)

// String returns the string representation of the error code.
func (c ErrorCode) String() string {
	return string(c)
}
