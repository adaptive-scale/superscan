# SuperScan

A fast and efficient file system scanner with support for multiple storage backends. Features an elegant ASCII tree display and comprehensive logging.

## Features

- Multiple storage backends:
  - Google Drive
  - Local filesystem
  - AWS S3
  - Google Cloud Storage (coming soon)
- ASCII tree visualization
- YAML configuration
- Structured logging
- Memory-efficient scanning

## Quick Start

```bash
# Build
make build

# List local files
./bin/superscan --source-type filesystem

# List Google Drive files
./bin/superscan --source-type google-drive

# List S3 files
./bin/superscan --source-type s3
```

## Usage

### Local Filesystem

```bash
# Current directory
./bin/superscan --source-type filesystem

# Specific directory
./bin/superscan --source-type filesystem --start-path /path/to/dir
```

Example output:
```
project/
├── 📁 src/
│   ├── 📄 main.go (1024 bytes)
│   └── 📁 pkg/
│       ├── 📄 config.go (512 bytes)
│       └── 📄 utils.go (768 bytes)
└── 📄 README.md (256 bytes)
```

### Google Drive

```bash
# Root directory
./bin/superscan --source-type google-drive

# Specific folder
./bin/superscan --source-type google-drive --start-path "My Drive/Folder"
```

Example output:
```
My Drive/
├── 📁 Documents/
│   ├── 📄 report.pdf (1024 bytes)
│   └── 📁 Projects/
│       ├── 📄 design.docx (512 bytes)
│       └── 📄 notes.txt (256 bytes)
└── 📁 Photos/
    └── 📄 vacation.jpg (2048 bytes)
```

### AWS S3

```bash
# List files from S3 bucket
./bin/superscan --source-type s3

# List files from specific prefix
./bin/superscan --source-type s3 --start-path "folder/subfolder"
```

Example output:
```
my-bucket/
├── 📁 folder/
│   ├── 📄 file1.txt (1024 bytes)
│   └── 📁 subfolder/
│       ├── 📄 file2.txt (512 bytes)
│       └── 📄 file3.txt (768 bytes)
└── 📄 root-file.txt (256 bytes)
```

## Configuration

Configuration file: `~/.superscan/config.yaml`

```yaml
google_drive:
  credentials_file: /path/to/credentials.json
  token_file: /path/to/token.json
  start_path: root

s3:
  bucket: my-bucket
  region: us-east-1
  start_path: ""
```

### Environment Variables

- `SUPERSCAN_CONFIG_GOOGLE`: Path to Google Drive credentials
- `AWS_S3_BUCKET`: S3 bucket name
- `AWS_REGION`: AWS region (default: us-east-1)
- `AWS_ACCESS_KEY_ID`: AWS access key
- `AWS_SECRET_ACCESS_KEY`: AWS secret key

## Google Drive Setup

1. Go to [Google Cloud Console](https://console.cloud.google.com)
2. Create/select project
3. Enable Drive API
4. Create OAuth 2.0 credentials
5. Download as `credentials.json`

## AWS S3 Setup

1. Create an AWS account if you don't have one
2. Create an S3 bucket
3. Create an IAM user with S3 access
4. Configure AWS credentials:
   ```bash
   export AWS_ACCESS_KEY_ID="your-access-key"
   export AWS_SECRET_ACCESS_KEY="your-secret-key"
   export AWS_REGION="your-region"
   export AWS_S3_BUCKET="your-bucket-name"
   ```

## Development

```bash
# Build
make build

# Test
make test

# Lint
make lint
```

## Project Structure

```
superscan/
├── bin/                    # Binaries
├── pkg/
│   ├── config/            # Configuration
│   ├── logger/            # Logging
│   └── source/            # Storage backends
├── .gitignore
├── go.mod
├── LICENSE
└── README.md
```

## License

MIT License - see [LICENSE](LICENSE)

## Contributing

1. Fork
2. Create feature branch
3. Commit changes
4. Push
5. Open PR 