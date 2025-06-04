package source

import (
	"fmt"
	"strings"
)

// FileNode represents a file or directory in the tree
type FileNode struct {
	Name     string
	IsDir    bool
	Size     int64
	Children []*FileNode
}

// displayTree displays the file tree in ASCII format
func displayTree(node *FileNode, level int) {
	// Print indentation
	indent := strings.Repeat("  ", level)

	// Print node
	if node.IsDir {
		fmt.Printf("%s📁 %s/\n", indent, node.Name)
	} else {
		fmt.Printf("%s📄 %s (%d bytes)\n", indent, node.Name, node.Size)
	}

	// Print children
	for _, child := range node.Children {
		displayTree(child, level+1)
	}
}
