# Scripts Tool

A powerful CLI tool for managing shell scripts and compiling binaries from multiple programming languages.

## Features

### Script Management
- **`scripts <name>`** - Run shell scripts from `scripts_bin/`
- **`scripts ready <script_name>`** - Make scripts in `scripts_bin` executable
- **`scripts add <script.sh>`** - Copy script to `scripts_bin/` and make executable
- **`scripts rm <script_name>`** - Remove script from `scripts_bin/`

### Binary Compilation & Management
- **`scripts compile <source>`** - Compile source code to executable binaries
- **`scripts compile <source> --name <custom_name>`** - Compile with custom binary name
- **`scripts rm --bin <binary_name>`** - Remove compiled binary from `~/opt/programs/`

### Supported Languages
- **Go** (.go)
- **Python** (.py) - requires PyInstaller
- **V** (.v)
- **Rust** (.rs) - supports both Cargo projects and single files
- **C** (.c)
- **C++** (.cpp, .cc, .cxx)

Compiled binaries are placed in `~/opt/programs/` and can be run directly from PATH.

## Configuration

Auto-creates `.config.json` in the scripts directory with:
- `scriptDir`: Path to scripts directory
- `binDir`: Path for compiled binaries (`~/opt/programs`)

**Setup:**
```bash
# Copy the example config
cp .config.json.example .config.json

# Edit paths as needed
# scriptDir: Where your scripts are stored
# binDir: Where compiled binaries go
```

**Note:** `.config.json` is gitignored - each user has their own configuration.

## Usage

```bash
# Run scripts
scripts my-script arg1 arg2

# Make scripts executable
scripts ready myscript        # single script
scripts ready -a             # all scripts

# Add scripts to scripts_bin
scripts add myscript.sh

# Remove scripts
scripts rm myscript

# Compile programs
scripts compile main.go
scripts compile main.go --name myapp
scripts compile program.py
scripts compile hello.c

# Remove binaries
scripts rm --bin myapp

# Show help
scripts help
```

## Requirements

- **PyInstaller** for Python compilation: `pip install pyinstaller`
- **Language compilers** (Go, GCC, etc.) must be installed for compilation
- Add `~/opt/programs` to your PATH for direct binary execution