# tr2rl (Trees to Reality)

**tr2rl** is a robust CLI tool that turns text-based tree specifications (like the output of the `tree` command) into actual directory structures on your disk.

It is designed to be "unbreakable": it accepts messy inputs, mixed indentation, and broken ASCII art, and safely materializes them.

## Features

*   **Magic Parsing**: Understands Unicode trees, ASCII trees (`|--`), indented lists, and path lists.
*   **Safety First**: Use `--dry-run` to preview changes. Requires `--force` to overwrite files.
*   **Cross-Platform**: Single static binary for Windows, Linux, and macOS.
*   **Templates**: Built-in blueprints for common project types (Go, React, Python).
*   **Formatter**: Clean up messy text trees into standard Unicode format.

## Why tr2rl? (Comparison)

| Feature | `tr2rl` | Shell Scripts (`mkdir -p`) | Other Tools (`tree-cli`) |
| :--- | :--- | :--- | :--- |
| **Messy Input** | ✅ **Magic Parser** handles anything | ❌ Fails on mismatched indents | ❌ Requires strict JSON/Format |
| **Safety** | ✅ **Dry-Run by default** | ❌ Destructive immediately | ⚠️ Varies |
| **Portability** | ✅ **Static Binary** (No deps) | ⚠️ Requires Bash/Python | ❌ Requires Node/Python |
| **Templates** | ✅ **Built-in Registry** | ❌ Manual | ❌ Rare |

## Installation

Download the latest binary from the releases page (or `dist/` folder if building locally).

### Build from source
```bash
go build -o tr2rl.exe .
```

## Usage

### 1. Build a directory from text
Create a file `spec.txt` with your tree:
```text
my-project/
├── src/
│   └── main.go
└── README.md
```

Run build (dry-run by default):
```bash
tr2rl build spec.txt ./output
```

Actually create files:
```bash
tr2rl build spec.txt ./output --dry-run=false
```

Generate boilerplate content (e.g. `package main` for Go files):
```bash
tr2rl build spec.txt ./output --populate
```

### 2. Use Built-in Templates
Generate a new React project:
```bash
tr2rl template show react-vite | tr2rl build - ./my-react-app --dry-run=false
```

### 3. Format/Clean a Tree
Turn a messy list into a clean tree:
```bash
tr2rl format messy_list.txt
```

### 4. Input Sources
*   **File**: `tr2rl build file.txt`
*   **Stdin (Pipe)**: `cat file.txt | tr2rl build -`
*   **Clipboard**: `tr2rl build --clipboard` (No file needed!)

### 5. Pro Tip: The "One-Liner"
Copy a tree structure to your clipboard, then run:
```bash
tr2rl build --clipboard --populate --dry-run=false
```
This instantly builds the folders **AND** generates boilerplate content (`main.go`, `index.html`, etc.) for you.


## Supported Formats

**Unicode Tree:**
```text
root/
├── child/
└── file
```

**ASCII / Windows Tree:**
```text
root
|-- child
`-- file
```

**Indented List:**
```text
root
  child
  file
```

**Path List:**
```text
root/child/
root/file
```

## License
MIT
