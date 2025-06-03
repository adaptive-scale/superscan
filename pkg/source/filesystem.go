package source

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/adaptive-scale/superscan/pkg/logger"
)

// FileNode represents a file or directory in the tree
type FileNode struct {
	Name     string
	Size     int64
	IsDir    bool
	Children []*FileNode
}

// FileSystemSource implements the Source interface for local filesystem
type FileSystemSource struct {
	rootPath string
	log      *logger.Logger
	tree     *FileNode
}

// NewFileSystemSource creates a new FileSystemSource
func NewFileSystemSource() *FileSystemSource {
	return &FileSystemSource{
		log: logger.New(logger.INFO),
	}
}

// ListFiles implements the Source interface for filesystem
func (fs *FileSystemSource) ListFiles(startPath string) error {
	fs.log.Debug("Starting filesystem scan with path: %s", startPath)

	// If no start path is provided, use current directory
	if startPath == "" {
		var err error
		startPath, err = os.Getwd()
		if err != nil {
			fs.log.Error("Failed to get current directory: %v", err)
			return fmt.Errorf("failed to get current directory: %v", err)
		}
		fs.log.Debug("Using current directory as start path: %s", startPath)
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(startPath)
	if err != nil {
		fs.log.Error("Failed to get absolute path: %v", err)
		return fmt.Errorf("failed to get absolute path: %v", err)
	}
	fs.log.Debug("Absolute path: %s", absPath)

	// Check if path exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		fs.log.Error("Path does not exist: %s", absPath)
		return fmt.Errorf("path does not exist: %s", absPath)
	}

	fs.rootPath = absPath
	fs.log.Info("Starting filesystem scan from: %s", absPath)

	// Create root node
	fs.tree = &FileNode{
		Name:  filepath.Base(absPath),
		IsDir: true,
	}

	// Use a stack to track directories to visit
	dirStack := []struct {
		path string
		node *FileNode
	}{{absPath, fs.tree}}

	// Iteratively process directories
	for len(dirStack) > 0 {
		// Pop the last directory from the stack
		current := dirStack[len(dirStack)-1]
		dirStack = dirStack[:len(dirStack)-1]

		// Read directory contents
		entries, err := os.ReadDir(current.path)
		if err != nil {
			fs.log.Error("Error reading directory %s: %v", current.path, err)
			continue
		}

		// Process each entry
		for _, entry := range entries {
			// Skip hidden files and directories
			if strings.HasPrefix(entry.Name(), ".") {
				fs.log.Debug("Skipping hidden entry: %s", entry.Name())
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

			// Create node for this entry
			node := &FileNode{
				Name:  entry.Name(),
				Size:  info.Size(),
				IsDir: entry.IsDir(),
			}

			// Add to parent's children
			current.node.Children = append(current.node.Children, node)

			// If it's a directory, add to stack
			if entry.IsDir() {
				dirStack = append(dirStack, struct {
					path string
					node *FileNode
				}{fullPath, node})
			}
		}
	}

	// Display the tree
	fs.displayTree(fs.tree, "", true)
	return nil
}

// displayTree recursively displays the file tree
func (fs *FileSystemSource) displayTree(node *FileNode, prefix string, isLast bool) {
	// Print current node
	if node == fs.tree {
		fmt.Printf("%s\n", node.Name)
	} else {
		// Determine the connector
		connector := "├── "
		if isLast {
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
	}

	// Calculate new prefix for children
	newPrefix := prefix
	if node != fs.tree {
		if isLast {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}
	}

	// Display children
	for i, child := range node.Children {
		isLastChild := i == len(node.Children)-1
		fs.displayTree(child, newPrefix, isLastChild)
	}
} 