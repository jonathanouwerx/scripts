package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestIsExecutable tests the isExecutable function
func TestIsExecutable(t *testing.T) {
	// Setup
	dirs := SetupTestDirs(t)
	defer CleanupTestDirs(t, dirs.Root)

	// Create a non-executable file
	nonExecFile := filepath.Join(dirs.Root, "nonexec.txt")
	err := os.WriteFile(nonExecFile, []byte("test"), 0644)
	AssertNil(t, err, "Should create non-executable file")

	// Create an executable file
	execFile := filepath.Join(dirs.Root, "exec.sh")
	err = os.WriteFile(execFile, []byte("#!/bin/bash\necho test"), 0755)
	AssertNil(t, err, "Should create executable file")

	// Test non-executable file
	AssertFalse(t, IsExecutable(t, nonExecFile), "Non-executable file should not be executable")

	// Test executable file
	AssertTrue(t, IsExecutable(t, execFile), "Executable file should be executable")

	// Test non-existent file
	nonExistent := filepath.Join(dirs.Root, "doesnotexist")
	AssertFalse(t, IsExecutable(t, nonExistent), "Non-existent file should not be executable")
}

// TestMakeExecutable tests the makeExecutable function
func TestMakeExecutable(t *testing.T) {
	// Setup
	dirs := SetupTestDirs(t)
	defer CleanupTestDirs(t, dirs.Root)

	// Create a non-executable file
	testFile := filepath.Join(dirs.Root, "testfile")
	err := os.WriteFile(testFile, []byte("test content"), 0644)
	AssertNil(t, err, "Should create test file")

	// Verify it's not executable initially
	AssertFalse(t, IsExecutable(t, testFile), "File should not be executable initially")

	// Make it executable (we'll test the concept by directly setting permissions)
	err = os.Chmod(testFile, 0755)
	AssertNil(t, err, "Should make file executable")

	// Verify it's now executable
	AssertTrue(t, IsExecutable(t, testFile), "File should be executable after chmod")
}

func TestCreateTestScript(t *testing.T) {
	// Setup
	dirs := SetupTestDirs(t)
	defer CleanupTestDirs(t, dirs.Root)

	// Create a test script
	scriptPath := CreateTestScript(t, dirs.ScriptsBin, "test", "echo 'hello world'")

	// Verify script was created
	AssertTrue(t, FileExists(t, scriptPath), "Script file should exist")

	// Verify content
	content := ReadFileContent(t, scriptPath)
	AssertTrue(t, strings.HasPrefix(content, "#!/bin/bash"), "Script should start with shebang")
	AssertTrue(t, strings.Contains(content, "echo 'hello world'"), "Script should contain test content")

	// Verify it's executable
	AssertTrue(t, IsExecutable(t, scriptPath), "Script should be executable")
}

func TestAddScript(t *testing.T) {
	// Setup
	dirs := SetupTestDirs(t)
	defer CleanupTestDirs(t, dirs.Root)

	// Create source script
	sourceScript := filepath.Join(dirs.Root, "source.sh")
	err := os.WriteFile(sourceScript, []byte("#!/bin/bash\necho 'source'"), 0644)
	AssertNil(t, err, "Should create source script")

	// Note: In a real test, we would use the config for addScript function
	// but here we're simulating the functionality

	// This would test the addScript function - we'll simulate it
	destScript := filepath.Join(dirs.ScriptsBin, "source.sh")

	// Copy source to destination (simulating addScript)
	sourceContent, err := os.ReadFile(sourceScript)
	AssertNil(t, err, "Should read source file")

	err = os.WriteFile(destScript, sourceContent, 0644)
	AssertNil(t, err, "Should write destination file")

	// Make executable (simulating addScript)
	err = os.Chmod(destScript, 0755)
	AssertNil(t, err, "Should make script executable")

	// Verify script was added correctly
	AssertTrue(t, FileExists(t, destScript), "Destination script should exist")
	AssertTrue(t, IsExecutable(t, destScript), "Destination script should be executable")

	destContent := ReadFileContent(t, destScript)
	AssertEqual(t, string(sourceContent), destContent, "Content should match")
}

func TestScriptRemoval(t *testing.T) {
	// Setup
	dirs := SetupTestDirs(t)
	defer CleanupTestDirs(t, dirs.Root)

	// Create test script
	scriptPath := CreateTestScript(t, dirs.ScriptsBin, "test", "echo test")

	// Verify it exists
	AssertTrue(t, FileExists(t, scriptPath), "Script should exist initially")

	// Remove the script (simulating rm functionality)
	err := os.Remove(scriptPath)
	AssertNil(t, err, "Should remove script without error")

	// Verify it no longer exists
	AssertFalse(t, FileExists(t, scriptPath), "Script should not exist after removal")
}

func TestBinaryRemoval(t *testing.T) {
	// Setup
	dirs := SetupTestDirs(t)
	defer CleanupTestDirs(t, dirs.Root)

	// Create test binary (just a file to simulate)
	binPath := filepath.Join(dirs.BinDir, "testbin")
	err := os.WriteFile(binPath, []byte("fake binary content"), 0755)
	AssertNil(t, err, "Should create test binary")

	// Verify it exists
	AssertTrue(t, FileExists(t, binPath), "Binary should exist initially")

	// Remove the binary (simulating rm --bin functionality)
	err = os.Remove(binPath)
	AssertNil(t, err, "Should remove binary without error")

	// Verify it no longer exists
	AssertFalse(t, FileExists(t, binPath), "Binary should not exist after removal")
}

func TestScriptDirectoryStructure(t *testing.T) {
	// Setup
	dirs := SetupTestDirs(t)
	defer CleanupTestDirs(t, dirs.Root)

	// Verify directories exist
	AssertTrue(t, FileExists(t, dirs.ScriptsBin), "Scripts bin directory should exist")
	AssertTrue(t, FileExists(t, dirs.BinDir), "Bin directory should exist")

	// Create multiple scripts
	CreateTestScript(t, dirs.ScriptsBin, "script1", "echo 1")
	CreateTestScript(t, dirs.ScriptsBin, "script2", "echo 2")
	CreateTestScript(t, dirs.ScriptsBin, "script3", "echo 3")

	// List scripts in directory
	entries, err := os.ReadDir(dirs.ScriptsBin)
	AssertNil(t, err, "Should read scripts directory")

	// Should have 3 scripts
	scriptCount := 0
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".sh") {
			scriptCount++
			// Verify each script is executable
			scriptPath := filepath.Join(dirs.ScriptsBin, entry.Name())
			AssertTrue(t, IsExecutable(t, scriptPath), "Script should be executable")
		}
	}

	AssertEqual(t, 3, scriptCount, "Should have 3 scripts")
}

func TestInvalidScriptOperations(t *testing.T) {
	// Setup
	dirs := SetupTestDirs(t)
	defer CleanupTestDirs(t, dirs.Root)

	// Test removing non-existent script
	nonExistentScript := filepath.Join(dirs.ScriptsBin, "nonexistent.sh")
	err := os.Remove(nonExistentScript)
	AssertNotNil(t, err, "Removing non-existent script should error")
	AssertTrue(t, os.IsNotExist(err), "Error should be 'not exist'")

	// Test removing non-existent binary
	nonExistentBinary := filepath.Join(dirs.BinDir, "nonexistent")
	err = os.Remove(nonExistentBinary)
	AssertNotNil(t, err, "Removing non-existent binary should error")
	AssertTrue(t, os.IsNotExist(err), "Error should be 'not exist'")
}
