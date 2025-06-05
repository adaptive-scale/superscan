package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/adaptive-scale/superscan/pkg/config"
	"github.com/adaptive-scale/superscan/pkg/logger"
	"github.com/adaptive-scale/superscan/pkg/source"
)

func main() {
	// Parse command line flags
	sourceType := flag.String("source", "", "Source type (google-drive, filesystem, s3)")
	path := flag.String("path", "", "Path for listing or downloading")
	destination := flag.String("destination", "", "Destination path for downloaded file (required for download)")
	recursive := flag.Bool("recursive", false, "Download entire directory tree (only works with directory paths)")
	sampleSize := flag.Int("sample", 0, "Number of files to download as a sample (0 means download all files)")
	flag.Parse()

	// Initialize logger
	log := logger.New(logger.INFO)

	// Validate source type
	if *sourceType == "" {
		log.Error("Source type is required")
		fmt.Println("Usage: superscan --source <source-type> [--path <path>] [--destination <dest-path>] [--recursive] [--sample <number>]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Error("Failed to load configuration: %v", err)
		os.Exit(1)
	}

	// Create source
	src, err := source.NewSource(*sourceType, cfg)
	if err != nil {
		log.Error("Failed to create source: %v", err)
		os.Exit(1)
	}

	// Handle download if destination is specified
	if *destination != "" {
		if *path == "" {
			log.Error("Path is required for download")
			fmt.Println("Usage: superscan --source <source-type> --path <file-path> --destination <dest-path> [--recursive] [--sample <number>]")
			os.Exit(1)
		}

		// Handle destination path
		destPath := *destination
		if !strings.HasSuffix(destPath, "/") && !strings.HasSuffix(destPath, string(os.PathSeparator)) {
			destPath = destPath + string(os.PathSeparator)
		}

		// Create destination directory if it doesn't exist
		if err := os.MkdirAll(destPath, 0755); err != nil {
			log.Error("Failed to create destination directory: %v", err)
			os.Exit(1)
		}

		if *recursive {
			// Download entire directory tree
			log.Info("Downloading directory tree from %s to %s", *path, destPath)
			if *sampleSize > 0 {
				log.Info("Downloading sample of %d files", *sampleSize)
			}
			downloadDirectoryTree(src, *path, destPath, *sampleSize, log)
			log.Info("Directory tree download completed successfully")
		} else {
			// Single file download
			if strings.HasSuffix(*path, "/") || strings.HasSuffix(*path, string(os.PathSeparator)) {
				log.Error("Path is a directory. Use --recursive flag to download directories")
				fmt.Printf("\nTips:\n")
				fmt.Printf("1. Use --recursive flag to download entire directory\n")
				fmt.Printf("2. Or specify a file path instead of a directory\n")
				os.Exit(1)
			}

			// Use source filename if destination is a directory
			filename := filepath.Base(*path)
			finalDestPath := filepath.Join(destPath, filename)

			log.Info("Downloading file from %s to %s", *path, finalDestPath)
			if err := src.DownloadFile(*path, finalDestPath); err != nil {
				if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "NoSuchKey") {
					log.Error("File not found: %s", *path)
					fmt.Printf("\nTips:\n")
					fmt.Printf("1. Verify the file path is correct\n")
					fmt.Printf("2. Use --path without --destination to list available files\n")
					fmt.Printf("3. Check if you have proper permissions to access the file\n")
				} else {
					log.Error("Failed to download file: %v", err)
				}
				os.Exit(1)
			}
			log.Info("Download completed successfully")
		}
		return
	}

	// List files
	if err := src.ListFiles(*path); err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "NoSuchKey") {
			log.Error("Path not found: %s", *path)
			fmt.Printf("\nTips:\n")
			fmt.Printf("1. Verify the path is correct\n")
			fmt.Printf("2. Try listing the root directory first (omit --path)\n")
			fmt.Printf("3. Check if you have proper permissions to access the path\n")
		} else {
			log.Error("Failed to list files: %v", err)
		}
		os.Exit(1)
	}
}

// downloadDirectoryTree downloads an entire directory tree while maintaining the structure
func downloadDirectoryTree(src source.Source, sourcePath, destPath string, sampleSize int, log *logger.Logger) {
	// Get the file tree from the source first
	tree, err := src.GetFileTree(sourcePath)
	if err != nil {
		log.Error("Failed to get file tree: %v", err)
		return
	}

	// Display the tree structure
	log.Info("Directory structure to be downloaded:")
	displayTree(tree, 0)
	fmt.Println() // Add a blank line for better readability

	// First pass: Create all directory structures
	log.Info("Creating directory structure...")
	createDirectories(tree, destPath, log)

	// Second pass: Download files
	log.Info("Starting file downloads...")
	if sampleSize > 0 {
		// Initialize random seed
		rand.Seed(time.Now().UnixNano())
		downloadSampleFiles(tree, "", destPath, src, log, sampleSize)
	} else {
		downloadFiles(tree, "", destPath, src, log)
	}
}

