// Package blobref models activity input and output values that may cross
// Temporal workflow history, where large payloads would otherwise bloat the
// event history and threaten replay cost and durability.
//
// A Value[T] carries exactly one of an inline typed value (small enough to
// travel in history) or a content-addressed Ref[T] handle to the same value
// held in a blob store. A BlobID is the content digest, so references are
// stable, deduplicating, and integrity-checked. Workflow code stays I/O-free:
// it may decode inline values but never fetches blobs, leaving fetch and
// offload decisions to the imperative shell around it.
package blobref

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
)

// Value carries exactly one of an inline typed value or a reference to a stored
// typed value for activity inputs and outputs that may cross Temporal history.
type Value[T any] struct {
	Inline *T      `json:"inline,omitempty"`
	Ref    *Ref[T] `json:"ref,omitempty"`
}

// Zero returns the zero Value[T].
func Zero[T any]() Value[T] {
	return Value[T]{}
}

// Inline constructs a Value whose data travels inline.
func Inline[T any](value T) Value[T] {
	return Value[T]{Inline: &value}
}

// RefValue constructs a Value whose data travels by blob reference.
func RefValue[T any](ref Ref[T]) Value[T] {
	return Value[T]{Ref: &ref}
}

// DecodeInline decodes an inline raw JSON value into T without performing I/O.
// Ref-shaped or invalid values fail because workflow code must not fetch blobs.
func DecodeInline[T any](value Value[json.RawMessage]) (T, error) {
	var decoded T
	if value.Inline == nil || value.Ref != nil {
		return decoded, ErrInvalidValue
	}
	decoder := json.NewDecoder(bytes.NewReader(*value.Inline))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&decoded); err != nil {
		return decoded, err
	}
	var trailing struct{}
	if err := decoder.Decode(&trailing); err != nil {
		if errors.Is(err, io.EOF) {
			return decoded, nil
		}
		return decoded, err
	}
	return decoded, ErrInvalidValue
}

// EncodeInline marshals value to raw JSON and returns it inline only when the
// encoded bytes fit within threshold.
func EncodeInline[T any](value T, threshold int) (Value[json.RawMessage], error) {
	data, err := json.Marshal(value)
	if err != nil {
		return Value[json.RawMessage]{}, err
	}
	if len(data) > threshold {
		return Value[json.RawMessage]{}, ErrInlineTooLarge
	}
	raw := json.RawMessage(data)
	return Inline(raw), nil
}
