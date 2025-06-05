package source

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/adaptive-scale/superscan/pkg/logger"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Source implements Source interface for AWS S3
type S3Source struct {
	client *s3.Client
	bucket string
	log    *logger.Logger
}

// NewS3Source creates a new S3 source
func NewS3Source(bucket string) (*S3Source, error) {
	log := logger.New(logger.INFO)
	log.Info("Initializing S3 source for bucket: %s", bucket)

	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Error("Failed to load AWS config: %v", err)
		return nil, fmt.Errorf("failed to load AWS config: %v", err)
	}

	// Create S3 client
	client := s3.NewFromConfig(cfg)

	return &S3Source{
		client: client,
		bucket: bucket,
		log:    log,
	}, nil
}

// ListFiles lists files in the S3 bucket
func (s *S3Source) ListFiles(startPath string) error {
	s.log.Info("Starting S3 scan from path: %s", startPath)

	// Ensure startPath doesn't start with /
	startPath = strings.TrimPrefix(startPath, "/")

	// Create root node for tree
	root := &FileNode{
		Name:     s.bucket,
		IsDir:    true,
		Children: make([]*FileNode, 0),
	}

	// Initialize paginator for listing objects
	paginator := s3.NewListObjectsV2Paginator(s.client, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(startPath),
	})

	// Process each page of results
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			s.log.Error("Failed to list objects: %v", err)
			return fmt.Errorf("failed to list objects: %v", err)
		}

		// Process each object
		for _, obj := range page.Contents {
			// Skip the start path itself
			if *obj.Key == startPath {
				continue
			}

			// Get relative path from startPath
			relPath := strings.TrimPrefix(*obj.Key, startPath)
			relPath = strings.TrimPrefix(relPath, "/")

			// Split path into components
			parts := strings.Split(relPath, "/")

			// Start from root
			current := root

			// Create directory structure
			for i, part := range parts {
				isLast := i == len(parts)-1
				isDir := !isLast || strings.HasSuffix(*obj.Key, "/")

				if isDir {
					// Find or create directory node
					var dirNode *FileNode
					for _, child := range current.Children {
						if child.Name == part && child.IsDir {
							dirNode = child
							break
						}
					}

					if dirNode == nil {
						dirNode = &FileNode{
							Name:     part,
							IsDir:    true,
							Children: make([]*FileNode, 0),
						}
						current.Children = append(current.Children, dirNode)
					}
					current = dirNode
				} else {
					// Create file node
					fileNode := &FileNode{
						Name:  part,
						IsDir: false,
						Size:  *obj.Size,
					}
					current.Children = append(current.Children, fileNode)
				}
			}
		}
	}

	// Display the tree
	displayTree(root, 0)
	return nil
}

// DownloadFile downloads a file from S3
func (s *S3Source) DownloadFile(filePath, destination string) error {
	// Clean and normalize the file path
	filePath = strings.TrimPrefix(filePath, "/")
	filePath = filepath.Clean(filePath)
	filePath = strings.ReplaceAll(filePath, "\\", "/") // Ensure forward slashes for S3

	// Check if the path is a directory
	if strings.HasSuffix(filePath, "/") || filePath == "" {
		// For directories, just create the directory structure
		if err := os.MkdirAll(destination, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}
		s.log.Info("Created directory structure: %s", destination)
		return nil
	}

	// For files, ensure the destination directory exists
	destDir := filepath.Dir(destination)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %v", err)
	}

	// Get the object from S3
	getResult, err := s.client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filePath),
	})
	if err != nil {
		return fmt.Errorf("failed to get object from S3: %v", err)
	}
	defer getResult.Body.Close()

	// Create the destination file
	file, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %v", err)
	}
	defer file.Close()

	// Copy the file contents
	if _, err := io.Copy(file, getResult.Body); err != nil {
		return fmt.Errorf("failed to copy file contents: %v", err)
	}

	s.log.Info("Successfully downloaded file: %s", filePath)
	return nil
}

// GetName returns the source name
func (s *S3Source) GetName() string {
	return "s3"
}

// GetDescription returns the source description
func (s *S3Source) GetDescription() string {
	return "AWS S3 Storage"
}

// GetFileTree returns the file tree structure for the given path
func (s *S3Source) GetFileTree(path string) (*FileNode, error) {
	// Clean and normalize the path
	path = strings.TrimPrefix(path, "/")
	path = filepath.Clean(path)
	path = strings.ReplaceAll(path, "\\", "/")

	// Create root node
	root := &FileNode{
		Name:     filepath.Base(path),
		IsDir:    true,
		Children: make([]*FileNode, 0),
	}

	// List objects with prefix
	prefix := path
	if !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}

	// List objects
	paginator := s3.NewListObjectsV2Paginator(s.client, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(prefix),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return nil, fmt.Errorf("failed to list objects: %v", err)
		}

		for _, obj := range page.Contents {
			// Skip the prefix itself
			if *obj.Key == prefix {
				continue
			}

			// Get relative path
			relPath := strings.TrimPrefix(*obj.Key, prefix)
			if relPath == "" {
				continue
			}

			// Split path into components
			parts := strings.Split(relPath, "/")
			currentNode := root

			// Create directory nodes if needed
			for i := 0; i < len(parts)-1; i++ {
				// Find or create directory node
				var dirNode *FileNode
				for _, child := range currentNode.Children {
					if child.Name == parts[i] && child.IsDir {
						dirNode = child
						break
					}
				}

				if dirNode == nil {
					dirNode = &FileNode{
						Name:     parts[i],
						IsDir:    true,
						Children: make([]*FileNode, 0),
					}
					currentNode.Children = append(currentNode.Children, dirNode)
				}
				currentNode = dirNode
			}

			// Add file node
			fileNode := &FileNode{
				Name:     parts[len(parts)-1],
				IsDir:    false,
				Size:     *obj.Size,
				Children: nil,
			}
			currentNode.Children = append(currentNode.Children, fileNode)
		}
	}

	return root, nil
} 