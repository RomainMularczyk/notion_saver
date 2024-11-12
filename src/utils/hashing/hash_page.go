package hashing

import (
	"crypto/sha256"
	"encoding/hex"
)

// Hash a page content
func HashPage(pageContent string) string {
	hashSum := sha256.Sum256([]byte(pageContent))
	return hex.EncodeToString(hashSum[:])
}

// Compare hashes of two pages
func IsSamePage(pageA string, pageB string) bool {
	hashA := HashPage(pageA)
	hashB := HashPage(pageB)
	return hashA == hashB
}
