# SuperScan

SuperScan is a powerful file scanning and management tool that supports multiple storage backends including Google Drive, local filesystem, and AWS S3. It provides a beautiful ASCII tree visualization of your files and directories, making it easy to navigate and manage your files across different storage systems.

## Features

- Multiple storage backend support (Google Drive, Local Filesystem, AWS S3)
- Beautiful ASCII tree visualization
- Configurable starting paths
- Environment variable support
- YAML configuration
- Comprehensive logging system
- Memory-efficient iterative scanning
- File download functionality for all storage backends

## Installation

### Prerequisites

- Go 1.16 or later
- Google Cloud Platform account (for Google Drive)
- AWS account (for S3)

### Building from Source

```bash
# Clone the repository
git clone https://github.com/adaptive-scale/superscan.git
cd superscan

# Build using the build script
./build.sh

# Or build manually
go build -o bin/superscan cmd/superscan/main.go
```

## Usage

### Listing Files

#### Google Drive

```bash
# List files from Google Drive
./bin/superscan --source google-drive

# List files from a specific folder
./bin/superscan --source google-drive --path "folder-id"
```

#### Local Filesystem

```bash
# List files from current directory
./bin/superscan --source filesystem

# List files from a specific path
./bin/superscan --source filesystem --path "/path/to/directory"
```

#### AWS S3

```bash
# List files from S3 bucket
./bin/superscan --source s3 --path "prefix/"

# List files from a specific prefix
./bin/superscan --source s3 --path "folder/subfolder/"
```

### Downloading Files

#### Google Drive

```bash
# Download a file from Google Drive
./bin/superscan --source google-drive --download "file-id" --destination "/path/to/save/file"
```

#### Local Filesystem

```bash
# Download (copy) a file from local filesystem
./bin/superscan --source filesystem --download "/path/to/file" --destination "/path/to/save/file"
```

#### AWS S3

```bash
# Download a file from S3
./bin/superscan --source s3 --download "path/to/file" --destination "/path/to/save/file"
```

## Configuration

### Environment Variables

- `GOOGLE_APPLICATION_CREDENTIALS`: Path to Google Cloud credentials file
- `AWS_S3_BUCKET`: S3 bucket name
- `AWS_ACCESS_KEY_ID`: AWS access key
- `AWS_SECRET_ACCESS_KEY`: AWS secret key
- `AWS_REGION`: AWS region

### YAML Configuration

Create a `config.yaml` file in your home directory under `.superscan/`:

```yaml
google_drive:
  credentials_file: /path/to/credentials.json
  token_file: /path/to/token.json

s3:
  bucket: your-bucket-name
  region: your-region
  access_key: your-access-key
  secret_key: your-secret-key
```

## Project Structure

```
.
├── cmd/
│   └── superscan/
│       └── main.go
├── pkg/
│   ├── config/
│   │   └── config.go
│   ├── logger/
│   │   └── logger.go
│   └── source/
│       ├── source.go
│       ├── gdrive.go
│       ├── filesystem.go
│       ├── s3.go
│       └── node.go
├── build.sh
├── go.mod
├── go.sum
└── README.md
```

## Logging

The application uses a structured logging system with three levels:

- DEBUG: Detailed information for debugging
- INFO: General operational information
- ERROR: Error conditions that need attention

Logs include timestamps, log levels, and caller information.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request 