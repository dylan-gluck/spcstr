package state

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// AtomicWriter provides atomic file write operations using temp file + rename pattern
type AtomicWriter struct {
	timeout time.Duration
}

// NewAtomicWriter creates a new AtomicWriter with specified timeout
func NewAtomicWriter(timeout time.Duration) *AtomicWriter {
	return &AtomicWriter{
		timeout: timeout,
	}
}

// WriteJSON atomically writes JSON data to a file
// Uses temp file + rename pattern to ensure atomic writes
func (a *AtomicWriter) WriteJSON(ctx context.Context, filename string, data interface{}) error {
	// Validate path safety
	if !filepath.IsAbs(filename) {
		return &FileError{
			Op:   "validate_path",
			Path: filename,
			Err:  ErrInvalidPath,
		}
	}

	// Create context with timeout
	writeCtx, cancel := context.WithTimeout(ctx, a.timeout)
	defer cancel()

	// Marshal data to JSON with proper indentation
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return &FileError{
			Op:   "marshal_json",
			Path: filename,
			Err:  err,
		}
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return &FileError{
			Op:   "create_directory",
			Path: dir,
			Err:  err,
		}
	}

	// Generate unique temporary file path to avoid race conditions
	tempPath := filename + ".tmp." + fmt.Sprintf("%d", time.Now().UnixNano())

	// Write to temporary file
	if err := a.writeFileWithContext(writeCtx, tempPath, jsonData); err != nil {
		return &FileError{
			Op:   "write_temp_file",
			Path: tempPath,
			Err:  err,
		}
	}

	// Atomic rename operation
	if err := os.Rename(tempPath, filename); err != nil {
		// Clean up temp file on failure
		os.Remove(tempPath)
		return &FileError{
			Op:   "atomic_rename",
			Path: filename,
			Err:  err,
		}
	}

	return nil
}

// writeFileWithContext writes data to file with context cancellation support
func (a *AtomicWriter) writeFileWithContext(ctx context.Context, filename string, data []byte) error {
	// Create channel for write operation
	done := make(chan error, 1)

	go func() {
		done <- os.WriteFile(filename, data, 0644)
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		// Try to clean up temp file on context cancellation
		os.Remove(filename)
		return ctx.Err()
	}
}

// FileError provides detailed error information for file operations
type FileError struct {
	Op   string
	Path string
	Err  error
}

func (e *FileError) Error() string {
	if e.Path == "" {
		return e.Op + ": " + e.Err.Error()
	}
	return e.Op + " " + e.Path + ": " + e.Err.Error()
}

func (e *FileError) Unwrap() error {
	return e.Err
}

// Custom error types
var (
	ErrInvalidPath = &StateError{Code: "invalid_path", Message: "path must be absolute"}
)

// StateError represents state management specific errors
type StateError struct {
	Code    string
	Message string
}

func (e *StateError) Error() string {
	return "state: " + e.Code + ": " + e.Message
}
