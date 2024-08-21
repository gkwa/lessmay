# lessmay

Lessmay is a command-line app to help manage sync conflicts in Obsidian vaults. It finds and displays differences between sync conflict files and their original versions, allowing users to easily resolve conflicts.

## Features

- Automatically finds sync conflict files in specified directories
- Compares conflict files with their original versions
- Deletes identical conflict files
- Displays differences for non-identical files
- Supports multiple directories
- Customizable skip paths to ignore certain directories

## Installation

To install Lessmay, use the following command:

```
go install github.com/gkwa/lessmay@latest
```

## Usage

### Basic Usage

To show conflicts in the default Obsidian vault location:

```
lessmay show-conflicts
```

### Specify Custom Directories

To show conflicts in specific directories:

```
lessmay show-conflicts /path/to/vault1 /path/to/vault2
```

### Use Custom Default Path

To use a custom default path:

```
lessmay show-conflicts --default-path /path/to/custom/vault
```

### Skip Specific Paths

To skip specific paths when searching for conflicts:

```
lessmay show-conflicts --skip-path .trash --skip-path .archive
```

### Verbose Output

For more detailed output:

```
lessmay show-conflicts -v
```

### JSON Logging

To output logs in JSON format:

```
lessmay show-conflicts --log-format json
```

## Configuration

Lessmay can be configured using a configuration file located at `$HOME/.lessmay.yaml`. You can specify the following options:

```yaml
verbose: true
log-format: json
skip-path:
  - .trash
  - .archive
```
