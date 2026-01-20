package parser

import (
	"strings"
)

// LineInfo captures the properties of a single scanned line,
// including its raw content, normalized indentation, and detected types.
type LineInfo struct {
	Raw        string
	CleanName  string
	Indent     int    // Normalized indentation (1 unit = 2 spaces or 1 tab)
	Marker     string // "├──", "└──", "|--", "+--", etc.
	IsPathLike bool   // Contains slashes but no spaces?
	IsComment  bool   // Starts with # or //
}

// ScanLines analyzes input text and returns structured info for each line.
// It normalizes tabs to 4 spaces for consistent indentation calculation.
func ScanLines(input string) []LineInfo {
	lines := strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n")
	result := make([]LineInfo, 0, len(lines))

	for _, raw := range lines {
		trim := strings.TrimSpace(raw)
		if trim == "" {
			continue // Skip empty lines here? Or keep them for line number tracking? skipping for now
		}

		// 0. Skip Absolute Windows Paths (for now, or treat as junk?)
		// If a line starts with "C:\" or similar, it's likely a path list entry,
		// BUT tr2rl is designed for relative trees.
		// Let's treat it as a "Path List" candidate if we support absolute paths later.
		// For now, if we see C:\, let's just clean it to relative?
		// Actually, let's just strip the drive letter to make it relative.
		if len(trim) > 3 && trim[1] == ':' && (trim[2] == '\\' || trim[2] == '/') {
			// C:\Path -> Path
			trim = trim[3:]
			raw = raw[3:] // Hacky adjust
		}

		info := LineInfo{Raw: raw}

		// 1. Check for comments
		if strings.HasPrefix(trim, "#") || strings.HasPrefix(trim, "//") {
			info.IsComment = true
			result = append(result, info)
			continue
		}

		// 2. Normalize tabs for calculation
		expanded := strings.ReplaceAll(raw, "\t", "    ")

		// 3. Find Tree Marker
		idx, marker := findBranchMarker(expanded)
		if idx >= 0 {
			info.Marker = marker
			// Depth from marker: count pipes/spaces before it
			prefix := expanded[:idx]
			// Count logical depth (simplified: 1 pipe or 2-4 spaces = 1 level)
			info.Indent = countGraphicDepth(prefix)

			// Clean name: everything after marker
			rest := expanded[idx+len(marker):]
			info.CleanName = strings.TrimSpace(rest)
		} else {
			// Indented list fallback
			leadingSpaces := len(expanded) - len(strings.TrimLeft(expanded, " "))
			info.Indent = leadingSpaces / 4 // Assume 4 spaces per level for indent mode
			info.CleanName = trim
		}

		// 4. Strip inline comments from name
		// Supports: #, //, <--, (comment)
		info.CleanName = stripInlineComment(info.CleanName)
		info.CleanName = strings.TrimSpace(info.CleanName)
		info.CleanName = strings.ReplaceAll(info.CleanName, "\\", "/") // Normalize Windows paths
		info.CleanName = strings.TrimSuffix(info.CleanName, "/")       // Remove trailing slash for consistency (added back by Kind)

		// Extra aggression: if it ends with " <-- ...", strip it
		if idx := strings.Index(info.CleanName, " <--"); idx != -1 {
			info.CleanName = info.CleanName[:idx]
		}

		// 5. Clean list bullets
		info.CleanName = strings.TrimPrefix(info.CleanName, "- ")
		info.CleanName = strings.TrimPrefix(info.CleanName, "* ")
		info.CleanName = strings.TrimSpace(info.CleanName)

		// 6. Path-like check
		// Contains slash, no spaces (unless escaped, which we ignore for now)
		if strings.Contains(info.CleanName, "/") && !strings.Contains(info.CleanName, " ") {
			info.IsPathLike = true
		}

		result = append(result, info)
	}
	return result
}

// Helper to count visual depth from tree graphics like "│   │   "
func countGraphicDepth(prefix string) int {
	// Simple heuristic: Count "│" or "|"
	// But we must be careful: "       |--" might be indented.
	// 4 spaces = 1 indent.
	// 1 pipe = 1 indent.
	// Mixing them is hard.
	// Let's normalize everything to spaces first?
	// Or just count pipes + spaces / 4?

	pipes := 0
	for _, r := range prefix {
		if r == '│' || r == '|' || r == '├' || r == '└' { // Include all vertical-ish markers
			pipes++
		}
	}

	// If messy input like "|       |--", prefix is "|       ".
	// 1 pipe + 7 spaces.
	// 7 spaces -> ~2 indents?
	// Total depth = 1 + 2 = 3?
	// Let's try to be smart.
	// Calculate effective width.
	// pipe = 4 spaces width in this context?
	// Standard tree is "│   ". Pipe uses 1 char, plus 3 spaces = 4 chars.

	// Let's assume standard "tree" width of 4 chars per level.
	// We count total visual width (pipe=1, space=1) and divide by 4.
	// But if we have just pipes "|||", that's 3 levels but only width 3.
	// So: Max(pipes, visualWidth/4)?

	width := len(prefix)
	depthByHash := width / 4

	if pipes > depthByHash {
		return pipes
	}
	return depthByHash
}
