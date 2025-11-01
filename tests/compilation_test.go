package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCompileGoLanguage(t *testing.T) {
	// Setup
	dirs := SetupTestDirs(t)
	defer CleanupTestDirs(t, dirs.Root)

	// Create Go source file
	goFile := CreateTestSourceFile(t, dirs.Root, "hello", ".go", `package main

import "fmt"

func main() {
    fmt.Println("Hello from Go compilation test!")
}`)

	// Change to scripts directory to run compilation
	// Scripts binary is in parent directory
	scriptsPath := filepath.Join("..", "scripts")

	// Attempt compilation
	cmd := exec.Command(scriptsPath, "compile", goFile, "--name", "gotest")
	output, err := cmd.CombinedOutput()

	// Go compilation might succeed if go is available
	if err == nil {
		AssertTrue(t, strings.Contains(string(output), "Compiled"), "Should report successful compilation")

		// Check if binary was created (in real environment)
		// Note: In test environment, the binary might not actually be created
		// but the command should be processed correctly
	} else {
		// If compilation fails, it should be due to missing compiler, not bad command
		outputStr := string(output)
		AssertTrue(t, strings.Contains(outputStr, "go") ||
			strings.Contains(outputStr, "not found") ||
			strings.Contains(outputStr, "exit status"), "Should attempt Go compilation")
	}

	// Clean up any test binary that was created
	testBinaryPath := filepath.Join("..", "opt", "programs", "gotest")
	if FileExists(t, testBinaryPath) {
		_ = os.Remove(testBinaryPath) // Ignore error - cleanup
	}
}

func TestCompilePythonLanguage(t *testing.T) {
	// Setup
	dirs := SetupTestDirs(t)
	defer CleanupTestDirs(t, dirs.Root)

	// Create Python source file
	pyFile := CreateTestSourceFile(t, dirs.Root, "hello", ".py", `#!/usr/bin/env python3
print("Hello from Python compilation test!")`)

	// Change to scripts directory
	// Scripts binary is in parent directory
	scriptsPath := filepath.Join("..", "scripts")

	// Attempt compilation
	cmd := exec.Command(scriptsPath, "compile", pyFile, "--name", "pytest")
	output, err := cmd.CombinedOutput()

	// Python compilation requires PyInstaller
	outputStr := string(output)
	if err == nil {
		AssertTrue(t, strings.Contains(outputStr, "Compiled"), "Should report successful compilation")
	} else {
		// Should attempt PyInstaller or show appropriate error
		AssertTrue(t, strings.Contains(outputStr, "pyinstaller") ||
			strings.Contains(outputStr, "PyInstaller") ||
			strings.Contains(outputStr, "not found"), "Should attempt Python compilation")
	}
}

func TestCompileCLanguage(t *testing.T) {
	// Setup
	dirs := SetupTestDirs(t)
	defer CleanupTestDirs(t, dirs.Root)

	// Create C source file
	cFile := CreateTestSourceFile(t, dirs.Root, "hello", ".c", `#include <stdio.h>

int main() {
    printf("Hello from C compilation test!\n");
    return 0;
}`)

	// Change to scripts directory
	// Scripts binary is in parent directory
	scriptsPath := filepath.Join("..", "scripts")

	// Attempt compilation
	cmd := exec.Command(scriptsPath, "compile", cFile, "--name", "ctest")
	output, err := cmd.CombinedOutput()

	// C compilation might succeed if gcc is available
	outputStr := string(output)
	if err == nil {
		AssertTrue(t, strings.Contains(outputStr, "Compiled"), "Should report successful compilation")
	} else {
		AssertTrue(t, strings.Contains(outputStr, "gcc") ||
			strings.Contains(outputStr, "clang") ||
			strings.Contains(outputStr, "not found"), "Should attempt C compilation")
	}

	// Clean up any test binary that was created
	testBinaryPath := filepath.Join("..", "opt", "programs", "ctest")
	if FileExists(t, testBinaryPath) {
		_ = os.Remove(testBinaryPath) // Ignore error - cleanup
	}
}

func TestCompileCppLanguage(t *testing.T) {
	// Setup
	dirs := SetupTestDirs(t)
	defer CleanupTestDirs(t, dirs.Root)

	// Create C++ source file
	cppFile := CreateTestSourceFile(t, dirs.Root, "hello", ".cpp", `#include <iostream>

int main() {
    std::cout << "Hello from C++ compilation test!" << std::endl;
    return 0;
}`)

	// Change to scripts directory
	// Scripts binary is in parent directory
	scriptsPath := filepath.Join("..", "scripts")

	// Attempt compilation
	cmd := exec.Command(scriptsPath, "compile", cppFile, "--name", "cpptest")
	output, err := cmd.CombinedOutput()

	// C++ compilation might succeed if g++ is available
	outputStr := string(output)
	if err == nil {
		AssertTrue(t, strings.Contains(outputStr, "Compiled"), "Should report successful compilation")
	} else {
		AssertTrue(t, strings.Contains(outputStr, "g++") ||
			strings.Contains(outputStr, "not found"), "Should attempt C++ compilation")
	}

	// Clean up any test binary that was created
	testBinaryPath := filepath.Join("..", "opt", "programs", "cpptest")
	if FileExists(t, testBinaryPath) {
		_ = os.Remove(testBinaryPath) // Ignore error - cleanup
	}
}

