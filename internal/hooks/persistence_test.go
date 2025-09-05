package hooks

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestAtomicWriteJSON(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.json")

	testData := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
		"key3": []string{"a", "b", "c"},
	}

	err := AtomicWriteJSON(filePath, testData)
	if err != nil {
		t.Fatalf("Failed to write JSON: %v", err)
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("File was not created")
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	var loaded map[string]interface{}
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if loaded["key1"] != "value1" {
		t.Errorf("key1 mismatch: got %v, want value1", loaded["key1"])
	}

	if int(loaded["key2"].(float64)) != 42 {
		t.Errorf("key2 mismatch: got %v, want 42", loaded["key2"])
	}
}

func TestAtomicReadJSON(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.json")

	testData := map[string]string{
		"field1": "value1",
		"field2": "value2",
	}

	data, _ := json.MarshalIndent(testData, "", "  ")
	os.WriteFile(filePath, data, 0644)

	var loaded map[string]string
	err := AtomicReadJSON(filePath, &loaded)
	if err != nil {
		t.Fatalf("Failed to read JSON: %v", err)
	}

	if loaded["field1"] != "value1" {
		t.Errorf("field1 mismatch: got %s, want value1", loaded["field1"])
	}

	if loaded["field2"] != "value2" {
		t.Errorf("field2 mismatch: got %s, want value2", loaded["field2"])
	}
}

func TestAtomicWriteWithRetry(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.json")

	testData := map[string]string{"key": "value"}

	err := AtomicWriteWithRetry(filePath, testData, 3)
	if err != nil {
		t.Fatalf("Failed to write with retry: %v", err)
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("File was not created")
	}
}

func TestConcurrentAtomicWrites(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "concurrent.json")

	var wg sync.WaitGroup
	numWriters := 20

	for i := 0; i < numWriters; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			data := map[string]interface{}{
				"writer_id": id,
				"timestamp": time.Now().UnixNano(),
				"data":      "test data",
			}

			err := AtomicWriteJSON(filePath, data)
			if err != nil {
				t.Errorf("Writer %d failed: %v", id, err)
			}
		}(i)
	}

	wg.Wait()

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("File was not created after concurrent writes")
	}

	var finalData map[string]interface{}
	err := AtomicReadJSON(filePath, &finalData)
	if err != nil {
		t.Fatalf("Failed to read final data: %v", err)
	}

	if finalData["data"] != "test data" {
		t.Error("Final data is corrupted")
	}
}

func TestFileExists(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("existing file", func(t *testing.T) {
		filePath := filepath.Join(tmpDir, "exists.txt")
		os.WriteFile(filePath, []byte("content"), 0644)

		if !FileExists(filePath) {
			t.Error("FileExists returned false for existing file")
		}
	})

	t.Run("non-existing file", func(t *testing.T) {
		filePath := filepath.Join(tmpDir, "not_exists.txt")

		if FileExists(filePath) {
			t.Error("FileExists returned true for non-existing file")
		}
	})
}

func TestEnsureDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	dirPath := filepath.Join(tmpDir, "new", "nested", "directory")

	err := EnsureDirectory(dirPath)
	if err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		t.Error("Directory was not created")
	}

	err = EnsureDirectory(dirPath)
	if err != nil {
		t.Error("EnsureDirectory failed on existing directory")
	}
}
