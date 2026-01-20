# Package `parser`

The heart of `tr2rl`. This package converts messy text input into a normalized list of file system nodes.

## Core Logic (`Parse`)

The parser uses a "Magic Parsing" strategy that employs heuristics to understand user intent.

### 1. Unified Scanning (`scanner.go`)
Instead of strict regex, we scan each line to identify:
*   **Indentation Level**: Normalized to 4 spaces (tabs are expanded).
*   **Tree Markers**: `├──`, `└──`, `|--`, `+--`.
*   **Path-like characteristics**: Does it contain slashes? Extensions?

### 2. Strategy Selection
The parser analyzes the first pass to decide between:
*   **Tree Strategy**: Standard hierarchical parsing.
*   **Path List Strategy**: If >50% of lines look like paths (`src/main.go`), we treat it as a flat list.

### 3. Forgiving Tree Construction
*   **Indentation Jumps**: If a line jumps from indent 0 to 5, we clamp it to `parent_depth + 1` instead of erroring.
*   **Missing Roots**: If the text starts with children but no root, we infer a root or attach to current directory.
*   **Junk Skipping**: Lines that don't look like files/dirs and have no markers are skipped (e.g., random sentence text).

## Supported Formats

| Format | Example |
| :--- | :--- |
| **Unicode Tree** | `├── src/` |
| **Windows Tree** | `|-- src\` |
| **Indented** | `  src` |
| **Path List** | `src/main.go` |
