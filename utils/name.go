package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func IDByAddress(address string) string {

	hash := sha256.Sum256([]byte(address))

	return hex.EncodeToString(hash[:8])
}
