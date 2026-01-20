# tr2rl Architecture

## Overview
`tr2rl` is a static Go binary designed for portability and robustness. It follows the Standard Go Project Layout.

## Data Flow

1.  **Input**: User provided text (File, Stdin, Clipboard).
2.  **Parser**: `internal/parser` converts text -> `[]Node` (Flat list of paths + types).
3.  **Command Layer**: `cmd/` decides what to do with nodes (Build, Format, Verify).
4.  **Action Layer**:
    *   **Build**: `internal/fs` applies nodes to disk.
    *   **Format**: `internal/printer` renders nodes as ASCII tree.

## Directory Structure

*   **/cmd**: Entry points for CLI commands (cobra).
*   **/internal**: Private library code.
    *   **/parser**: "Magic Parser" logic.
    *   **/fs**: Safer filesystem operations (Dry-run logic).
    *   **/printer**: ASCII tree generation.
    *   **/templates**: Built-in project blueprints.
    *   **/clipboard**: Cross-platform clipboard access (no CGO).
*   **/testdata**: Fixtures for integration testing.

## Key Design Decisions

### 1. Flat Node Representation
We convert trees to a flat list of normalized paths early.
*   *Input*: `src -> main.go`
*   *Internal*: `Node{Path: "src/main.go", Kind: File}`
*   *Why*: Simplifies "Path List" parsing and makes sorting/deduplication trivial.

### 2. No CGO
We use `os/exec` to call system clipboard tools (`powershell`, `pbpaste`) instead of linking C libraries. This ensures the binary runs on any machine without runtime dependencies.
