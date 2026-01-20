// Package printer handles the visual rendering of the tree structure.
package printer

import (
	"fmt"
	"sort"
	"strings"

	"github.com/tr2rl/tr2rl/internal/parser"
)

// PrintTree constructs and prints a hierarchical Unicode tree from a flat list of nodes.
// It handles directory grouping, sorting, and prefix generation (e.g., "├── ").
func PrintTree(nodes []parser.Node) {
	if len(nodes) == 0 {
		return
	}

	// 1. Build a Tree map structure for printing
	// Map parent path -> list of children
	childrenMap := make(map[string][]parser.Node)
	roots := make([]parser.Node, 0)

	// Track all paths to identify roots (items with no parent in the list)
	// Actually, we can just use path.Dir() logic.
	// But since input can be fragmented, we need to be careful.
	// We'll rely on our parser's normalization which ensures full paths.

	// First, find the common root or strictly top-level items.
	// The parser outputs relative paths. "src/main.go" -> parent is "src".
	// "src" -> parent is ".".

	nodeMap := make(map[string]parser.Node)
	for _, n := range nodes {
		// Normalize: strip trailing slash for key
		key := strings.TrimSuffix(n.Path, "/")
		nodeMap[key] = n
	}

	for _, n := range nodes {
		key := strings.TrimSuffix(n.Path, "/")
		// parent := strings.TrimSuffix(path.Dir(key), "/") // wrong, path.Dir handles separators

		// Manual parent finding to respect forward slashes universally
		lastSlash := strings.LastIndex(key, "/")
		if lastSlash == -1 {
			// Top level
			roots = append(roots, n)
		} else {
			parentPath := key[:lastSlash]
			childrenMap[parentPath] = append(childrenMap[parentPath], n)

			// If parent doesn't exist in our node list, we should probably print it implicitly?
			// The parser should ideally have filled gaps, but if not:
			// For now assume strictly parsed nodes.
		}
	}

	// Sort roots
	sortNodes(roots)

	// Print recursively
	for i, root := range roots {
		printNode(root, "", i == len(roots)-1, childrenMap)
	}
}

func printNode(node parser.Node, prefix string, isLast bool, childrenMap map[string][]parser.Node) {
	// Prepare current line marker
	marker := "├── "
	if isLast {
		marker = "└── "
	}

	name := node.Path
	// Extract basename
	if idx := strings.LastIndex(strings.TrimSuffix(name, "/"), "/"); idx != -1 {
		name = strings.TrimSuffix(name, "/")
		name = name[idx+1:]
	}
	name = strings.TrimSuffix(name, "/")
	if node.Kind == parser.Dir {
		name += "/"
	}

	fmt.Printf("%s%s%s\n", prefix, marker, name)

	// Calculate prefix for children
	childPrefix := prefix
	if isLast {
		childPrefix += "    "
	} else {
		childPrefix += "│   "
	}

	// Get children
	key := strings.TrimSuffix(node.Path, "/")
	children := childrenMap[key]
	sortNodes(children)

	for i, child := range children {
		printNode(child, childPrefix, i == len(children)-1, childrenMap)
	}
}

func sortNodes(nodes []parser.Node) {
	sort.Slice(nodes, func(i, j int) bool {
		// Sort Dirs first, then Files. Both alphabetical.
		if nodes[i].Kind != nodes[j].Kind {
			return nodes[i].Kind == parser.Dir
		}
		return nodes[i].Path < nodes[j].Path
	})
}
