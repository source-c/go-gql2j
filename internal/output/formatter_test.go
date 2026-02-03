package output

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFormatter(t *testing.T) {
	f := NewFormatter()

	require.NotNil(t, f)
	assert.Equal(t, "    ", f.indentString)
	assert.Equal(t, "\n", f.lineEnding)
}

func TestFormatter_SetIndent(t *testing.T) {
	f := NewFormatter()
	f.SetIndent("\t")

	assert.Equal(t, "\t", f.indentString)
}

func TestFormatter_SetLineEnding(t *testing.T) {
	f := NewFormatter()
	f.SetLineEnding("\r\n")

	assert.Equal(t, "\r\n", f.lineEnding)
}

func TestFormatter_Indent_ZeroLevel(t *testing.T) {
	f := NewFormatter()

	result := f.Indent("test", 0)

	assert.Equal(t, "test", result)
}

func TestFormatter_Indent_SingleLevel(t *testing.T) {
	f := NewFormatter()

	result := f.Indent("test", 1)

	assert.Equal(t, "    test", result)
}

func TestFormatter_Indent_MultipleLevel(t *testing.T) {
	f := NewFormatter()

	result := f.Indent("test", 2)

	assert.Equal(t, "        test", result)
}

func TestFormatter_Indent_MultiLine(t *testing.T) {
	f := NewFormatter()

	result := f.Indent("line1\nline2\nline3", 1)

	assert.Equal(t, "    line1\n    line2\n    line3", result)
}

func TestFormatter_Indent_EmptyLines(t *testing.T) {
	f := NewFormatter()

	result := f.Indent("line1\n\nline2", 1)

	// Empty lines should not be indented
	assert.Equal(t, "    line1\n\n    line2", result)
}

func TestFormatter_Indent_EmptyString(t *testing.T) {
	f := NewFormatter()

	result := f.Indent("", 1)

	assert.Equal(t, "", result)
}

func TestFormatter_FormatBlock(t *testing.T) {
	f := NewFormatter()

	result := f.FormatBlock("public class Test", "private int x;", 0)

	expected := "public class Test {\n    private int x;\n}\n"
	assert.Equal(t, expected, result)
}

func TestFormatter_FormatBlock_WithIndent(t *testing.T) {
	f := NewFormatter()

	result := f.FormatBlock("if (true)", "return;", 1)

	expected := "    if (true) {\n        return;\n    }\n"
	assert.Equal(t, expected, result)
}

func TestFormatter_JoinLines(t *testing.T) {
	f := NewFormatter()

	result := f.JoinLines("line1", "line2", "line3")

	assert.Equal(t, "line1\nline2\nline3", result)
}

func TestFormatter_BlankLine(t *testing.T) {
	f := NewFormatter()

	assert.Equal(t, "\n", f.BlankLine())
}

func TestFormatter_BlankLine_CustomEnding(t *testing.T) {
	f := NewFormatter()
	f.SetLineEnding("\r\n")

	assert.Equal(t, "\r\n", f.BlankLine())
}

func TestFormatter_WrapInBlock(t *testing.T) {
	f := NewFormatter()

	result := f.WrapInBlock("content")

	assert.Equal(t, "{\ncontent}\n", result)
}

func TestFormatter_FormatJavadoc_Empty(t *testing.T) {
	f := NewFormatter()

	result := f.FormatJavadoc("", 0)

	assert.Equal(t, "", result)
}

func TestFormatter_FormatJavadoc_SingleLine(t *testing.T) {
	f := NewFormatter()

	result := f.FormatJavadoc("A simple description", 0)

	expected := "/**\n * A simple description\n */\n"
	assert.Equal(t, expected, result)
}

func TestFormatter_FormatJavadoc_MultiLine(t *testing.T) {
	f := NewFormatter()

	result := f.FormatJavadoc("Line one\nLine two", 0)

	expected := "/**\n * Line one\n * Line two\n */\n"
	assert.Equal(t, expected, result)
}

func TestFormatter_FormatJavadoc_WithIndent(t *testing.T) {
	f := NewFormatter()

	result := f.FormatJavadoc("Description", 1)

	expected := "    /**\n     * Description\n     */\n"
	assert.Equal(t, expected, result)
}

func TestFormatter_FormatAnnotation_NoParams(t *testing.T) {
	f := NewFormatter()

	result := f.FormatAnnotation("Override", nil)

	assert.Equal(t, "@Override", result)
}

func TestFormatter_FormatAnnotation_EmptyParams(t *testing.T) {
	f := NewFormatter()

	result := f.FormatAnnotation("Override", map[string]string{})

	assert.Equal(t, "@Override", result)
}

func TestFormatter_FormatAnnotation_WithParams(t *testing.T) {
	f := NewFormatter()

	params := map[string]string{
		"value": "\"test\"",
	}
	result := f.FormatAnnotation("JsonProperty", params)

	assert.Contains(t, result, "@JsonProperty")
	assert.Contains(t, result, "value = \"test\"")
}

func TestFormatter_NormalizeLineEndings_Unix(t *testing.T) {
	f := NewFormatter()

	result := f.NormalizeLineEndings("line1\nline2\nline3")

	assert.Equal(t, "line1\nline2\nline3", result)
}

func TestFormatter_NormalizeLineEndings_Windows(t *testing.T) {
	f := NewFormatter()

	result := f.NormalizeLineEndings("line1\r\nline2\r\nline3")

	assert.Equal(t, "line1\nline2\nline3", result)
}

func TestFormatter_NormalizeLineEndings_OldMac(t *testing.T) {
	f := NewFormatter()

	result := f.NormalizeLineEndings("line1\rline2\rline3")

	assert.Equal(t, "line1\nline2\nline3", result)
}

func TestFormatter_NormalizeLineEndings_Mixed(t *testing.T) {
	f := NewFormatter()

	result := f.NormalizeLineEndings("line1\r\nline2\rline3\nline4")

	assert.Equal(t, "line1\nline2\nline3\nline4", result)
}

func TestFormatter_NormalizeLineEndings_ToWindows(t *testing.T) {
	f := NewFormatter()
	f.SetLineEnding("\r\n")

	result := f.NormalizeLineEndings("line1\nline2")

	assert.Equal(t, "line1\r\nline2", result)
}

func TestFormatter_RemoveTrailingWhitespace(t *testing.T) {
	f := NewFormatter()

	result := f.RemoveTrailingWhitespace("line1   \nline2\t\nline3")

	assert.Equal(t, "line1\nline2\nline3", result)
}

func TestFormatter_EnsureFinalNewline_NoNewline(t *testing.T) {
	f := NewFormatter()

	result := f.EnsureFinalNewline("content")

	assert.Equal(t, "content\n", result)
}

func TestFormatter_EnsureFinalNewline_HasNewline(t *testing.T) {
	f := NewFormatter()

	result := f.EnsureFinalNewline("content\n")

	assert.Equal(t, "content\n", result)
}

func TestFormatter_EnsureFinalNewline_CustomEnding(t *testing.T) {
	f := NewFormatter()
	f.SetLineEnding("\r\n")

	result := f.EnsureFinalNewline("content")

	assert.Equal(t, "content\r\n", result)
}
