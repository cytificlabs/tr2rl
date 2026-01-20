package parser

type NodeKind string

const (
	Dir  NodeKind = "dir"
	File NodeKind = "file"
)

type Node struct {
	Path string   // normalized relative path like "core/model/entity.py"
	Kind NodeKind // dir or file
}

type Result struct {
	Nodes        []Node
	Normalized   string
	Warnings     []string
	RootInferred bool // Did we guess the root directory?
}
