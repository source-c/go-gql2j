package output

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/source-c/go-gql2j/internal/generator"
)

func TestNewWriter(t *testing.T) {
	w := NewWriter("/tmp/test")

	require.NotNil(t, w)
	assert.Equal(t, "/tmp/test", w.outputDir)
	assert.True(t, w.overwrite)
}

func TestWriter_SetOverwrite(t *testing.T) {
	w := NewWriter("/tmp/test")
	w.SetOverwrite(false)

	assert.False(t, w.overwrite)
}

func TestWriter_WriteFile(t *testing.T) {
	tmpDir := t.TempDir()
	w := NewWriter(tmpDir)

	file := &generator.GeneratedFile{
		FileName: "Test.java",
		Content:  "public class Test {}",
	}

	err := w.WriteFile(file)
	require.NoError(t, err)

	// Verify file was written
	content, err := os.ReadFile(filepath.Join(tmpDir, "Test.java"))
	require.NoError(t, err)
	assert.Equal(t, "public class Test {}", string(content))
}

func TestWriter_WriteFile_OverwriteDisabled(t *testing.T) {
	tmpDir := t.TempDir()
	w := NewWriter(tmpDir)
	w.SetOverwrite(false)

	// Create existing file
	existingPath := filepath.Join(tmpDir, "Test.java")
	err := os.WriteFile(existingPath, []byte("original content"), 0644)
	require.NoError(t, err)

	file := &generator.GeneratedFile{
		FileName: "Test.java",
		Content:  "new content",
	}

	err = w.WriteFile(file)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "file already exists")

	// Verify original content is preserved
	content, err := os.ReadFile(existingPath)
	require.NoError(t, err)
	assert.Equal(t, "original content", string(content))
}

func TestWriter_WriteFile_OverwriteEnabled(t *testing.T) {
	tmpDir := t.TempDir()
	w := NewWriter(tmpDir)
	w.SetOverwrite(true)

	// Create existing file
	existingPath := filepath.Join(tmpDir, "Test.java")
	err := os.WriteFile(existingPath, []byte("original content"), 0644)
	require.NoError(t, err)

	file := &generator.GeneratedFile{
		FileName: "Test.java",
		Content:  "new content",
	}

	err = w.WriteFile(file)
	require.NoError(t, err)

	// Verify content was overwritten
	content, err := os.ReadFile(existingPath)
	require.NoError(t, err)
	assert.Equal(t, "new content", string(content))
}

func TestWriter_WriteAll(t *testing.T) {
	tmpDir := t.TempDir()
	w := NewWriter(tmpDir)

	files := []*generator.GeneratedFile{
		{FileName: "User.java", Content: "public class User {}"},
		{FileName: "Post.java", Content: "public class Post {}"},
	}

	err := w.WriteAll(files)
	require.NoError(t, err)

	// Verify both files were written
	_, err = os.Stat(filepath.Join(tmpDir, "User.java"))
	assert.NoError(t, err)

	_, err = os.Stat(filepath.Join(tmpDir, "Post.java"))
	assert.NoError(t, err)
}

func TestWriter_WriteAll_CreatesDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "subdir", "generated")
	w := NewWriter(outputDir)

	files := []*generator.GeneratedFile{
		{FileName: "Test.java", Content: "public class Test {}"},
	}

	err := w.WriteAll(files)
	require.NoError(t, err)

	// Verify directory was created and file exists
	_, err = os.Stat(filepath.Join(outputDir, "Test.java"))
	assert.NoError(t, err)
}

func TestWriter_WriteAllWithResult(t *testing.T) {
	tmpDir := t.TempDir()
	w := NewWriter(tmpDir)

	files := []*generator.GeneratedFile{
		{FileName: "User.java", Content: "public class User {}"},
		{FileName: "Post.java", Content: "public class Post {}"},
	}

	result := w.WriteAllWithResult(files)

	assert.Len(t, result.Written, 2)
	assert.Empty(t, result.Skipped)
	assert.Empty(t, result.Errors)
}

func TestWriter_WriteAllWithResult_WithSkipped(t *testing.T) {
	tmpDir := t.TempDir()
	w := NewWriter(tmpDir)
	w.SetOverwrite(false)

	// Create existing file
	err := os.WriteFile(filepath.Join(tmpDir, "Existing.java"), []byte("original"), 0644)
	require.NoError(t, err)

	files := []*generator.GeneratedFile{
		{FileName: "New.java", Content: "public class New {}"},
		{FileName: "Existing.java", Content: "public class Existing {}"},
	}

	result := w.WriteAllWithResult(files)

	assert.Len(t, result.Written, 1)
	assert.Len(t, result.Skipped, 1)
	assert.Empty(t, result.Errors)
}

func TestWriter_GetOutputPath(t *testing.T) {
	w := NewWriter("/output/dir")

	path := w.GetOutputPath("Test.java")

	assert.Equal(t, "/output/dir/Test.java", path)
}

func TestWriter_Clean(t *testing.T) {
	tmpDir := t.TempDir()
	w := NewWriter(tmpDir)

	// Create some Java files
	err := os.WriteFile(filepath.Join(tmpDir, "User.java"), []byte("class User {}"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(tmpDir, "Post.java"), []byte("class Post {}"), 0644)
	require.NoError(t, err)
	// Create a non-Java file
	err = os.WriteFile(filepath.Join(tmpDir, "readme.txt"), []byte("readme"), 0644)
	require.NoError(t, err)

	err = w.Clean()
	require.NoError(t, err)

	// Java files should be removed
	_, err = os.Stat(filepath.Join(tmpDir, "User.java"))
	assert.True(t, os.IsNotExist(err))

	_, err = os.Stat(filepath.Join(tmpDir, "Post.java"))
	assert.True(t, os.IsNotExist(err))

	// Non-Java file should remain
	_, err = os.Stat(filepath.Join(tmpDir, "readme.txt"))
	assert.NoError(t, err)
}

func TestWriter_Clean_NonexistentDirectory(t *testing.T) {
	w := NewWriter("/nonexistent/directory")

	err := w.Clean()
	assert.NoError(t, err)
}

func TestWriter_Exists(t *testing.T) {
	tmpDir := t.TempDir()

	w := NewWriter(tmpDir)
	assert.True(t, w.Exists())

	w2 := NewWriter("/nonexistent/directory")
	assert.False(t, w2.Exists())
}

func TestWriter_EnsureDir(t *testing.T) {
	tmpDir := t.TempDir()
	newDir := filepath.Join(tmpDir, "new", "nested", "dir")
	w := NewWriter(newDir)

	assert.False(t, w.Exists())

	err := w.EnsureDir()
	require.NoError(t, err)

	assert.True(t, w.Exists())
}
