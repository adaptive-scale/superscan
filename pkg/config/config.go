package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds the application configuration
type Config struct {
	GoogleDrive GoogleDriveConfig `yaml:"google_drive,omitempty"`
	S3         S3Config          `yaml:"s3,omitempty"`
}

// GoogleDriveConfig holds Google Drive specific configuration
type GoogleDriveConfig struct {
	CredentialsFile string `yaml:"credentials_file"`
	TokenFile      string `yaml:"token_file"`
	StartPath      string `yaml:"start_path"`
}

// S3Config holds AWS S3 specific configuration
type S3Config struct {
	Bucket    string `yaml:"bucket"`
	Region    string `yaml:"region"`
	StartPath string `yaml:"start_path"`
}

// LoadConfig loads the configuration from a file or environment variable
func LoadConfig(configPath string) (*Config, error) {
	// If no config path is provided, use default
	if configPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		configPath = filepath.Join(homeDir, ".superscan", "config.yaml")
	}

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, err
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Create default config if it doesn't exist
			config := &Config{
				GoogleDrive: GoogleDriveConfig{
					CredentialsFile: filepath.Join(configDir, "credentials.json"),
					TokenFile:      filepath.Join(configDir, "token.json"),
					StartPath:      "root",
				},
				S3: S3Config{
					Bucket:    "",
					Region:    "us-east-1",
					StartPath: "",
				},
			}
			if err := SaveConfig(configPath, config); err != nil {
				return nil, err
			}
			return config, nil
		}
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}

	// Override credentials file path if environment variable is set
	if credsPath := os.Getenv("SUPERSCAN_CONFIG_GOOGLE"); credsPath != "" {
		config.GoogleDrive.CredentialsFile = credsPath
	}

	// Override S3 bucket if environment variable is set
	if bucket := os.Getenv("AWS_S3_BUCKET"); bucket != "" {
		config.S3.Bucket = bucket
	}

	// Override S3 region if environment variable is set
	if region := os.Getenv("AWS_REGION"); region != "" {
		config.S3.Region = region
	}

	return &config, nil
}

// SaveConfig saves the configuration to a file
func SaveConfig(configPath string, config *Config) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("error marshaling config: %v", err)
	}

	return os.WriteFile(configPath, data, 0644)
}

// GetConfigAsYAML returns the configuration as a YAML string
func GetConfigAsYAML(config *Config) (string, error) {
	data, err := yaml.Marshal(config)
	if err != nil {
		return "", fmt.Errorf("error marshaling config: %v", err)
	}
	return string(data), nil
} 