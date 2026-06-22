package blobref

import "errors"

var (
	// ErrInvalidID means a blob ID is not in the canonical sha256 digest form.
	ErrInvalidID = errors.New("blobref: invalid blob ID")

	// ErrInvalidValue means a claim-check value has zero or multiple active arms,
	// or a helper was asked to materialize a ref in workflow-safe code.
	ErrInvalidValue = errors.New("blobref: invalid value")

	// ErrInlineTooLarge means an inline-only encoding would exceed its configured
	// history-safe byte threshold.
	ErrInlineTooLarge = errors.New("blobref: inline value exceeds threshold")
)
