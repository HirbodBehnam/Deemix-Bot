package rng

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"time"
)

// QueryToken creates a query token for inline queries
// The text contains the unix time in hex + random hex bytes
// The result is always 64 chars
func QueryToken() string {
	resultBytes := make([]byte, 32)
	binary.LittleEndian.PutUint64(resultBytes, uint64(time.Now().Unix()))
	_, _ = rand.Read(resultBytes[8:])
	return hex.EncodeToString(resultBytes)
}
