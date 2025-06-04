package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/adaptive-scale/superscan/pkg/source"
)

func main() {
	// Define command line flags
	sourceTypeStr := flag.String("source-type", "filesystem", "Type of source (filesystem|gdrive|s3|gcs)")
	startPath := flag.String("start-path", "/", "Starting path for scanning (default: /)")
	showVersion := flag.Bool("version", false, "Show version information")

	// Parse the flags
	flag.Parse()

	// Show version and exit if requested
	if *showVersion {
		PrintVersion()
		os.Exit(0)
	}

	// Parse source type
	var sourceType source.SourceType
	if err := sourceType.Set(*sourceTypeStr); err != nil {
		fmt.Printf("Error: %v\n", err)
		flag.Usage()
		os.Exit(1)
	}

	// Create source
	src, err := source.NewSource(*sourceTypeStr, nil	)
	if err != nil {
		fmt.Printf("Error creating source: %v\n", err)
		os.Exit(1)
	}

	// List files
	if err := src.ListFiles(*startPath); err != nil {
		fmt.Printf("Error listing files: %v\n", err)
		os.Exit(1)
	}
}