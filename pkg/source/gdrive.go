package source

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/adaptive-scale/superscan/pkg/config"
	"github.com/adaptive-scale/superscan/pkg/logger"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

// GoogleDriveSource implements the Source interface for Google Drive
type GoogleDriveSource struct {
	service *drive.Service
	log     *logger.Logger
}

// NewGoogleDriveSource creates a new GoogleDriveSource
func NewGoogleDriveSource(cfg *config.Config) (Source, error) {
	return &GoogleDriveSource{
		log: logger.New(logger.INFO),
	}, nil
}

func (gds *GoogleDriveSource) GetName() string {
//TODO
	return "Google Drive"
}

// GetDescription returns the source description
func (gds *GoogleDriveSource) GetDescription() string {
	return "Google Drive Storage"
}

// ListFiles implements the Source interface for Google Drive
func (gds *GoogleDriveSource) ListFiles(startPath string) error {
	gds.log.Debug("Starting Google Drive scan with path: %s", startPath)
	ctx := context.Background()

	// Get credentials file path from environment or use default
	credentialsFile := os.Getenv("SUPERSCAN_CONFIG_GOOGLE")
	if credentialsFile == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			gds.log.Error("Failed to get home directory: %v", err)
			return fmt.Errorf("failed to get home directory: %v", err)
		}
		credentialsFile = fmt.Sprintf("%s/.superscan/credentials.json", homeDir)
		gds.log.Debug("Using default credentials file: %s", credentialsFile)
	} else {
		gds.log.Debug("Using credentials file from environment: %s", credentialsFile)
	}

	// Read credentials file
	b, err := os.ReadFile(credentialsFile)
	if err != nil {
		gds.log.Error("Unable to read credentials file: %v", err)
		return fmt.Errorf("unable to read credentials file: %v", err)
	}
	gds.log.Debug("Successfully read credentials file")

	// Configure OAuth2
	config, err := google.ConfigFromJSON(b, drive.DriveReadonlyScope)
	if err != nil {
		gds.log.Error("Unable to parse client secret file: %v", err)
		return fmt.Errorf("unable to parse client secret file to config: %v", err)
	}
	gds.log.Debug("Successfully configured OAuth2")

	// Get token file path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		gds.log.Error("Failed to get home directory: %v", err)
		return fmt.Errorf("failed to get home directory: %v", err)
	}
	tokenFile := fmt.Sprintf("%s/.superscan/token.json", homeDir)
	gds.log.Debug("Using token file: %s", tokenFile)

	// Get token from file or create new one
	tok, err := getTokenFromFile(tokenFile, config)
	if err != nil {
		gds.log.Info("Token not found in file, requesting new token")
		tok, err = getTokenFromWeb(config)
		if err != nil {
			gds.log.Error("Unable to get token: %v", err)
			return fmt.Errorf("unable to get token: %v", err)
		}
		saveToken(tokenFile, tok)
		gds.log.Info("New token saved to file")
	} else {
		gds.log.Debug("Successfully loaded token from file")
	}

	// Create Drive service
	service, err := drive.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, tok)))
	if err != nil {
		gds.log.Error("Unable to retrieve Drive client: %v", err)
		return fmt.Errorf("unable to retrieve Drive client: %v", err)
	}
	gds.log.Debug("Successfully created Drive service")

	gds.service = service

	// If no start path is provided, use root
	if startPath == "" {
		startPath = "root"
		gds.log.Debug("Using root as start path")
	}

	// List files
	gds.log.Info("Starting Google Drive scan from: %s", startPath)
	return gds.listFiles(startPath)
}

// listFiles lists files and folders in Google Drive
func (gds *GoogleDriveSource) listFiles(folderId string) error {
	gds.log.Debug("Listing files in folder: %s", folderId)
	query := fmt.Sprintf("'%s' in parents and trashed = false", folderId)
	r, err := gds.service.Files.List().
		Q(query).
		Fields("files(id, name, mimeType, size)").
		Do()
	if err != nil {
		gds.log.Error("Unable to retrieve files: %v", err)
		return fmt.Errorf("unable to retrieve files: %v", err)
	}

	for _, file := range r.Files {
		if file.MimeType == "application/vnd.google-apps.folder" {
			gds.log.Info("Found directory: %s/", file.Name)
			fmt.Printf("📁 %s/\n", file.Name)
			// Recursively list files in subfolder
			if err := gds.listFiles(file.Id); err != nil {
				return err
			}
		} else {
			gds.log.Info("Found file: %s (%d bytes)", file.Name, file.Size)
			fmt.Printf("📄 %s (%d bytes)\n", file.Name, file.Size)
		}
	}

	return nil
}

