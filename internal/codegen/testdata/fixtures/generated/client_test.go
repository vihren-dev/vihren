package generated

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/client"
)

// TestClientUsesNarrowWorkflowStarter proves generated client sync and async
// methods can be unit-tested without a live Temporal server.
func TestClientUsesNarrowWorkflowStarter(t *testing.T) {
	t.Parallel()
	starter := &fakeWorkflowStarter{
		run: fakeWorkflowRun{
			result: CheckoutResult{ReceiptID: "receipt", Clean: "clean"},
		},
	}
	cl := Client{c: starter}
	result, err := cl.Checkout(context.Background(), client.StartWorkflowOptions{TaskQueue: "billing"}, CheckoutRequest{Amount: 10})
	require.NoError(t, err)
	require.Equal(t, CheckoutResult{ReceiptID: "receipt", Clean: "clean"}, result)
	require.Equal(t, CheckoutWorkflowName, starter.workflow)
	require.Equal(t, []interface{}{CheckoutRequest{Amount: 10}}, starter.args)
	run, err := cl.CheckoutAsync(context.Background(), client.StartWorkflowOptions{TaskQueue: "billing"}, CheckoutRequest{Amount: 11})
	require.NoError(t, err)
	require.Equal(t, "run-id", run.GetRunID())
}

type fakeWorkflowStarter struct {
	workflow interface{}
	args     []interface{}
	run      client.WorkflowRun
}

func (starter *fakeWorkflowStarter) ExecuteWorkflow(
	ctx context.Context,
	options client.StartWorkflowOptions,
	workflow interface{},
	args ...interface{},
) (client.WorkflowRun, error) {
	_ = ctx
	_ = options
	starter.workflow = workflow
	starter.args = append([]interface{}(nil), args...)
	return starter.run, nil
}

type fakeWorkflowRun struct {
	result CheckoutResult
}

func (run fakeWorkflowRun) GetID() string { return "workflow-id" }

func (run fakeWorkflowRun) GetRunID() string { return "run-id" }

func (run fakeWorkflowRun) Get(ctx context.Context, valuePtr interface{}) error {
	return run.GetWithOptions(ctx, valuePtr, client.WorkflowRunGetOptions{})
}

func (run fakeWorkflowRun) GetWithOptions(
	ctx context.Context,
	valuePtr interface{},
	options client.WorkflowRunGetOptions,
) error {
	_ = ctx
	_ = options
	if out, ok := valuePtr.(*CheckoutResult); ok {
		*out = run.result
	}
	return nil
}
