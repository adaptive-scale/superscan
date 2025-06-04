package source

import (
	"fmt"
	"os"

	"github.com/adaptive-scale/superscan/pkg/config"
	"github.com/adaptive-scale/superscan/pkg/logger"
)

// SourceType represents the type of source to scan
type SourceType string

const (
	// GoogleDrive represents Google Drive source
	GoogleDrive SourceType = "google-drive"
	// FileSystem represents local filesystem source
	FileSystem SourceType = "filesystem"
	// S3Bucket represents AWS S3 bucket source
	S3Bucket SourceType = "s3"
	// GoogleStorage represents Google Cloud Storage source
	GoogleStorage SourceType = "gcs"
)

// Source defines the interface for different storage backends
type Source interface {
	ListFiles(startPath string) error
	GetName() string
}

// Set validates and sets the source type
func (st *SourceType) Set(value string) error {
	switch SourceType(value) {
	case GoogleDrive, FileSystem, S3Bucket, GoogleStorage:
		*st = SourceType(value)
		return nil
	default:
		return fmt.Errorf("invalid source type: %s", value)
	}
}

// String returns the string representation of the source type
func (st SourceType) String() string {
	return string(st)
}

// NewSource creates a new source based on the source type
func NewSource(sourceType string, cfg *config.Config) (Source, error) {
	log := logger.New(logger.INFO)
	log.Info("Creating new source of type: %s", sourceType)

	switch sourceType {
	case "google-drive":
		return NewGoogleDriveSource(), nil
	case "filesystem":
		return NewFileSystemSource(), nil
	case "s3":
		bucket := os.Getenv("AWS_S3_BUCKET")
		if bucket == "" {
			return nil, fmt.Errorf("AWS_S3_BUCKET environment variable is required for S3 source")
		}
		return NewS3Source(bucket)
	default:
		return nil, fmt.Errorf("unsupported source type: %s", sourceType)
	}
}