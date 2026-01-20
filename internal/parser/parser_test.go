package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParse_Golden(t *testing.T) {
	// These tests expect the golden files to exist in the project root.
	// We look up two directories because tests run in internal/parser/
	rootDir := "../../testdata/"

	tests := []struct {
		name     string
		filename string
	}{
		{"Unicode Tree", "correct_unicode.tree"},
		{"Windows Tree", "correct_windows.tree"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(rootDir, tt.filename)
			content, err := os.ReadFile(path)
			if err != nil {
				t.Skipf("Golden file not found: %v", err)
			}

			res := Parse(string(content))

			if len(res.Nodes) == 0 {
				t.Errorf("Expected nodes, got none")
			}
			if len(res.Warnings) > 0 {
				t.Logf("Warnings: %v", res.Warnings)
			}

			// Basic sanity check: ensure main.py exists
			found := false
			for _, n := range res.Nodes {
				if strings.HasSuffix(n.Path, "main.py") {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Parser failed to find 'main.py' in %s", tt.filename)
			}
		})
	}
}

func TestParse_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int // number of nodes
	}{
		{
			name: "Path List (Siblings)",
			input: `
/etc/nginx/
/var/log/
`,
			expected: 2,
		},
		{
			name: "Mixed Windows Paths",
			input: `
src\main.go
src\utils\helper.go
`,
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := Parse(tt.input)
			if len(res.Nodes) != tt.expected {
				t.Errorf("Expected %d nodes, got %d. Nodes: %v", tt.expected, len(res.Nodes), res.Nodes)
			}
		})
	}
}
