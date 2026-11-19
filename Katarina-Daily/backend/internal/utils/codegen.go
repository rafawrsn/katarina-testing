package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"time"
)

const randomCodeAlphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" // no 0/O/1/I ambiguity

// RandomCode returns a cryptographically random alphanumeric string of
// the given length, suitable for order numbers and uniqueness suffixes.
func RandomCode(length int) string {
	var sb strings.Builder
	sb.Grow(length)
	max := big.NewInt(int64(len(randomCodeAlphabet)))
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			// Extremely unlikely; fall back to time-based entropy.
			sb.WriteByte(randomCodeAlphabet[time.Now().UnixNano()%int64(len(randomCodeAlphabet))])
			continue
		}
		sb.WriteByte(randomCodeAlphabet[n.Int64()])
	}
	return sb.String()
}

// GenerateOrderNumber builds a human-readable, mostly-unique order
// number. The repository layer still enforces the DB unique
// constraint and retries on the rare collision.
func GenerateOrderNumber() string {
	return fmt.Sprintf("ORD%s%s", time.Now().UTC().Format("20060102"), RandomCode(8))
}
