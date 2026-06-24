package main

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestRunGreetsThroughEmbeddedTemporal proves the all-in-one binary runs a real
// durable workflow end to end with no external Temporal server.
func TestRunGreetsThroughEmbeddedTemporal(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var out bytes.Buffer
	require.NoError(t, run(ctx, &out))
	require.Equal(t, "Hello, Ada\n", out.String())
}
