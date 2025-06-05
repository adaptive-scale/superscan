package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	GoogleDrive struct {
		CredentialsFile string `yaml:"credentials_file"`
		TokenFile       string `yaml:"token_file"`
	} `yaml:"google_drive"`

	S3 struct {
		Bucket    string `yaml:"bucket"`
		Region    string `yaml:"region"`
		AccessKey string `yaml:"access_key"`
		SecretKey string `yaml:"secret_key"`
	} `yaml:"s3"`
}

// Load loads the configuration from the config file
func Load() (*Config, error) {
	// Get home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %v", err)
	}

	// Construct config file path
	configDir := filepath.Join(homeDir, ".superscan")
	configFile := filepath.Join(configDir, "config.yaml")

	// Create default config
	config := &Config{}

	// Check if config file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// Create config directory if it doesn't exist
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create config directory: %v", err)
		}

		// Create default config file
		if err := createDefaultConfig(configFile); err != nil {
			return nil, fmt.Errorf("failed to create default config: %v", err)
		}
	} else {
		// Read existing config file
		data, err := os.ReadFile(configFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %v", err)
		}

		// Parse YAML
		if err := yaml.Unmarshal(data, config); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %v", err)
		}
	}

	// Override with environment variables if set
	if bucket := os.Getenv("AWS_S3_BUCKET"); bucket != "" {
		config.S3.Bucket = bucket
	}
	if region := os.Getenv("AWS_REGION"); region != "" {
		config.S3.Region = region
	}
	if accessKey := os.Getenv("AWS_ACCESS_KEY_ID"); accessKey != "" {
		config.S3.AccessKey = accessKey
	}
	if secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY"); secretKey != "" {
		config.S3.SecretKey = secretKey
	}
	if credsFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"); credsFile != "" {
		config.GoogleDrive.CredentialsFile = credsFile
	}

	return config, nil
}

// createDefaultConfig creates a default configuration file
func createDefaultConfig(configFile string) error {
	config := &Config{
		GoogleDrive: struct {
			CredentialsFile string `yaml:"credentials_file"`
			TokenFile       string `yaml:"token_file"`
		}{
			CredentialsFile: filepath.Join(filepath.Dir(configFile), "credentials.json"),
			TokenFile:       filepath.Join(filepath.Dir(configFile), "token.json"),
		},
		S3: struct {
			Bucket    string `yaml:"bucket"`
			Region    string `yaml:"region"`
			AccessKey string `yaml:"access_key"`
			SecretKey string `yaml:"secret_key"`
		}{
			Region: "us-east-1",
		},
	}

	// Marshal to YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal default config: %v", err)
	}

	// Write to file
	if err := os.WriteFile(configFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write default config: %v", err)
	}

	return nil
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