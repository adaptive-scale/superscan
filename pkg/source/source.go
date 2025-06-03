package source

import (
	"fmt"

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

// Source defines the interface for different source types
type Source interface {
	ListFiles(startPath string) error
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
func NewSource(sourceType SourceType) (Source, error) {
	switch sourceType {
	case GoogleDrive:
		return NewGoogleDriveSource(), nil
	case FileSystem:
		return NewFileSystemSource(), nil
	case S3Bucket:
		return nil, fmt.Errorf("S3 bucket source not implemented yet")
	case GoogleStorage:
		return nil, fmt.Errorf("Google Cloud Storage source not implemented yet")
	default:
		return nil, fmt.Errorf("unknown source type: %s", sourceType)
	}
}