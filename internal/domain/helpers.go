package domain

import (
	"crypto/rand"
	"fmt"
)

//TODO: Move this to internal

// Generates an id with a prefix and length of 16 chars
// conflict is possible but very unlikely, could be mitigated by checking
// for existing ids in the database and regenerating if a conflict is found.
func GenerateID(prefix string) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	bytes := make([]byte, 16)
	rand.Read(bytes)
	for i, b := range bytes {
		// Modulo bias exists since 62 does not evenly divide 256, but is negligible for ID generation.
		bytes[i] = charset[b%byte(len(charset))]
	}

	return fmt.Sprintf("%s-%s", prefix, string(bytes))
}
