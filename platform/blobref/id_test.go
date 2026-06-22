package blobref

import (
	"errors"
	"strings"
	"testing"
)

// TestContentIDStable verifies that content addressing is deterministic.
func TestContentIDStable(t *testing.T) {
	t.Parallel()

	first := ContentID([]byte("same"))
	second := ContentID([]byte("same"))
	if first != second {
		t.Fatalf("equal bytes produced different IDs: %q != %q", first, second)
	}
	if first == ContentID([]byte("different")) {
		t.Fatalf("different bytes produced the same ID: %q", first)
	}
}

// TestParseBlobIDAcceptsCanonical verifies the accepted wire form.
func TestParseBlobIDAcceptsCanonical(t *testing.T) {
	t.Parallel()

	expected := ContentID([]byte("canonical"))
	parsed, err := ParseBlobID(string(expected))
	if err != nil {
		t.Fatalf("ParseBlobID returned error for canonical ID: %v", err)
	}
	if parsed != expected {
		t.Fatalf("parsed ID mismatch: got %q want %q", parsed, expected)
	}
}

// TestParseBlobIDRejectsNonCanonical verifies hostile and malformed IDs fail.
func TestParseBlobIDRejectsNonCanonical(t *testing.T) {
	t.Parallel()

	validDigest := strings.TrimPrefix(string(ContentID([]byte("canonical"))), blobIDPrefix)
	invalidIDs := []string{
		"md5:" + validDigest,
		"sha256:" + validDigest[:len(validDigest)-1],
		strings.ToUpper(string(ContentID([]byte("canonical")))),
		"sha256:" + validDigest[:len(validDigest)-1] + "g",
		"sha256:../../etc/passwd",
		"sha256:" + validDigest[:8] + "/" + validDigest[9:],
		"sha256:" + validDigest[:8] + ".." + validDigest[10:],
	}
	for _, invalidID := range invalidIDs {
		if _, err := ParseBlobID(invalidID); !errors.Is(err, ErrInvalidID) {
			t.Fatalf("ParseBlobID(%q) error = %v, want ErrInvalidID", invalidID, err)
		}
	}
}
