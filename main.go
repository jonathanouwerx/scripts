package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Config struct {
	ScriptDir string `json:"scriptDir"`
	BinDir    string `json:"binDir"`
}

func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	mode := info.Mode()
	return mode&0100 != 0
}

func makeExecutable(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	mode := info.Mode()
	newMode := mode | 0100 // Add execute permission for owner
	return os.Chmod(path, newMode)
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return strings.Replace(path, "~", homeDir, 1)
	}
	return path
}

func loadConfig() (*Config, error) {
	// Try to find the config file in the correct location
	var scriptsDir string

	// First, try to get the actual executable path
	if execPath, err := os.Executable(); err == nil {
		execDir := filepath.Dir(execPath)
		// Check if this looks like a scripts installation directory
		// (contains the scripts binary and possibly scripts_bin)
		if info, err := os.Stat(filepath.Join(execDir, "scripts_bin")); err == nil && info.IsDir() {
			scriptsDir = execDir
		} else if info, err := os.Stat(filepath.Join(execDir, "scripts")); err == nil && info.Mode()&0100 != 0 {
			// Check if there's a scripts binary in this directory
			scriptsDir = execDir
		}
	}

	// If we couldn't find the scripts directory from the executable,
	// check if we're running from the source directory
	if scriptsDir == "" {
		if cwd, err := os.Getwd(); err == nil {
			if info, err := os.Stat(filepath.Join(cwd, "scripts_bin")); err == nil && info.IsDir() {
				scriptsDir = cwd
			}
		}
	}

	// As a last resort, use user config directory
	if scriptsDir == "" {
		if homeDir, err := os.UserHomeDir(); err == nil {
			scriptsDir = filepath.Join(homeDir, ".config", "scripts")
		} else {
			return nil, fmt.Errorf("could not determine config directory")
		}
	}

	configPath := filepath.Join(scriptsDir, ".config.json")

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config
		defaultConfig := &Config{
			ScriptDir: expandPath("~/code/personal/scripts/scripts_bin"),
			BinDir:    expandPath("~/opt/programs"),
		}
		if err := saveConfig(defaultConfig); err != nil {
			return nil, fmt.Errorf("failed to create default config: %v", err)
		}
		return defaultConfig, nil
	}

	// Read existing config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	return &config, nil
}

func saveConfig(config *Config) error {
	// Use the same logic as loadConfig to find the scripts directory
	var scriptsDir string

	// First, try to get the actual executable path
	if execPath, err := os.Executable(); err == nil {
		execDir := filepath.Dir(execPath)
		// Check if this looks like a scripts installation directory
		// (contains the scripts binary and possibly scripts_bin)
		if info, err := os.Stat(filepath.Join(execDir, "scripts_bin")); err == nil && info.IsDir() {
			scriptsDir = execDir
		} else if info, err := os.Stat(filepath.Join(execDir, "scripts")); err == nil && info.Mode()&0100 != 0 {
			// Check if there's a scripts binary in this directory
			scriptsDir = execDir
		}
	}

	// If we couldn't find the scripts directory from the executable,
	// check if we're running from the source directory
	if scriptsDir == "" {
		if cwd, err := os.Getwd(); err == nil {
			if info, err := os.Stat(filepath.Join(cwd, "scripts_bin")); err == nil && info.IsDir() {
				scriptsDir = cwd
			}
		}
	}

	// As a last resort, use user config directory
	if scriptsDir == "" {
		if homeDir, err := os.UserHomeDir(); err == nil {
			scriptsDir = filepath.Join(homeDir, ".config", "scripts")
		} else {
			return fmt.Errorf("could not determine config directory")
		}
	}

	configPath := filepath.Join(scriptsDir, ".config.json")

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

func readyScripts(paths []string) error {
	for _, path := range paths {
		// If path is a directory, find all .sh files in it
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			files, err := filepath.Glob(filepath.Join(path, "*.sh"))
			if err != nil {
				return fmt.Errorf("failed to glob %s: %v", path, err)
			}
			for _, file := range files {
				if !isExecutable(file) {
					fmt.Printf("Making %s executable\n", filepath.Base(file))
					if err := makeExecutable(file); err != nil {
						return fmt.Errorf("failed to make %s executable: %v", file, err)
					}
				} else {
					fmt.Printf("%s is already executable\n", filepath.Base(file))
				}
			}
		} else {
			// Handle single file
			if !strings.HasSuffix(path, ".sh") {
				path = path + ".sh"
			}
			if !isExecutable(path) {
				fmt.Printf("Making %s executable\n", filepath.Base(path))
				if err := makeExecutable(path); err != nil {
					return fmt.Errorf("failed to make %s executable: %v", path, err)
				}
			} else {
				fmt.Printf("%s is already executable\n", filepath.Base(path))
			}
		}
	}
	return nil
}