// downloadSampleFiles downloads a random sample of files from the tree
func downloadSampleFiles(node *source.FileNode, sourceBase, destBase string, src source.Source, log *logger.Logger, sampleSize int) error {
	// Collect all files in the tree
	var allFiles []struct {
		sourcePath string
		destPath   string
	}
	collectFiles(node, sourceBase, destBase, &allFiles)

	// If we have fewer files than the sample size, download all
	if len(allFiles) <= sampleSize {
		log.Info("Found %d files, downloading all", len(allFiles))
		for _, file := range allFiles {
			log.Info("Downloading: %s", file.sourcePath)
			if err := src.DownloadFile(file.sourcePath, file.destPath); err != nil {
				log.Error("Failed to download %s: %v", file.sourcePath, err)
				continue
			}
			log.Info("Successfully downloaded: %s", file.sourcePath)
		}
		return nil
	}

	// Create a random number generator with current time as seed
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Create a map to track selected indices
	selected := make(map[int]bool)
	selectedCount := 0

	// Randomly select files
	log.Info("Found %d files, randomly selecting %d files", len(allFiles), sampleSize)
	for selectedCount < sampleSize {
		// Generate a random index
		idx := rng.Intn(len(allFiles))
		
		// Skip if already selected
		if selected[idx] {
			continue
		}

		// Mark as selected
		selected[idx] = true
		selectedCount++

		// Download the file
		file := allFiles[idx]
		log.Info("Downloading random file %d/%d: %s", selectedCount, sampleSize, file.sourcePath)
		if err := src.DownloadFile(file.sourcePath, file.destPath); err != nil {
			log.Error("Failed to download %s: %v", file.sourcePath, err)
			continue
		}
		log.Info("Successfully downloaded: %s", file.sourcePath)
	}

	log.Info("Completed downloading %d random files", selectedCount)
	return nil
}

// collectFiles recursively collects all files in the tree
func collectFiles(node *source.FileNode, sourceBase, destBase string, files *[]struct{ sourcePath, destPath string }) {
	if !node.IsDir {
		sourcePath := node.Name
		if sourceBase != "" {
			sourcePath = filepath.Join(sourceBase, node.Name)
		}
		sourcePath = filepath.Clean(sourcePath)
		sourcePath = strings.ReplaceAll(sourcePath, "\\", "/")
		destPath := filepath.Join(destBase, node.Name)
		*files = append(*files, struct{ sourcePath, destPath string }{sourcePath, destPath})
	}

	for _, child := range node.Children {
		childSourceBase := sourceBase
		childDestBase := destBase
		if node.Name != "" {
			childSourceBase = filepath.Join(sourceBase, node.Name)
			childDestBase = filepath.Join(destBase, node.Name)
		}
		collectFiles(child, childSourceBase, childDestBase, files)
	}
}

// createDirectories recursively creates directory structure
func createDirectories(node *source.FileNode, destBase string, log *logger.Logger) error {
	destPath := filepath.Join(destBase, node.Name)
	if node.IsDir {
		if err := os.MkdirAll(destPath, 0755); err != nil {
			log.Error("Failed to create directory %s: %v", destPath, err)
			return err
		}
		log.Info("Created directory: %s", destPath)
	}

	// Process children
	for _, child := range node.Children {
		if err := createDirectories(child, destPath, log); err != nil {
			return err
		}
	}
	return nil
}

// downloadFiles recursively downloads files
func downloadFiles(node *source.FileNode, sourceBase, destBase string, src source.Source, log *logger.Logger) error {
	if !node.IsDir {
		// For files, use the current node's name
		sourcePath := node.Name
		if sourceBase != "" {
			sourcePath = filepath.Join(sourceBase, node.Name)
		}
		sourcePath = filepath.Clean(sourcePath)
		sourcePath = strings.ReplaceAll(sourcePath, "\\", "/") // Ensure forward slashes for S3
		
		destPath := filepath.Join(destBase, node.Name)
		log.Info("Downloading: %s", sourcePath)
		if err := src.DownloadFile(sourcePath, destPath); err != nil {
			log.Error("Failed to download %s: %v", sourcePath, err)
			return err
		}
		log.Info("Successfully downloaded: %s", sourcePath)
	}

	// Process children
	for _, child := range node.Children {
		// For children, use the current node's path as the base
		childSourceBase := sourceBase
		childDestBase := destBase
		
		// Only append the current node's name if it's not the root
		if node.Name != "" {
			childSourceBase = filepath.Join(sourceBase, node.Name)
			childDestBase = filepath.Join(destBase, node.Name)
		}
		
		if err := downloadFiles(child, childSourceBase, childDestBase, src, log); err != nil {
			return err
		}
	}
	return nil
}

// displayTree recursively displays the tree structure with file sizes
func displayTree(node *source.FileNode, level int) {
	// Create indentation
	indent := strings.Repeat("  ", level)
	
	// Display node name with size if it's a file
	if node.IsDir {
		fmt.Printf("%s📁 %s/\n", indent, node.Name)
	} else {
		size := formatSize(node.Size)
		fmt.Printf("%s📄 %s (%s)\n", indent, node.Name, size)
	}

	// Process children
	for _, child := range node.Children {
		displayTree(child, level+1)
	}
}

// formatSize formats the file size in a human-readable format
func formatSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
} 