package parser

import (
	"strings"
)

// Helper functions restored from previous parse.go

func findBranchMarker(s string) (idx int, marker string) {
	// Prefer longer markers first.
	markers := []string{
		"├───", "└───", "├─", "└─", // Windows Unicode (3 dashes)
		"├──", "└──", // Linux Standard (2 dashes)
		"|---", "+---", "\\---", // Extended Windows/ASCII
		"|--", "+--", "\\--", "`--", // Standard ASCII
		"┠──", "┗━━", // Variants
		"├ ", "└ ", "| ", "+ ", // Broken variants
	}
	bestIdx := -1
	bestMarker := ""
	for _, m := range markers {
		i := strings.LastIndex(s, m)
		if i >= 0 {
			// Logic: Prefer RIGHT-MOST match.
			// If tied (same index), prefer LONGER marker (e.g. "├──" over "├─").
			if i > bestIdx {
				bestIdx = i
				bestMarker = m
			} else if i == bestIdx {
				if len(m) > len(bestMarker) {
					bestMarker = m
				}
			}
		}
	}
	return bestIdx, bestMarker
}

func hasBranchMarker(raw string) bool {
	_, m := findBranchMarker(raw)
	return m != ""
}

func isDirLike(name string) bool {
	return strings.HasSuffix(strings.TrimSpace(name), "/")
}

func looksLikeFile(name string) bool {
	n := strings.TrimSpace(name)
	if strings.HasSuffix(n, "/") {
		return false
	}
	// Common file heuristics: contains '.' not at beginning/end, or known filenames.
	if n == "Makefile" || n == "Dockerfile" || n == "LICENSE" || n == "README" {
		return true
	}
	dot := strings.LastIndex(n, ".")
	return dot > 0 && dot < len(n)-1
}

func stripInlineComment(s string) string {
	// Remove " # ..." and " // ..." when preceded by whitespace.
	// This keeps URLs like "http://..." mostly intact because they usually
	// won’t have a leading whitespace before the slashes in a filename.
	if i := strings.Index(s, " #"); i >= 0 {
		return s[:i]
	}
	if i := strings.Index(s, "\t#"); i >= 0 {
		return s[:i]
	}
	if i := strings.Index(s, " //"); i >= 0 {
		return s[:i]
	}
	if i := strings.Index(s, "\t//"); i >= 0 {
		return s[:i]
	}
	return s
}
