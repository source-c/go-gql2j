package output

import (
	"strings"
)

// Formatter provides code formatting utilities.
type Formatter struct {
	indentString string
	lineEnding   string
}

// NewFormatter creates a new formatter with default settings.
func NewFormatter() *Formatter {
	return &Formatter{
		indentString: "    ", // 4 spaces
		lineEnding:   "\n",
	}
}

// SetIndent sets the indentation string.
func (f *Formatter) SetIndent(indent string) {
	f.indentString = indent
}

// SetLineEnding sets the line ending.
func (f *Formatter) SetLineEnding(ending string) {
	f.lineEnding = ending
}

// Indent adds indentation to each line of the input.
func (f *Formatter) Indent(s string, level int) string {
	if level <= 0 || s == "" {
		return s
	}

	indent := strings.Repeat(f.indentString, level)
	lines := strings.Split(s, "\n")

	for i, line := range lines {
		if line != "" {
			lines[i] = indent + line
		}
	}

	return strings.Join(lines, f.lineEnding)
}

// FormatBlock formats a block of code with proper indentation.
func (f *Formatter) FormatBlock(header string, body string, indentLevel int) string {
	var sb strings.Builder

	indent := strings.Repeat(f.indentString, indentLevel)

	sb.WriteString(indent)
	sb.WriteString(header)
	sb.WriteString(" {")
	sb.WriteString(f.lineEnding)

	// Indent body
	bodyLines := strings.Split(body, "\n")
	for _, line := range bodyLines {
		if line != "" {
			sb.WriteString(indent)
			sb.WriteString(f.indentString)
			sb.WriteString(line)
		}
		sb.WriteString(f.lineEnding)
	}

	sb.WriteString(indent)
	sb.WriteString("}")
	sb.WriteString(f.lineEnding)

	return sb.String()
}

// JoinLines joins lines with proper line endings.
func (f *Formatter) JoinLines(lines ...string) string {
	return strings.Join(lines, f.lineEnding)
}

// BlankLine returns a blank line.
func (f *Formatter) BlankLine() string {
	return f.lineEnding
}

// WrapInBlock wraps content in a block.
func (f *Formatter) WrapInBlock(content string) string {
	return "{" + f.lineEnding + content + "}" + f.lineEnding
}

// FormatJavadoc formats a Javadoc comment.
func (f *Formatter) FormatJavadoc(description string, indentLevel int) string {
	if description == "" {
		return ""
	}

	indent := strings.Repeat(f.indentString, indentLevel)
	var sb strings.Builder

	sb.WriteString(indent)
	sb.WriteString("/**")
	sb.WriteString(f.lineEnding)

	lines := strings.Split(description, "\n")
	for _, line := range lines {
		sb.WriteString(indent)
		sb.WriteString(" * ")
		sb.WriteString(strings.TrimSpace(line))
		sb.WriteString(f.lineEnding)
	}

	sb.WriteString(indent)
	sb.WriteString(" */")
	sb.WriteString(f.lineEnding)

	return sb.String()
}

// FormatAnnotation formats an annotation with optional parameters.
func (f *Formatter) FormatAnnotation(name string, params map[string]string) string {
	if len(params) == 0 {
		return "@" + name
	}

	var parts []string
	for k, v := range params {
		parts = append(parts, k+" = "+v)
	}

	return "@" + name + "(" + strings.Join(parts, ", ") + ")"
}

// NormalizeLineEndings normalizes line endings to the configured style.
func (f *Formatter) NormalizeLineEndings(s string) string {
	// First normalize all to \n
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")

	// Then convert to target
	if f.lineEnding != "\n" {
		s = strings.ReplaceAll(s, "\n", f.lineEnding)
	}

	return s
}

// RemoveTrailingWhitespace removes trailing whitespace from each line.
func (f *Formatter) RemoveTrailingWhitespace(s string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t")
	}
	return strings.Join(lines, f.lineEnding)
}

// EnsureFinalNewline ensures the string ends with a newline.
func (f *Formatter) EnsureFinalNewline(s string) string {
	if !strings.HasSuffix(s, f.lineEnding) {
		return s + f.lineEnding
	}
	return s
}
