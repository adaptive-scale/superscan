# SuperScan

SuperScan is a versatile file system scanner that supports multiple storage backends including Google Drive, local filesystem, AWS S3, and Google Cloud Storage. It provides a unified interface to list and visualize files and directories across different storage systems.

## Features

- Multiple storage backend support:
  - Google Drive
  - Local filesystem
  - AWS S3 (coming soon)
  - Google Cloud Storage (coming soon)
- Directory tree visualization
- Configurable starting paths
- Environment variable support
- YAML configuration

## Installation

### Prerequisites

- Go 1.16 or higher
- Google Cloud credentials (for Google Drive support)

### Building from Source

```bash
# Clone the repository
git clone https://github.com/adaptive-scale/superscan.git
cd superscan

# Build the project
make build

# The binary will be available in the bin directory
```

## Usage

### List Files from Google Drive

```bash
# List files from Google Drive root
./bin/superscan --source-type google-drive

# List files from a specific Google Drive folder
./bin/superscan --source-type google-drive --start-path "My Drive/Folder"

# Show only directory tree
./bin/superscan --source-type google-drive --only-tree
```

### List Files from Local Filesystem

```bash
# List files from current directory
./bin/superscan --source-type filesystem

# List files from a specific directory
./bin/superscan --source-type filesystem --start-path /path/to/directory

# Show only directory tree
./bin/superscan --source-type filesystem --only-tree
```

## Command Line Options

- `--source-type`: Type of storage to scan (google-drive, filesystem, s3, gcs)
- `--start-path`: Starting path for scanning (default: root for Google Drive, current directory for filesystem)
- `--only-tree`: Show only directory tree without file details
- `--version`: Show version information
- `--config`: Path to configuration file (default: ~/.superscan/config.yaml)

## Configuration

The configuration file is stored in YAML format at `~/.superscan/config.yaml` by default. You can specify a different location using the `--config` flag.

Example configuration:

```yaml
google_drive:
  credentials_file: /path/to/credentials.json
  token_file: /path/to/token.json
  start_path: root
```

### Environment Variables

- `SUPERSCAN_CONFIG_GOOGLE`: Path to Google Drive credentials file

## Google Drive Setup

1. Go to the [Google Cloud Console](https://console.cloud.google.com)
2. Create a new project or select an existing one
3. Enable the Google Drive API
4. Create credentials (OAuth 2.0 Client ID)
5. Download the credentials and save them as `credentials.json`
6. Place the credentials file in your config directory or specify its location in the config file

## Development

### Building

```bash
# Build the project
make build

# Build for specific platforms
make build-linux
make build-darwin
make build-windows
```

### Testing

```bash
# Run tests
make test

# Run tests with coverage
make test-coverage
```

### Linting

```bash
# Run linters
make lint
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request 