package computation

import (
    "fmt"
    "math"
    "path/filepath"
    "strings"
)

// File struct
type File struct {
    Name string
    Size int64 // in bytes
}

// Helper to get extension, normalized
func getExtension(filename string) string {
    ext := strings.ToLower(filepath.Ext(filename))
    if ext != "" && ext[0] == '.' {
        return ext[1:] // strip dot
    }
    return ext
}

// Similarity function
func fileSimilarity(a, b File) float64 {
    extA := getExtension(a.Name)
    extB := getExtension(b.Name)
    if extA != extB {
        return 0
    }
    maxSize := math.Max(float64(a.Size), float64(b.Size))
    if maxSize == 0 {
        if a.Size == b.Size {
            return 1
        }
        return 0
    }
    diff := math.Abs(float64(a.Size - b.Size))
    return 1 - diff/maxSize
}