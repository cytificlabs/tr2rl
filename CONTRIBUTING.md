# Contributing to tr2rl

Thank you for your interest in contributing to `tr2rl`!

## Development Setup

1.  **Prerequisites**: Go 1.21+
2.  **Clone**: `git clone https://github.com/tr2rl/tr2rl`
3.  **Run**: `go run main.go --help`

## Project Structure

*   `cmd/`: CLI commands (cobra).
*   `internal/parser/`: The "Magic Parsing" logic.
*   `internal/fs/`: Filesystem operations (Dry-run safety).
*   `testdata/`: Text files for integration testing.

## Running Tests
 
```bash
go test ./...
```

To run a specific verification test:
```bash
go run main.go spec testdata/mixed_indent.txt
```

## Pull Requests

1.  Fork the repo.
2.  Create a branch (`feature/my-cool-feature`).
3.  Ensure `go test ./...` passes.
4.  Submit a PR!