func addScript(scriptPath string, config *Config) error {
	// Check if source script exists
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return fmt.Errorf("script %s does not exist", scriptPath)
	}

	// Ensure it's a .sh file
	if !strings.HasSuffix(scriptPath, ".sh") {
		return fmt.Errorf("script must have .sh extension")
	}

	// Get the script name without extension
	scriptName := strings.TrimSuffix(filepath.Base(scriptPath), ".sh")
	destPath := filepath.Join(config.ScriptDir, scriptName+".sh")

	// Create scripts_bin directory if it doesn't exist
	if err := os.MkdirAll(config.ScriptDir, 0755); err != nil {
		return fmt.Errorf("failed to create scripts directory: %v", err)
	}

	// Copy the script
	sourceData, err := os.ReadFile(scriptPath)
	if err != nil {
		return fmt.Errorf("failed to read source script: %v", err)
	}

	if err := os.WriteFile(destPath, sourceData, 0644); err != nil {
		return fmt.Errorf("failed to write script to scripts_bin: %v", err)
	}

	// Make it executable
	if err := makeExecutable(destPath); err != nil {
		return fmt.Errorf("failed to make script executable: %v", err)
	}

	fmt.Printf("Added %s to scripts_bin\n", scriptName+".sh")
	return nil
}

func compileSource(sourcePath, binaryName string, config *Config) error {
	// Check if source file exists
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		return fmt.Errorf("source file %s does not exist", sourcePath)
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(config.BinDir, 0755); err != nil {
		return fmt.Errorf("failed to create bin directory: %v", err)
	}

	// Get file extension to determine language
	ext := strings.ToLower(filepath.Ext(sourcePath))

	// Use provided binary name or default to source file name
	name := binaryName
	if name == "" {
		name = strings.TrimSuffix(filepath.Base(sourcePath), filepath.Ext(sourcePath))
	}
	outputPath := filepath.Join(config.BinDir, name)

	var err error
	switch ext {
	case ".go":
		err = compileGo(sourcePath, outputPath)
	case ".py":
		err = compilePython(sourcePath, outputPath)
	case ".v":
		err = compileV(sourcePath, outputPath)
	case ".rs":
		err = compileRust(sourcePath, outputPath)
	case ".c":
		err = compileC(sourcePath, outputPath)
	case ".cpp", ".cc", ".cxx":
		err = compileCpp(sourcePath, outputPath)
	default:
		return fmt.Errorf("unsupported file extension: %s", ext)
	}

	if err != nil {
		return err
	}

	// Make binary executable
	if err := makeExecutable(outputPath); err != nil {
		return fmt.Errorf("failed to make binary executable: %v", err)
	}

	fmt.Printf("Compiled %s to %s\n", sourcePath, outputPath)
	return nil
}