// getTokenFromFile retrieves a token from a local file
func getTokenFromFile(file string, config *oauth2.Config) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// getTokenFromWeb requests a token from the web
func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, fmt.Errorf("unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve token from web: %v", err)
	}
	return tok, nil
}

// saveToken saves a token to a file
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// DownloadFile downloads a file from Google Drive to the local filesystem
func (g *GoogleDriveSource) DownloadFile(filePath string, destination string) error {
	g.log.Info("Downloading file from Google Drive: %s to %s", filePath, destination)

	// Create destination directory if it doesn't exist
	destDir := filepath.Dir(destination)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		g.log.Error("Failed to create destination directory: %v", err)
		return fmt.Errorf("failed to create destination directory: %v", err)
	}

	// Create destination file
	destFile, err := os.Create(destination)
	if err != nil {
		g.log.Error("Failed to create destination file: %v", err)
		return fmt.Errorf("failed to create destination file: %v", err)
	}
	defer destFile.Close()

	// Download file content
	resp, err := g.service.Files.Get(filePath).Download()
	if err != nil {
		g.log.Error("Failed to download file: %v", err)
		return fmt.Errorf("failed to download file: %v", err)
	}
	defer resp.Body.Close()

	// Copy file contents
	_, err = io.Copy(destFile, resp.Body)
	if err != nil {
		g.log.Error("Failed to copy file contents: %v", err)
		return fmt.Errorf("failed to copy file contents: %v", err)
	}

	g.log.Info("Successfully downloaded file to %s", destination)
	return nil
}

// GetFileTree returns the file tree structure for the given path
func (s *GoogleDriveSource) GetFileTree(path string) (*FileNode, error) {
	// Create root node
	root := &FileNode{
		Name:     filepath.Base(path),
		IsDir:    true,
		Children: make([]*FileNode, 0),
	}

	// Get the root folder ID
	folderID, err := s.getFolderID(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get folder ID: %v", err)
	}

	// Use a stack for iterative traversal
	stack := []struct {
		node     *FileNode
		folderID string
	}{
		{root, folderID},
	}

	for len(stack) > 0 {
		// Pop from stack
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// List files in current folder
		files, err := s.service.Files.List().
			Q(fmt.Sprintf("'%s' in parents and trashed = false", current.folderID)).
			Fields("files(id, name, mimeType, size)").
			Do()
		if err != nil {
			return nil, fmt.Errorf("failed to list files: %v", err)
		}

		// Process each file
		for _, file := range files.Files {
			isDir := file.MimeType == "application/vnd.google-apps.folder"
			node := &FileNode{
				Name:     file.Name,
				IsDir:    isDir,
				Size:     file.Size,
				Children: make([]*FileNode, 0),
			}

			current.node.Children = append(current.node.Children, node)

			// If it's a directory, add it to the stack
			if isDir {
				stack = append(stack, struct {
					node     *FileNode
					folderID string
				}{node, file.Id})
			}
		}
	}

	return root, nil
}

// getFolderID returns the folder ID for the given path
func (gds *GoogleDriveSource) getFolderID(path string) (string, error) {
	// If path is empty or "root", return root folder ID
	if path == "" || path == "root" {
		return "root", nil
	}

	// Split path into components
	parts := strings.Split(path, "/")
	currentID := "root"

	// Traverse the path
	for _, part := range parts {
		if part == "" {
			continue
		}

		// Search for the folder
		query := fmt.Sprintf("name = '%s' and mimeType = 'application/vnd.google-apps.folder' and '%s' in parents and trashed = false", part, currentID)
		files, err := gds.service.Files.List().
			Q(query).
			Fields("files(id, name)").
			Do()
		if err != nil {
			return "", fmt.Errorf("failed to search for folder %s: %v", part, err)
		}

		if len(files.Files) == 0 {
			return "", fmt.Errorf("folder not found: %s", part)
		}

		currentID = files.Files[0].Id
	}

	return currentID, nil
} 