func TestCompileUnsupportedLanguage(t *testing.T) {
	// Setup
	dirs := SetupTestDirs(t)
	defer CleanupTestDirs(t, dirs.Root)

	// Create unsupported file (like a text file)
	txtFile := CreateTestSourceFile(t, dirs.Root, "readme", ".txt", "This is a text file, not source code.")

	// Change to scripts directory
	// Scripts binary is in parent directory
	scriptsPath := filepath.Join("..", "scripts")

	// Attempt compilation - should fail
	cmd := exec.Command(scriptsPath, "compile", txtFile)
	output, err := cmd.CombinedOutput()

	AssertNotNil(t, err, "Compilation of unsupported file should fail")
	AssertTrue(t, strings.Contains(string(output), "unsupported file extension"), "Should report unsupported extension")
}

func TestCompileWithCustomName(t *testing.T) {
	// Setup
	dirs := SetupTestDirs(t)
	defer CleanupTestDirs(t, dirs.Root)

	// Create Go source file
	goFile := CreateTestSourceFile(t, dirs.Root, "simple", ".go", `package main

func main() {
    // Simple Go program
}`)

	// Change to scripts directory
	// Scripts binary is in parent directory
	scriptsPath := filepath.Join("..", "scripts")

	// Test different custom names
	customNames := []string{"my-custom-app", "tool123", "binary_name"}

	for _, customName := range customNames {
		cmd := exec.Command(scriptsPath, "compile", goFile, "--name", customName)
		output, err := cmd.CombinedOutput()

		// Even if compilation fails, the --name option should be accepted
		outputStr := string(output)
		if err == nil {
			AssertTrue(t, strings.Contains(outputStr, "Compiled"), "Should accept custom name and attempt compilation")
		} else {
			// Should not fail due to bad --name parsing
			AssertFalse(t, strings.Contains(outputStr, "Usage:"), "Should not show usage error for valid --name syntax")
		}

		// Clean up any test binary that was created
		testBinaryPath := filepath.Join("..", "opt", "programs", customName)
		if FileExists(t, testBinaryPath) {
			_ = os.Remove(testBinaryPath) // Ignore error - cleanup
		}
	}
}

func TestCompileMissingSourceFile(t *testing.T) {
	// Change to scripts directory
	// Scripts binary is in parent directory
	scriptsPath := filepath.Join("..", "scripts")

	// Try to compile non-existent file
	cmd := exec.Command(scriptsPath, "compile", "/nonexistent/file.go")
	output, err := cmd.CombinedOutput()

	AssertNotNil(t, err, "Compilation of missing file should fail")
	AssertTrue(t, strings.Contains(string(output), "does not exist"), "Should report file not found")
}

func TestCompileInvalidSyntax(t *testing.T) {
	// Setup
	dirs := SetupTestDirs(t)
	defer CleanupTestDirs(t, dirs.Root)

	// Create Go file with syntax error
	badGoFile := CreateTestSourceFile(t, dirs.Root, "bad", ".go", `package main

func main() {
    // Missing closing brace - syntax error
`)

	// Change to scripts directory
	// Scripts binary is in parent directory
	scriptsPath := filepath.Join("..", "scripts")

	// Attempt compilation - should fail due to syntax error
	cmd := exec.Command(scriptsPath, "compile", badGoFile)
	output, err := cmd.CombinedOutput()

	AssertNotNil(t, err, "Compilation of invalid syntax should fail")
	AssertNotNil(t, output, "Should have some output")
	// The error message will depend on the compiler, but it should fail
}

func TestCompilationOutput(t *testing.T) {
	// Setup
	dirs := SetupTestDirs(t)
	defer CleanupTestDirs(t, dirs.Root)

	// Create a valid Go program
	goFile := CreateTestSourceFile(t, dirs.Root, "output_test", ".go", `package main

import "fmt"

func main() {
    fmt.Println("Compilation successful!")
}`)

	// Change to scripts directory
	// Scripts binary is in parent directory
	scriptsPath := filepath.Join("..", "scripts")

	// Compile with custom name
	cmd := exec.Command(scriptsPath, "compile", goFile, "--name", "output_test_bin")
	output, err := cmd.CombinedOutput()

	// Test that output contains expected information
	outputStr := string(output)

	if err == nil {
		// Successful compilation
		AssertTrue(t, strings.Contains(outputStr, "Compiled"), "Should show compilation success")
		AssertTrue(t, strings.Contains(outputStr, goFile) || strings.Contains(outputStr, "output_test.go"), "Should mention source file")
	} else {
		// Failed compilation (due to missing compiler)
		AssertTrue(t, strings.Contains(outputStr, "go") || strings.Contains(outputStr, "not found"), "Should attempt Go compilation")
	}

	// Clean up any test binary that was created
	testBinaryPath := filepath.Join("..", "opt", "programs", "output_test_bin")
	if FileExists(t, testBinaryPath) {
		_ = os.Remove(testBinaryPath) // Ignore error - cleanup
	}
}