func compileGo(sourcePath, outputPath string) error {
	cmd := exec.Command("go", "build", "-o", outputPath, sourcePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func compilePython(sourcePath, outputPath string) error {
	// Use PyInstaller to create standalone executable
	cmd := exec.Command("pyinstaller", "--onefile", "--distpath", filepath.Dir(outputPath), "--name", filepath.Base(outputPath), sourcePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("PyInstaller compilation failed: %v (make sure PyInstaller is installed)", err)
	}

	// PyInstaller creates files in dist directory, move to final location
	distPath := filepath.Join(filepath.Dir(outputPath), filepath.Base(outputPath))
	if _, err := os.Stat(distPath); err == nil {
		return os.Rename(distPath, outputPath)
	}
	return nil
}

func compileV(sourcePath, outputPath string) error {
	cmd := exec.Command("v", "-prod", "-o", outputPath, sourcePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func compileRust(sourcePath, outputPath string) error {
	// Check if this is a Cargo project
	dir := filepath.Dir(sourcePath)
	if _, err := os.Stat(filepath.Join(dir, "Cargo.toml")); err == nil {
		// Cargo project
		cmd := exec.Command("cargo", "build", "--release")
		cmd.Dir = dir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
		// Copy binary from target/release/ to output path
		binaryName := strings.TrimSuffix(filepath.Base(sourcePath), ".rs")
		srcPath := filepath.Join(dir, "target", "release", binaryName)
		return exec.Command("cp", srcPath, outputPath).Run()
	} else {
		// Single file compilation with rustc
		cmd := exec.Command("rustc", "-o", outputPath, sourcePath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
}

func compileC(sourcePath, outputPath string) error {
	cmd := exec.Command("gcc", "-o", outputPath, sourcePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func compileCpp(sourcePath, outputPath string) error {
	cmd := exec.Command("g++", "-o", outputPath, sourcePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func printHelp() {
	fmt.Println("scripts - A tool for managing and running shell scripts and compiling binaries")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  scripts <script_name> [args...]    Run a script from scripts_bin/")
	fmt.Println("  scripts list                        List available scripts and binaries")
	fmt.Println("  scripts ready <script_name> [-a]    Make scripts in scripts_bin executable")
	fmt.Println("  scripts add <script.sh>             Add script to scripts_bin/")
	fmt.Println("  scripts compile <source> [--name <binary>]    Compile source to binary")
	fmt.Println("  scripts rm <script_name> [--bin]    Remove script or binary")
	fmt.Println("  scripts help                        Show this help message")
	fmt.Println("  scripts -h                          Show this help message")
	fmt.Println("  scripts --help                      Show this help message")
	fmt.Println()
	fmt.Println("COMMANDS:")
	fmt.Println("  <script_name>    Run the specified script (must be in scripts_bin/)")
	fmt.Println("                   Example: scripts gitprune --dry-run")
	fmt.Println()
	fmt.Println("  list             List all available scripts in scripts_bin/ and binaries in ~/opt/programs/")
	fmt.Println("                   Shows script names with executable status and available binaries")
	fmt.Println("                   Example: scripts list")
	fmt.Println()
	fmt.Println("  ready            Make scripts in scripts_bin executable")
	fmt.Println("                   - <script_name> makes script_name.sh in scripts_bin executable")
	fmt.Println("                   - -a or --all makes all .sh files in scripts_bin executable")
	fmt.Println("                   Examples:")
	fmt.Println("                     scripts ready myscript")
	fmt.Println("                     scripts ready -a")
	fmt.Println()
	fmt.Println("  add              Copy script to scripts_bin and make executable")
	fmt.Println("                   Examples:")
	fmt.Println("                     scripts add myscript.sh")
	fmt.Println("                     scripts add ./path/to/script.sh")
	fmt.Println()
	fmt.Println("  compile          Compile source code to binary in ~/opt/programs/")
	fmt.Println("                   Supported: Go, Python, V, Rust, C, C++")
	fmt.Println("                   Use --name to specify custom binary name")
	fmt.Println("                   Examples:")
	fmt.Println("                     scripts compile main.go")
	fmt.Println("                     scripts compile main.go --name myapp")
	fmt.Println("                     scripts compile program.py --name tool")
	fmt.Println("                     scripts compile hello.c -n utility")
	fmt.Println()
	fmt.Println("  rm               Remove script from scripts_bin or binary from ~/opt/programs")
	fmt.Println("                   Use --bin to remove compiled binaries")
	fmt.Println("                   Examples:")
	fmt.Println("                     scripts rm myscript")
	fmt.Println("                     scripts rm --bin myapp")
	fmt.Println()
	fmt.Println("  help             Show this help message")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  scripts list                  # List all available scripts and binaries")
	fmt.Println("  scripts gitprune              # Run gitprune.sh")
	fmt.Println("  scripts test arg1 arg2        # Run test.sh with arguments")
	fmt.Println("  scripts ready myscript        # Make myscript.sh executable")
	fmt.Println("  scripts ready -a              # Make all scripts in scripts_bin executable")
	fmt.Println("  scripts add myscript.sh       # Add script to scripts_bin/")
	fmt.Println("  scripts compile main.go       # Compile Go program to binary")
	fmt.Println("  scripts rm myscript           # Remove myscript.sh from scripts_bin")
	fmt.Println("  scripts rm --bin myapp        # Remove myapp binary from ~/opt/programs")
	fmt.Println("  scripts help                  # Show this help")
	fmt.Println()
	fmt.Println("NOTES:")
	fmt.Println("  - Scripts must be in the scripts_bin/ directory")
	fmt.Println("  - Use 'scripts ready' if you get 'permission denied' errors")
	fmt.Println("  - Compiled binaries are placed in ~/opt/programs/ (add to PATH)")
	fmt.Println("  - PyInstaller required for Python compilation")
	fmt.Println("  - No sudo needed - uses your user permissions")
}

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	// Load configuration
	config, err := loadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	command := os.Args[1]

	// Handle help commands
	if command == "help" || command == "-h" || command == "--help" {
		printHelp()
		return
	}

	if command == "ready" {
		// Handle ready command (make scripts in scripts_bin executable)
		if len(os.Args) < 3 {
			fmt.Println("Usage: scripts ready <script_name> [-a|--all]")
			fmt.Println("  <script_name> makes script_name.sh in scripts_bin executable")
			fmt.Println("  -a|--all makes all .sh files in scripts_bin executable")
			os.Exit(1)
		}

		if os.Args[2] == "-a" || os.Args[2] == "--all" {
			// Make all scripts in scripts_bin executable
			if err := readyScripts([]string{config.ScriptDir}); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			return
		}

		// Handle specific script name (no flags allowed)
		for i := 2; i < len(os.Args); i++ {
			arg := os.Args[i]
			if strings.HasPrefix(arg, "-") {
				fmt.Printf("Unknown flag: %s\n", arg)
				fmt.Println("Usage: scripts ready <script_name>")
				os.Exit(1)
			}
		}

		// Only one script name allowed
		if len(os.Args) != 3 {
			fmt.Println("Usage: scripts ready <script_name>")
			os.Exit(1)
		}

		scriptName := os.Args[2]
		scriptPath := filepath.Join(config.ScriptDir, scriptName+".sh")

		// Check if script exists in scripts_bin
		if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
			fmt.Printf("Script %s not found in scripts_bin (%s)\n", scriptName, config.ScriptDir)
			os.Exit(1)
		}

		// Make the script executable
		if err := makeExecutable(scriptPath); err != nil {
			fmt.Printf("Error making %s executable: %v\n", scriptName, err)
			os.Exit(1)
		}

		fmt.Printf("Made %s executable\n", scriptName)
		return
	}

	if command == "add" {
		// Handle new add command (copy script to scripts_bin)
		if len(os.Args) != 3 {
			fmt.Println("Usage: scripts add <script.sh>")
			fmt.Println("  Copy script to scripts_bin and make executable")
			os.Exit(1)
		}

		scriptPath := os.Args[2]
		if err := addScript(scriptPath, config); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if command == "compile" {
		// Handle compile command
		if len(os.Args) < 3 {
			fmt.Println("Usage: scripts compile <source> [--name <binary_name>]")
			fmt.Println("  Compile source code to binary in ~/opt/programs/")
			fmt.Println("  Supported: Go, Python, V, Rust, C, C++")
			fmt.Println("  --name: specify custom binary name (default: source file name)")
			os.Exit(1)
		}

		sourcePath := os.Args[2]
		binaryName := "" // empty means use default name

		// Parse optional --name flag
		if len(os.Args) >= 4 {
			if os.Args[3] == "--name" || os.Args[3] == "-n" {
				if len(os.Args) != 5 {
					fmt.Println("Usage: scripts compile <source> --name <binary_name>")
					os.Exit(1)
				}
				binaryName = os.Args[4]
			} else {
				fmt.Println("Usage: scripts compile <source> [--name <binary_name>]")
				os.Exit(1)
			}
		}

		if err := compileSource(sourcePath, binaryName, config); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if command == "rm" {
		// Handle rm command
		if len(os.Args) < 3 {
			fmt.Println("Usage: scripts rm <name> [--bin]")
			fmt.Println("  Remove script from scripts_bin/ or binary from ~/opt/programs/")
			fmt.Println("  Use --bin to remove compiled binaries")
			os.Exit(1)
		}

		var name string
		isBinary := false

		// Check if second argument is a flag
		if strings.HasPrefix(os.Args[2], "--") || strings.HasPrefix(os.Args[2], "-") {
			if os.Args[2] == "--bin" || os.Args[2] == "-b" {
				isBinary = true
				if len(os.Args) < 4 {
					fmt.Println("Usage: scripts rm --bin <binary_name>")
					os.Exit(1)
				}
				name = os.Args[3]
			} else {
				fmt.Println("Usage: scripts rm <name> [--bin]")
				os.Exit(1)
			}
		} else {
			// os.Args[2] is the name
			name = os.Args[2]
			// Check for extra arguments
			if len(os.Args) > 3 {
				fmt.Println("Usage: scripts rm <name>")
				os.Exit(1)
			}
		}

		if isBinary {
			// Remove binary from ~/opt/programs
			binPath := filepath.Join(config.BinDir, name)
			if _, err := os.Stat(binPath); os.IsNotExist(err) {
				fmt.Printf("Binary %s not found in %s\n", name, config.BinDir)
				os.Exit(1)
			}

			if err := os.Remove(binPath); err != nil {
				fmt.Printf("Error removing binary %s: %v\n", name, err)
				os.Exit(1)
			}

			fmt.Printf("Removed binary %s\n", name)
		} else {
			// Remove script from scripts_bin
			scriptPath := filepath.Join(config.ScriptDir, name+".sh")
			if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
				fmt.Printf("Script %s not found in %s\n", name, config.ScriptDir)
				os.Exit(1)
			}

			if err := os.Remove(scriptPath); err != nil {
				fmt.Printf("Error removing script %s: %v\n", name, err)
				os.Exit(1)
			}

			fmt.Printf("Removed script %s\n", name)
		}
		return
	}

	if command == "list" {
		// Handle list command (show available scripts and binaries)
		if len(os.Args) > 2 {
			fmt.Println("Usage: scripts list")
			fmt.Println("  Show all available scripts in scripts_bin/ and binaries in ~/opt/programs/")
			os.Exit(1)
		}

		hasOutput := false

		// List scripts
		if _, err := os.Stat(config.ScriptDir); err == nil {
			// Get all .sh files in scripts_bin
			files, err := filepath.Glob(filepath.Join(config.ScriptDir, "*.sh"))
			if err == nil && len(files) > 0 {
				fmt.Println("Available scripts:")
				for _, file := range files {
					scriptName := strings.TrimSuffix(filepath.Base(file), ".sh")
					status := "not executable"
					if isExecutable(file) {
						status = "executable"
					}
					fmt.Printf("  %s (%s)\n", scriptName, status)
				}
				hasOutput = true
			}
		}

		// List binaries
		if _, err := os.Stat(config.BinDir); err == nil {
			// Get all files in bin directory (excluding directories and the scripts binary itself)
			entries, err := os.ReadDir(config.BinDir)
			if err == nil {
				var binaries []string
				for _, entry := range entries {
					if !entry.IsDir() && entry.Name() != "scripts" {
						// Check if it's executable
						binPath := filepath.Join(config.BinDir, entry.Name())
						if isExecutable(binPath) {
							binaries = append(binaries, entry.Name())
						}
					}
				}

				if len(binaries) > 0 {
					if hasOutput {
						fmt.Println()
					}
					fmt.Printf("Available binaries (%s):\n", config.BinDir)
					for _, binary := range binaries {
						fmt.Printf("  %s\n", binary)
					}
					hasOutput = true
				}
			}
		}

		if !hasOutput {
			fmt.Println("No scripts or binaries found.")
			fmt.Printf("Scripts directory: %s\n", config.ScriptDir)
			fmt.Printf("Binaries directory: %s\n", config.BinDir)
		}
		return
	}

	// Handle running scripts
	scriptName := command
	scriptPath := filepath.Join(config.ScriptDir, scriptName+".sh")

	// Check if the script exists
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		fmt.Printf("Script %s not found in %s\n", scriptName, config.ScriptDir)
		os.Exit(1)
	}

	// Check if the script is executable
	if !isExecutable(scriptPath) {
		fmt.Printf("Script %s is not executable. Run 'scripts ready %s' to make it executable.\n", scriptName, scriptName)
		os.Exit(1)
	}

	// Execute the script
	cmd := exec.Command(scriptPath, os.Args[2:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error running script %s: %v\n", scriptName, err)
		os.Exit(1)
	}
}
