package output

import (
	"os"
	"path/filepath"

	"github.com/source-c/go-gql2j/internal/errors"
	"github.com/source-c/go-gql2j/internal/generator"
)

// Writer handles writing generated files to disk.
type Writer struct {
	outputDir string
	overwrite bool
}

// NewWriter creates a new file writer.
func NewWriter(outputDir string) *Writer {
	return &Writer{
		outputDir: outputDir,
		overwrite: true,
	}
}

// SetOverwrite sets whether to overwrite existing files.
func (w *Writer) SetOverwrite(overwrite bool) {
	w.overwrite = overwrite
}

// WriteAll writes all generated files to the output directory.
func (w *Writer) WriteAll(files []*generator.GeneratedFile) error {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(w.outputDir, 0755); err != nil {
		return errors.NewOutputError("failed to create output directory", err).
			WithFilePath(w.outputDir)
	}

	errs := errors.NewErrorCollection()

	for _, file := range files {
		if err := w.WriteFile(file); err != nil {
			errs.Add(err)
		}
	}

	return errs.ToError()
}

// WriteFile writes a single generated file to the output directory.
func (w *Writer) WriteFile(file *generator.GeneratedFile) error {
	path := filepath.Join(w.outputDir, file.FileName)

	// Check if file exists and overwrite is disabled
	if !w.overwrite {
		if _, err := os.Stat(path); err == nil {
			return errors.NewOutputError("file already exists and overwrite is disabled", nil).
				WithFilePath(path)
		}
	}

	// Write the file
	if err := os.WriteFile(path, []byte(file.Content), 0644); err != nil {
		return errors.NewOutputError("failed to write file", err).
			WithFilePath(path)
	}

	return nil
}

// GetOutputPath returns the full output path for a file.
func (w *Writer) GetOutputPath(fileName string) string {
	return filepath.Join(w.outputDir, fileName)
}

// Clean removes all Java files from the output directory.
func (w *Writer) Clean() error {
	entries, err := os.ReadDir(w.outputDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return errors.NewOutputError("failed to read output directory", err).
			WithFilePath(w.outputDir)
	}

	errs := errors.NewErrorCollection()

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if filepath.Ext(entry.Name()) == ".java" {
			path := filepath.Join(w.outputDir, entry.Name())
			if err := os.Remove(path); err != nil {
				errs.Add(errors.NewOutputError("failed to remove file", err).
					WithFilePath(path))
			}
		}
	}

	return errs.ToError()
}

// Exists checks if the output directory exists.
func (w *Writer) Exists() bool {
	info, err := os.Stat(w.outputDir)
	return err == nil && info.IsDir()
}

// EnsureDir ensures the output directory exists.
func (w *Writer) EnsureDir() error {
	if err := os.MkdirAll(w.outputDir, 0755); err != nil {
		return errors.NewOutputError("failed to create output directory", err).
			WithFilePath(w.outputDir)
	}
	return nil
}

// WriteResult represents the result of writing files.
type WriteResult struct {
	Written []string
	Skipped []string
	Errors  []error
}

// WriteAllWithResult writes all files and returns detailed results.
func (w *Writer) WriteAllWithResult(files []*generator.GeneratedFile) *WriteResult {
	result := &WriteResult{}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(w.outputDir, 0755); err != nil {
		result.Errors = append(result.Errors,
			errors.NewOutputError("failed to create output directory", err).
				WithFilePath(w.outputDir))
		return result
	}

	for _, file := range files {
		path := filepath.Join(w.outputDir, file.FileName)

		// Check if file exists and overwrite is disabled
		if !w.overwrite {
			if _, err := os.Stat(path); err == nil {
				result.Skipped = append(result.Skipped, path)
				continue
			}
		}

		// Write the file
		if err := os.WriteFile(path, []byte(file.Content), 0644); err != nil {
			result.Errors = append(result.Errors,
				errors.NewOutputError("failed to write file", err).
					WithFilePath(path))
			continue
		}

		result.Written = append(result.Written, path)
	}

	return result
}
