package content

import (
	"fmt"
	"path/filepath"
	"strings"
)

// GetContent returns smart default content for a file based on its name/extension.
func GetContent(path string) string {
	base := filepath.Base(path)
	ext := strings.ToLower(filepath.Ext(path))

	// 1. Exact Filename Matches
	switch strings.ToLower(base) {
	case "makefile":
		return "all:\n\t@echo 'Hello World'\n"
	case "dockerfile":
		return "FROM alpine:latest\nCMD [\"echo\", \"Hello World\"]\n"
	case ".gitignore":
		return "# Ignore list\n.DS_Store\nnode_modules/\ndist/\nbin/\n"
	case "license", "license.txt", "license.md":
		return "MIT License\n\nCopyright (c) 2026\n"
	}

	// 2. Extension Matches
	switch ext {
	// Go
	case ".go":
		if base == "main.go" {
			return "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello World\")\n}\n"
		}
		// Try to guess package name from parent dir?
		// For now simple package declaration
		return "package " + guessPackage(path) + "\n"

	// Web
	case ".html":
		return "<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n    <meta charset=\"UTF-8\">\n    <title>New Page</title>\n</head>\n<body>\n    <h1>Hello World</h1>\n</body>\n</html>\n"
	case ".css":
		return "body {\n    font-family: sans-serif;\n    margin: 0;\n}\n"
	case ".js":
		return "console.log('Hello World');\n"
	case ".jsx", ".tsx":
		return "import React from 'react';\n\nexport const Component = () => {\n    return <div>Hello</div>;\n};\n"
	case ".json":
		return "{}\n"

	// Python
	case ".py":
		if base == "main.py" || base == "app.py" {
			return "def main():\n    print(\"Hello World\")\n\nif __name__ == \"__main__\":\n    main()\n"
		}
		return "# New Python Module\n"

	// Shell
	case ".sh":
		return "#!/bin/bash\nset -euo pipefail\n\necho \"Hello from script\"\n"

	// Config / Data
	case ".yaml", ".yml":
		return "version: '1.0'\n"
	case ".md":
		title := strings.TrimSuffix(base, ext)
		return fmt.Sprintf("# %s\n\nDescription goes here.\n", strings.Title(title))
	case ".txt":
		return "" // Keep text empty? Or "New Text"? Empty is probably better for .txt
	}

	// Default: Empty
	return ""
}

func guessPackage(path string) string {
	// parent dir name
	dir := filepath.Base(filepath.Dir(path))
	if dir == "." || dir == "/" {
		return "main"
	}
	// Sanitize package name (lowercase, no symbols)
	clean := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			return r
		}
		if r >= 'A' && r <= 'Z' {
			return r + 32 // toLower
		}
		return -1 // drop
	}, dir)

	if clean == "" {
		return "pkg"
	}
	return clean
}
