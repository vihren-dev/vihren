package blobref

// Ref is a small, history-safe handle to a stored blob. Size is advisory
// metadata for fetch/inline decisions; identity and integrity remain the digest.
type Ref[T any] struct {
	ID   BlobID `json:"id"`
	Size int64  `json:"size"`
}
