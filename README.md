# SuperScan

A fast and efficient file system scanner with support for multiple storage backends. Features an elegant ASCII tree display and comprehensive logging.

## Features

- Multiple storage backends:
  - Google Drive
  - Local filesystem
  - AWS S3 (coming soon)
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

## Configuration

Configuration file: `~/.superscan/config.yaml`

```yaml
google_drive:
  credentials_file: /path/to/credentials.json
  token_file: /path/to/token.json
  start_path: root
```

### Environment Variables

- `SUPERSCAN_CONFIG_GOOGLE`: Path to Google Drive credentials

## Google Drive Setup

1. Go to [Google Cloud Console](https://console.cloud.google.com)
2. Create/select project
3. Enable Drive API
4. Create OAuth 2.0 credentials
5. Download as `credentials.json`

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