package computation

import (
    "fmt"
    "math/rand"
    "time"
)

// Generic random selection function
func SelectRandom[T any](items []T, n int) []T {
    if n > len(items) {
        n = len(items)
    }
    rand.Seed(time.Now().UnixNano())
    perm := rand.Perm(len(items))
    result := make([]T, n)
    for i := 0; i < n; i++ {
        result[i] = items[perm[i]]
    }
    return result
}
