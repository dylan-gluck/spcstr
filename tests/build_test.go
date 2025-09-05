package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestBuildSystem(t *testing.T) {
	// Get project root directory
	projectRoot, err := getProjectRoot()
	if err != nil {
		t.Fatalf("Failed to find project root: %v", err)
	}

	// Test that Makefile exists
	makefilePath := filepath.Join(projectRoot, "Makefile")
	if _, err := os.Stat(makefilePath); os.IsNotExist(err) {
		t.Error("Makefile should exist in project root")
	}

	// Test that go.mod exists
	goModPath := filepath.Join(projectRoot, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		t.Error("go.mod should exist in project root")
	}

	// Test that main.go exists
	mainGoPath := filepath.Join(projectRoot, "cmd", "spcstr", "main.go")
	if _, err := os.Stat(mainGoPath); os.IsNotExist(err) {
		t.Error("main.go should exist in cmd/spcstr/")
	}
}

func TestCompilation(t *testing.T) {
	projectRoot, err := getProjectRoot()
	if err != nil {
		t.Fatalf("Failed to find project root: %v", err)
	}

	// Change to project root for build
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	err = os.Chdir(projectRoot)
	if err != nil {
		t.Fatalf("Failed to change to project root: %v", err)
	}

	// Test that project compiles without errors
	cmd := exec.Command("go", "build", "./cmd/spcstr")
	if err := cmd.Run(); err != nil {
		t.Errorf("Project should compile without errors: %v", err)
	}

	// Cleanup generated binary
	_ = os.Remove("spcstr")
}

func TestBinarySize(t *testing.T) {
	projectRoot, err := getProjectRoot()
	if err != nil {
		t.Fatalf("Failed to find project root: %v", err)
	}

	// Change to project root for build
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	err = os.Chdir(projectRoot)
	if err != nil {
		t.Fatalf("Failed to change to project root: %v", err)
	}

	// Build binary
	cmd := exec.Command("go", "build", "-o", "test_binary", "./cmd/spcstr")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary for size test: %v", err)
	}
	defer os.Remove("test_binary")

	// Check binary size
	stat, err := os.Stat("test_binary")
	if err != nil {
		t.Fatalf("Failed to stat test binary: %v", err)
	}

	const maxSize = 50 * 1024 * 1024 // 50MB
	if stat.Size() > maxSize {
		t.Errorf("Binary size %d bytes exceeds maximum of %d bytes", stat.Size(), maxSize)
	}
}

func getProjectRoot() (string, error) {
	// Start from current directory and walk up until we find go.mod
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached filesystem root
			break
		}
		dir = parent
	}

	return "", os.ErrNotExist
}
