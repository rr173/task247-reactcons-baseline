package version

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

// ContentHash returns a stable SHA-256 hex digest of the canonical JSON encoding
// of payload. Maps are encoded with sorted keys so the hash is reproducible.
func ContentHash(payload any) (string, error) {
	var data []byte
	var err error
	switch p := payload.(type) {
	case []byte:
		data = p
	case string:
		data = []byte(p)
	default:
		data, err = json.Marshal(payload)
		if err != nil {
			return "", err
		}
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}
