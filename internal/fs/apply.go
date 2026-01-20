// Package fs handles the actual filesystem operations for tr2rl.
// It includes safety features like dry-run mode and overwrite protection
// to ensure users don't accidentally destroy their work.
package fs

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/tr2rl/tr2rl/internal/content"
	"github.com/tr2rl/tr2rl/internal/parser"
)

type ApplyOptions struct {
	DryRun   bool
	Force    bool
	Populate bool
}

// Apply materializes the parsed nodes into the filesystem at rootDir.
func Apply(rootDir string, nodes []parser.Node, opts ApplyOptions) error {
	for _, node := range nodes {
		fullPath := filepath.Join(rootDir, node.Path)

		if opts.DryRun {
			// In dry-run, just print what we would do
			suffix := ""
			if node.Kind == parser.Dir {
				suffix = "/"
			}
			fmt.Printf("[DRY-RUN] Create %s%s\n", fullPath, suffix)
			continue
		}

		if node.Kind == parser.Dir {
			if err := os.MkdirAll(fullPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", fullPath, err)
			}
		} else {
			// Ensure parent dir exists
			parent := filepath.Dir(fullPath)
			if err := os.MkdirAll(parent, 0755); err != nil {
				return fmt.Errorf("failed to create parent dir %s: %w", parent, err)
			}

			// Check if file exists
			if _, err := os.Stat(fullPath); err == nil {
				if !opts.Force {
					fmt.Printf("[SKIP] File exists: %s (use --force to overwrite)\n", node.Path)
					continue
				}
			}

			// Prepare content
			data := ""
			if opts.Populate {
				data = content.GetContent(fullPath)
			}

			// Create/Truncate file
			// Use os.WriteFile for simplicity
			if err := os.WriteFile(fullPath, []byte(data), 0644); err != nil {
				return fmt.Errorf("failed to create file %s: %w", fullPath, err)
			}
			fmt.Printf("[OK] Created %s\n", node.Path)
		}
	}
	return nil
}
