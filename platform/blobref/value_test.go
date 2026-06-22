package blobref

import (
	"encoding/json"
	"errors"
	"testing"
)

type testPayload struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func TestZeroReturnsEmptyValue(t *testing.T) {
	t.Parallel()

	structZero := Zero[testPayload]()
	if structZero.Inline != nil || structZero.Ref != nil {
		t.Fatalf("struct zero = %#v, want empty value", structZero)
	}
	scalarZero := Zero[int]()
	if scalarZero.Inline != nil || scalarZero.Ref != nil {
		t.Fatalf("scalar zero = %#v, want empty value", scalarZero)
	}
}

// TestRefJSONIncludesIDAndSize verifies that phantom types add no wire data
// beyond the content reference metadata.
func TestRefJSONIncludesIDAndSize(t *testing.T) {
	t.Parallel()

	ref := Ref[testPayload]{ID: ContentID([]byte("payload")), Size: 7}
	data, err := json.Marshal(ref)
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}
	expected := `{"id":"` + string(ref.ID) + `","size":7}`
	if string(data) != expected {
		t.Fatalf("ref JSON = %s, want %s", data, expected)
	}
}

// TestValueJSONRoundTrips verifies both union arms are serializable.
func TestValueJSONRoundTrips(t *testing.T) {
	t.Parallel()

	inline := Inline(testPayload{Name: "inline", Count: 1})
	var decodedInline Value[testPayload]
	roundTripValue(t, inline, &decodedInline)
	if decodedInline.Inline == nil || decodedInline.Ref != nil || *decodedInline.Inline != *inline.Inline {
		t.Fatalf("inline value mismatch after JSON round trip: %#v", decodedInline)
	}

	ref := Ref[testPayload]{ID: ContentID([]byte("stored")), Size: 6}
	refValue := RefValue(ref)
	var decodedRef Value[testPayload]
	roundTripValue(t, refValue, &decodedRef)
	if decodedRef.Ref == nil || decodedRef.Inline != nil || decodedRef.Ref.ID != ref.ID || decodedRef.Ref.Size != ref.Size {
		t.Fatalf("ref value mismatch after JSON round trip: %#v", decodedRef)
	}
}

// TestDecodeInlineDecodesRawJSON verifies workflow-safe raw JSON decoding.
func TestDecodeInlineDecodesRawJSON(t *testing.T) {
	t.Parallel()

	raw := json.RawMessage(`{"name":"inline","count":2}`)
	decoded, err := DecodeInline[testPayload](Inline(raw))
	if err != nil {
		t.Fatalf("DecodeInline returned error: %v", err)
	}
	if decoded != (testPayload{Name: "inline", Count: 2}) {
		t.Fatalf("decoded value mismatch: %#v", decoded)
	}
}

// TestDecodeInlineRejectsRef verifies workflow-safe decoding never fetches refs.
func TestDecodeInlineRejectsRef(t *testing.T) {
	t.Parallel()

	ref := Ref[json.RawMessage]{ID: ContentID([]byte(`{"name":"stored"}`)), Size: 17}
	if _, err := DecodeInline[testPayload](RefValue(ref)); !errors.Is(err, ErrInvalidValue) {
		t.Fatalf("DecodeInline ref error = %v, want ErrInvalidValue", err)
	}
}

// TestEncodeInlineEnforcesThreshold verifies inline-only encoding is bounded.
func TestEncodeInlineEnforcesThreshold(t *testing.T) {
	t.Parallel()

	value := testPayload{Name: "abc", Count: 1}
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}
	encoded, err := EncodeInline(value, len(data))
	if err != nil {
		t.Fatalf("EncodeInline at threshold returned error: %v", err)
	}
	if encoded.Inline == nil || string(*encoded.Inline) != string(data) {
		t.Fatalf("unexpected inline encoding: %#v", encoded)
	}
	if _, err := EncodeInline(value, len(data)-1); !errors.Is(err, ErrInlineTooLarge) {
		t.Fatalf("EncodeInline over threshold error = %v, want ErrInlineTooLarge", err)
	}
}

// roundTripValue marshals input and unmarshals it into output.
func roundTripValue[T any](t *testing.T, input Value[T], output *Value[T]) {
	t.Helper()
	data, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}
	if err := json.Unmarshal(data, output); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}
}
