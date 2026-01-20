// Package parser provides the core logic for checking and converting text-based
// tree specifications into a structured list of nodes.
//
// It employs a "Magic Parsing" strategy that attempts to recover structure from
// messy inputs, such as mixed indentation, missing markers, and path lists.
package parser

import (
	"path"
	"strings"
)

// Parse turns a text tree (Windows ASCII tree, Unicode tree, indented lists, mixed)
// into a flat list of Nodes with normalized slash-separated paths.
func Parse(input string) Result {
	lines := ScanLines(input)

	// Pre-filter lines to remove comments and junk
	// This ensures lines[0] is the actual first tree item for root detection.
	validLines := make([]LineInfo, 0, len(lines))
	for _, l := range lines {
		if l.IsComment {
			continue
		}
		// Also skip empty CleanName here?
		if l.CleanName == "" {
			continue
		}

		// Windows Tree Headers
		if strings.HasPrefix(l.Raw, "Folder PATH listing") ||
			strings.HasPrefix(l.Raw, "Volume serial number") {
			continue
		}

		// Windows Drive Anchor "C:." -> Treat as "." (or skip if we want implicit root)
		// Actually, let's normalize it to "."
		if len(l.CleanName) == 3 && l.CleanName[1] == ':' && l.CleanName[2] == '.' {
			l.CleanName = "."
			l.IsPathLike = true // Force it to look like a path so it's not filtered later
			// We need to modify 'l' but 'l' is a copy.
			// But we append 'l' to validLines. So if we modify l, it works.
		}

		validLines = append(validLines, l)
	}
	lines = validLines

	if len(lines) == 0 {
		return Result{}
	}

	result := Result{
		Warnings: make([]string, 0),
	}

	// Phase 1: Heuristic Analysis
	pathLikeCount := 0
	markerCount := 0
	for _, l := range lines {
		if l.IsPathLike {
			pathLikeCount++
		}
		if l.Marker != "" {
			markerCount++
		}
	}

	// Decision: Is this a Path List or a Tree?
	isPathList := false
	if pathLikeCount > len(lines)/2 && markerCount == 0 {
		isPathList = true
		// result.Warnings = append(result.Warnings, "Detected Path List format")
	}

	if isPathList {
		result.Nodes = parsePathList(lines)
	} else {
		result.Nodes, result.Warnings = parseTree(lines, result.Warnings)
	}

	// Normalize Output
	norm := make([]string, 0, len(result.Nodes))
	for _, n := range result.Nodes {
		p := n.Path
		if n.Kind == Dir && !strings.HasSuffix(p, "/") {
			p += "/"
		}
		norm = append(norm, p)
	}
	result.Normalized = strings.Join(norm, "\n")

	return result
}

func parsePathList(lines []LineInfo) []Node {
	nodes := make([]Node, 0, len(lines))
	seen := make(map[string]bool)

	for _, l := range lines {
		clean := strings.TrimSpace(l.Raw) // Use raw for path lists, but trim
		// Remove ./ prefix if present
		clean = strings.TrimPrefix(clean, "./")
		clean = strings.ReplaceAll(clean, "\\", "/")

		if clean == "" {
			continue
		}

		kind := File
		if strings.HasSuffix(clean, "/") {
			kind = Dir
			clean = strings.TrimSuffix(clean, "/")
		}

		if !seen[clean] {
			nodes = append(nodes, Node{Path: clean, Kind: kind})
			seen[clean] = true
		}
	}
	return nodes
}

