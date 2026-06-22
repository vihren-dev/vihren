package blobref

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

const (
	blobIDAlgorithm = "sha256"
	blobIDPrefix    = blobIDAlgorithm + ":"
	blobIDHexLength = sha256.Size * 2
)

// BlobID identifies immutable blob bytes by content digest. Its canonical form
// is "sha256:<64 lowercase hex characters>".
type BlobID string

// ContentID forms the canonical content ID for already-serialized bytes.
func ContentID(data []byte) BlobID {
	sum := sha256.Sum256(data)
	return BlobID(blobIDPrefix + hex.EncodeToString(sum[:]))
}

// ParseBlobID validates and returns a canonical sha256 content ID.
func ParseBlobID(s string) (BlobID, error) {
	algorithm, digest, ok := strings.Cut(s, ":")
	if !ok || algorithm != blobIDAlgorithm || len(digest) != blobIDHexLength {
		return "", ErrInvalidID
	}
	for _, character := range digest {
		if !isLowerHex(character) {
			return "", ErrInvalidID
		}
	}
	return BlobID(s), nil
}

// isLowerHex reports whether character is valid lowercase hexadecimal.
func isLowerHex(character rune) bool {
	return ('0' <= character && character <= '9') ||
		('a' <= character && character <= 'f')
}
