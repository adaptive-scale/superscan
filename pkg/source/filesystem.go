package source

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adaptive-scale/superscan/pkg/logger"
)

// FileSystemSource implements Source interface for local filesystem
type FileSystemSource struct {
	log *logger.Logger
}

// NewFileSystemSource creates a new filesystem source
func NewFileSystemSource() *FileSystemSource {
	return &FileSystemSource{
		log: logger.New(logger.INFO),
	}
}

// ListFiles lists files in the filesystem
func (fs *FileSystemSource) ListFiles(startPath string) error {
	fs.log.Info("Starting filesystem scan from path: %s", startPath)

	// Get current directory if startPath is empty
	if startPath == "" {
		var err error
		startPath, err = os.Getwd()
		if err != nil {
			fs.log.Error("Failed to get current directory: %v", err)
			return fmt.Errorf("failed to get current directory: %v", err)
		}
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(startPath)
	if err != nil {
		fs.log.Error("Failed to convert path to absolute: %v", err)
		return fmt.Errorf("failed to convert path to absolute: %v", err)
	}

	// Check if path exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		fs.log.Error("Path does not exist: %s", absPath)
		return fmt.Errorf("path does not exist: %s", absPath)
	}

	// Create root node
	root := &FileNode{
		Name:     filepath.Base(absPath),
		IsDir:    true,
		Children: make([]*FileNode, 0),
	}

	// Create a stack for iterative traversal
	stack := []struct {
		path   string
		parent *FileNode
	}{
		{absPath, root},
	}

	// Process directories iteratively
	for len(stack) > 0 {
		// Pop from stack
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// Read directory
		entries, err := os.ReadDir(current.path)
		if err != nil {
			fs.log.Error("Failed to read directory %s: %v", current.path, err)
			continue
		}

		// Process each entry
		for _, entry := range entries {
			// Skip hidden files and directories
			if entry.Name()[0] == '.' {
				continue
			}

			// Get full path
			fullPath := filepath.Join(current.path, entry.Name())

			// Get file info
			info, err := entry.Info()
			if err != nil {
				fs.log.Error("Failed to get file info for %s: %v", fullPath, err)
				continue
			}

			// Create node
			node := &FileNode{
				Name:  entry.Name(),
				IsDir: entry.IsDir(),
				Size:  info.Size(),
			}

			// Add to parent's children
			current.parent.Children = append(current.parent.Children, node)

			// If directory, add to stack
			if entry.IsDir() {
				node.Children = make([]*FileNode, 0)
				stack = append(stack, struct {
					path   string
					parent *FileNode
				}{fullPath, node})
			}
		}
	}

	// Display the tree
	fsdisplayTree(root, 0)
	return nil
}

// GetName returns the source name
func (fs *FileSystemSource) GetName() string {
	return "filesystem"
}

// fsdisplayTree recursively displays the file tree
func fsdisplayTree(node *FileNode, level int) {
	// Print current node
	prefix := ""
	for i := 0; i < level; i++ {
		prefix += "│   "
	}

	// Determine the connector
	connector := "├── "
	if len(node.Children) == 0 {
		connector = "└── "
	}

	// Print the node with appropriate icon
	icon := "📄"
	if node.IsDir {
		icon = "📁"
	}

	if node.IsDir {
		fmt.Printf("%s%s%s%s/\n", prefix, connector, icon, node.Name)
	} else {
		fmt.Printf("%s%s%s%s (%d bytes)\n", prefix, connector, icon, node.Name, node.Size)
	}

	// Display children
	for _, child := range node.Children {
		fsdisplayTree(child, level+1)
	}
} 
