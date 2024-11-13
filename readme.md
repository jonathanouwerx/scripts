# Scripts

This CLI executes scripts in scripts_bin using the ```scripts``` command.

The scripts binary itself is stored in ~/opt/programs.

## Usage

```bash
scripts <script_name> <additional_args>
```

## Installation

Copy the scripts binary to ~/opt/programs.

```bash
cp scripts ~/opt/programs
```

Add the following to your .bashrc or .bash_profile:

```bash
export PATH=$PATH:~/opt/programs
```
