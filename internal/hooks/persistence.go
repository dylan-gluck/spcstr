package hooks

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var fileLocks = make(map[string]*sync.Mutex)
var lockMapMutex sync.Mutex

func getFileLock(path string) *sync.Mutex {
	lockMapMutex.Lock()
	defer lockMapMutex.Unlock()

	if lock, exists := fileLocks[path]; exists {
		return lock
	}

	lock := &sync.Mutex{}
	fileLocks[path] = lock
	return lock
}

func AtomicWriteJSON(filePath string, data interface{}) error {
	lock := getFileLock(filePath)
	lock.Lock()
	defer lock.Unlock()

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	tempFile := fmt.Sprintf("%s.tmp.%d", filePath, time.Now().UnixNano())

	if err := os.WriteFile(tempFile, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	if err := os.Rename(tempFile, filePath); err != nil {
		os.Remove(tempFile)
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}

func AtomicReadJSON(filePath string, data interface{}) error {
	lock := getFileLock(filePath)
	lock.Lock()
	defer lock.Unlock()

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	if err := json.Unmarshal(fileData, data); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return nil
}

func AtomicWriteWithRetry(filePath string, data interface{}, maxRetries int) error {
	var lastErr error
	backoff := 10 * time.Millisecond

	for i := 0; i < maxRetries; i++ {
		if err := AtomicWriteJSON(filePath, data); err == nil {
			return nil
		} else {
			lastErr = err
			time.Sleep(backoff)
			backoff *= 2
			if backoff > 500*time.Millisecond {
				backoff = 500 * time.Millisecond
			}
		}
	}

	return fmt.Errorf("failed after %d retries: %w", maxRetries, lastErr)
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func EnsureDirectory(dirPath string) error {
	return os.MkdirAll(dirPath, 0755)
}
