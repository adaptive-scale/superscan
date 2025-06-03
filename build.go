package main

import (
	"fmt"
	"runtime"
)

var (
	// Version is the application version
	Version = "dev"
	// BuildTime is the time when the application was built
	BuildTime = "unknown"
	// GitCommit is the git commit hash
	GitCommit = "unknown"
	// GoVersion is the Go version used to build the application
	GoVersion = runtime.Version()
)

// PrintVersion prints the version information
func PrintVersion() {
	fmt.Printf("SuperScan Version: %s\n", Version)
	fmt.Printf("Build Time: %s\n", BuildTime)
	fmt.Printf("Git Commit: %s\n", GitCommit)
	fmt.Printf("Go Version: %s\n", GoVersion)
	fmt.Printf("OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
} 