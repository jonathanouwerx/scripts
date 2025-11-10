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

## Installation

Follow these steps to set up the Scripts Tool from scratch.

### Prerequisites

- **Go 1.21+** - Download from [golang.org](https://golang.org/dl/)
- **Git** - For cloning the repository

### Step 1: Clone the Repository

```bash
git clone https://github.com/jonathanouwerx/scripts.git
cd scripts
```

### Step 2: Build the Binary

```bash
# Build the scripts binary
make build

# Or build manually with Go
go build -o scripts .
```

### Step 3: Install the Binary

```bash
# Option 1: Move to a directory in your PATH (recommended)
sudo mv scripts /usr/local/bin/
```

### Step 4: Verify Installation

```bash
# Check that scripts is available
scripts --help

# You should see the help output with available commands
```

### Step 5: Initial Setup

The tool will automatically create its configuration when first run. It will:

1. Create a `scripts_bin/` directory in your scripts repository
2. Set up `~/opt/programs/` for compiled binaries
3. Create a `.config.json` file with your paths

```bash
# First run will auto-configure
scripts list

# Should show empty lists (no scripts or binaries yet)
```

### Step 6: Add Scripts Directory to PATH (Optional)

For easier access to compiled binaries:

```bash
# Add to your shell profile
echo 'export PATH="$HOME/opt/programs:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

## Configuration

The tool **automatically creates** a `.config.json` file in the scripts directory on first run with these default paths:
- `scriptDir`: `~/code/personal/scripts/scripts_bin` (where your scripts are stored)
- `binDir`: `~/opt/programs` (where compiled binaries are placed)

**Note:** `.config.json` is gitignored - each user gets their own personalized configuration.
