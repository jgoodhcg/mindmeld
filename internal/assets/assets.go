package assets

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"os"
)

// ComputeFileHash reads a file and returns its MD5 hash as a hex string.
// Returns "dev" if the file cannot be read (useful for development).
func ComputeFileHash(path string) string {
	f, err := os.Open(path)
	if err != nil {
		log.Printf("Warning: could not open %s for hashing: %v", path, err)
		return "dev"
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Printf("Warning: could not read %s for hashing: %v", path, err)
		return "dev"
	}

	return hex.EncodeToString(h.Sum(nil))[:8] // First 8 chars is enough
}