func parseTree(lines []LineInfo, warnings []string) ([]Node, []string) {
	// Root Handling: Check if first line is a root
	// A line is ROOT if:
	// 1. Depth is 0 (or very low compared to next)
	// 2. It has no markers
	// 3. Next line is deeper OR has markers

	nodes := make([]Node, 0, len(lines))
	stack := make([]string, 0, 32)

	// Heuristic: If first line has markers, it's probably NOT a root (it's a child of CWD)
	// But if first line has NO markers, and second line DOES, first line is Root.

	rootIdx := -1
	if len(lines) > 0 {
		l0 := lines[0]
		if l0.Marker == "" {
			// Check if it looks like a root wrapper?
			if len(lines) > 1 {
				l1 := lines[1]
				// If next line has indentation OR markers, we are the root.
				// Even if current line has no slash, if it heads a tree, it's a dir.
				// Specially, if l1 has a marker like "|--", that's indentation 0 usually.
				// But physically, if l0 is above it, l0 is the parent.
				if l1.Indent >= l0.Indent || l1.Marker != "" {
					rootIdx = 0
				}
			} else if len(lines) == 1 && isDirLike(lines[0].CleanName) {
				// Single line directory
				rootIdx = 0
			}
		}
	}

	// Initialize state
	prevIndent := -1
	if rootIdx == 0 {
		// We have a declared root
		rootName := strings.TrimSuffix(lines[0].CleanName, "/")
		stack = append(stack, rootName)
		nodes = append(nodes, Node{Path: rootName, Kind: Dir})
		prevIndent = lines[0].Indent
		// Start processing children from index 1
		lines = lines[1:]
	} else {
		// Implicit root (current directory), start from 0
		prevIndent = -1
	}
	// Determine Indent Offset for Flat Roots
	// If explicit root exists (Indent 0) and first child is also Indent 0,
	// we need to shift all children by +1 so they nest inside root.
	indentOffset := 0
	if rootIdx == 0 && len(lines) > 0 {
		firstChildIndent := lines[0].Indent
		rootIndent := prevIndent // Captured from l0 before slicing
		if firstChildIndent <= rootIndent {
			indentOffset = 1 + (rootIndent - firstChildIndent)
		}
	}
	prevRawIndent := -1

	for _, l := range lines {
		indent := l.Indent
		name := l.CleanName

		if name == "" {
			continue
		}

		// 2. Junk Filter: If it has spaces, isn't a known file, and isn't path-like, skip it.
		// EXCEPTION: If it has a valid marker, TRUST IT.
		if l.Marker == "" && !l.IsPathLike && strings.Contains(name, " ") && !looksLikeFile(name) {
			// Special check: sometimes "Program Files" is valid?
			// But "random garbage line" is not.
			// If it has NO extension and DOES have spaces, assume junk unless flagged otherwise.
			// warnings = append(warnings, fmt.Sprintf("Skipping potential junk line: '%s'", l.Raw))
			continue
		}

		// Forgiving Jump Protection logic fixed:
		// If we clamped the previous line, we must respect that for siblings.
		// If current RawIndent == PrevRawIndent, then LogicalIndent must be PrevLogicalIndent.

		if l.Indent == prevRawIndent {
			indent = prevIndent
		} else {
			// Normal clamp
			if indent > prevIndent+1 {
				indent = prevIndent + 1
			}
		}

		// Effective Indent: apply offset for flat roots
		effectiveIndent := indent + indentOffset

		// ... stack logic ...
		targetStackLen := effectiveIndent

		if rootIdx == 0 {
			// If we have an explicit root, we must never pop it off the stack.
			// The root is always at stack[0], so stack len must be at least 1.
			if targetStackLen < 1 {
				targetStackLen = 1
			}
		}

		if targetStackLen < 0 {
			targetStackLen = 0
		}

		for len(stack) > targetStackLen {
			stack = stack[:len(stack)-1]
		}

		// Infer Kind
		kind := File
		if isDirLike(name) {
			kind = Dir
			name = strings.TrimSuffix(name, "/")
		} else {
			// Lookahead check could go here
			if !looksLikeFile(name) {
				// default to file for safety
			}
		}

		stack = append(stack, name)
		fullPath := path.Join(stack...)
		nodes = append(nodes, Node{Path: fullPath, Kind: kind})

		prevIndent = indent
		prevRawIndent = l.Indent
	}

	// Post-pass: Fix "File" that became a parent
	// If Node A is a parent of Node B, Node A must be a Directory.
	// Since `nodes` is ordered, we can check basic containment?
	// Or better: In the loop, when we push to stack, we are declaring it a parent.
	// So anything remaining in `stack` (except the very last item) IS acting as a directory.

	// Let's fix the loop to handle "Make Parent" logic.
	// Actually, `stack` represents the current active directory chain.
	// So every time we push to stack, the *previous* stack top (if any) was effectively a specific node.
	// But `stack` only stores strings.
	// We need to update the `Node` kind in `nodes` list.

	// Map-based correction
	pathToKind := make(map[string]int) // Index in nodes
	for i, n := range nodes {
		pathToKind[n.Path] = i
	}

	for _, n := range nodes {
		// If "A/B" exists, then "A" must be a Dir
		parent := path.Dir(n.Path)
		if parent != "." && parent != "/" {
			if idx, ok := pathToKind[parent]; ok {
				nodes[idx].Kind = Dir
			}
		}
	}

	return nodes, warnings
}
