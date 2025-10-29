package tests

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestConfig represents the configuration structure
type TestConfig struct {
	ScriptDir string `json:"scriptDir"`
	BinDir    string `json:"binDir"`
}

// TestDirs holds temporary directories for testing
type TestDirs struct {
	Root       string
	ScriptsBin string
	BinDir     string
	ConfigFile string
}

// SetupTestDirs creates temporary directories for testing
func SetupTestDirs(t *testing.T) *TestDirs {
	t.Helper()

	// Create temporary root directory
	rootDir, err := os.MkdirTemp("", "scripts_test_")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Create subdirectories
	scriptsBin := filepath.Join(rootDir, "scripts_bin")
	binDir := filepath.Join(rootDir, "bin")
	configFile := filepath.Join(rootDir, ".config.json")

	dirs := []string{scriptsBin, binDir}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			CleanupTestDirs(t, rootDir)
			t.Fatalf("Failed to create dir %s: %v", dir, err)
		}
	}

	return &TestDirs{
		Root:       rootDir,
		ScriptsBin: scriptsBin,
		BinDir:     binDir,
		ConfigFile: configFile,
	}
}

// CleanupTestDirs removes temporary test directories
func CleanupTestDirs(t *testing.T, rootDir string) {
	t.Helper()
	if err := os.RemoveAll(rootDir); err != nil {
		t.Logf("Warning: failed to cleanup %s: %v", rootDir, err)
	}
}

// CreateTestConfig creates a test configuration file
func CreateTestConfig(t *testing.T, configFile string, scriptDir, binDir string) {
	t.Helper()

	// Create directory structure if it doesn't exist
	configDir := filepath.Dir(configFile)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	config := TestConfig{
		ScriptDir: scriptDir,
		BinDir:    binDir,
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}

	if err := os.WriteFile(configFile, data, 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}
}

// CreateTestScript creates a test script file
func CreateTestScript(t *testing.T, dir, name, content string) string {
	t.Helper()

	scriptPath := filepath.Join(dir, name+".sh")
	fullContent := "#!/bin/bash\n" + content

	if err := os.WriteFile(scriptPath, []byte(fullContent), 0755); err != nil {
		t.Fatalf("Failed to create test script: %v", err)
	}

	return scriptPath
}

// CreateTestSourceFile creates a test source file for compilation
func CreateTestSourceFile(t *testing.T, dir, name, ext, content string) string {
	t.Helper()

	filePath := filepath.Join(dir, name+"."+ext)

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test source file: %v", err)
	}

	return filePath
}

// FileExists checks if a file exists
func FileExists(t *testing.T, path string) bool {
	t.Helper()
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// IsExecutable checks if a file is executable
func IsExecutable(t *testing.T, path string) bool {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode()&0100 != 0
}

// ReadFileContent reads the entire content of a file
func ReadFileContent(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", path, err)
	}
	return string(data)
}

// MockExecCommand mocks exec.Command for testing
func MockExecCommand(command string, args ...string) *MockCmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return &MockCmd{cmd}
}

// MockCmd wraps exec.Cmd for mocking
type MockCmd struct {
	*exec.Cmd
}

// CombinedOutput returns mock output
func (m *MockCmd) CombinedOutput() ([]byte, error) {
	return m.Cmd.CombinedOutput()
}

// Run runs the mock command
func (m *MockCmd) Run() error {
	return m.Cmd.Run()
}

// TestHelperProcess is used for mocking exec commands
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	// Get the command that was supposed to be run
	args := os.Args
	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}
		args = args[1:]
	}

	if len(args) == 0 {
		os.Exit(1)
	}

	cmd, args := args[0], args[1:]

	// Mock different commands
	switch cmd {
	case "go":
		if len(args) >= 2 && args[0] == "build" && args[1] == "-o" {
			// Simulate successful go build
			os.Exit(0)
		}
	case "gcc", "g++":
		// Simulate successful C/C++ compilation
		os.Exit(0)
	case "pyinstaller":
		// Simulate successful Python compilation
		os.Exit(0)
	default:
		// Unknown command, fail
		os.Exit(1)
	}
}

// AssertEqual checks if two values are equal
func AssertEqual(t *testing.T, expected, actual interface{}, message string) {
	t.Helper()
	if expected != actual {
		t.Errorf("%s: expected %v, got %v", message, expected, actual)
	}
}

// AssertTrue checks if a condition is true
func AssertTrue(t *testing.T, condition bool, message string) {
	t.Helper()
	if !condition {
		t.Errorf("%s", message)
	}
}

// AssertFalse checks if a condition is false
func AssertFalse(t *testing.T, condition bool, message string) {
	t.Helper()
	if condition {
		t.Errorf("%s", message)
	}
}

// AssertNotEqual checks if two values are not equal
func AssertNotEqual(t *testing.T, expected, actual interface{}, message string) {
	t.Helper()
	if expected == actual {
		t.Errorf("%s: expected %v, got %v (should not be equal)", message, expected, actual)
	}
}

// AssertNil checks if a value is nil
func AssertNil(t *testing.T, value interface{}, message string) {
	t.Helper()
	if value != nil {
		t.Errorf("%s: expected nil, got %v", message, value)
	}
}

// AssertNotNil checks if a value is not nil
func AssertNotNil(t *testing.T, value interface{}, message string) {
	t.Helper()
	if value == nil {
		t.Errorf("%s: expected not nil", message)
	}
}
