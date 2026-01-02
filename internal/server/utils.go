package server

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
)

func generateCode() string {
	bytes := make([]byte, 3)
	rand.Read(bytes)
	return strings.ToUpper(hex.EncodeToString(bytes))
}
