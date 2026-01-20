package fs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/tr2rl/tr2rl/internal/parser"
)

func TestApply_Integration(t *testing.T) {
	// Setup temp directory
	tmpDir, err := os.MkdirTemp("", "tr2rl-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	nodes := []parser.Node{
		{Path: "src/main.go", Kind: parser.File},
		{Path: "docs/readme.md", Kind: parser.File},
		{Path: "tests/", Kind: parser.Dir},
	}

	// 1. Test Dry Run (Should not create files)
	err = Apply(tmpDir, nodes, ApplyOptions{DryRun: true, Force: false})
	if err != nil {
		t.Errorf("Dry run failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(tmpDir, "src/main.go")); !os.IsNotExist(err) {
		t.Error("Dry run created file 'src/main.go', but shouldn't have")
	}

	// 2. Test Real Run
	err = Apply(tmpDir, nodes, ApplyOptions{DryRun: false, Force: false})
	if err != nil {
		t.Errorf("Apply failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(tmpDir, "src/main.go")); os.IsNotExist(err) {
		t.Error("Apply failed to create 'src/main.go'")
	}
	if _, err := os.Stat(filepath.Join(tmpDir, "docs/readme.md")); os.IsNotExist(err) {
		t.Error("Apply failed to create 'docs/readme.md'")
	}
	if _, err := os.Stat(filepath.Join(tmpDir, "tests")); os.IsNotExist(err) {
		t.Error("Apply failed to create 'tests' directory")
	}

	// 3. Test Overwrite Protection
	// Write content to file
	mainPath := filepath.Join(tmpDir, "src/main.go")
	os.WriteFile(mainPath, []byte("content"), 0644)

	// Run Apply again (no force)
	err = Apply(tmpDir, nodes, ApplyOptions{DryRun: false, Force: false})
	if err != nil {
		t.Error(err)
	}

	content, _ := os.ReadFile(mainPath)
	if string(content) != "content" {
		t.Error("Apply overwrote file without --force")
	}

	// 4. Test Force
	err = Apply(tmpDir, nodes, ApplyOptions{DryRun: false, Force: true})
	if err != nil {
		t.Error(err)
	}

	content, _ = os.ReadFile(mainPath)
	if string(content) != "" {
		// os.Create truncates, so it should be empty now
		t.Error("Apply with --force did not truncate file")
	}
}
