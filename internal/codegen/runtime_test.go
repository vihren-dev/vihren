package codegen

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
	"go.temporal.io/sdk/worker"

	"github.com/vihren-dev/vihren/internal/codegen/testdata/fixtures/generated"
	"github.com/vihren-dev/vihren/internal/codegen/testdata/fixtures/workflowonly"
)

var (
	// worker.Worker must satisfy worker.Registry for native worker registration.
	_ worker.Registry = (worker.Worker)(nil)

	// TestWorkflowEnvironment must satisfy worker.Registry for unit tests.
	_ worker.Registry = (*testsuite.TestWorkflowEnvironment)(nil)
)

// TestGeneratedRegisterRunsWorkflowThroughProxy proves generated Register and
// Activity proxy code executes in the Temporal workflow test environment.
func TestGeneratedRegisterRunsWorkflowThroughProxy(t *testing.T) {
	t.Parallel()
	var suite testsuite.WorkflowTestSuite
	env := suite.NewTestWorkflowEnvironment()
	generated.Register(
		env,
		&generated.BillingActivities{Prefix: "paid"},
		&generated.RefundActivities{Ledger: "ledger"},
	)
	env.ExecuteWorkflow(generated.CheckoutWorkflowName, generated.CheckoutRequest{
		Amount: 10,
		Raw:    "  county record  ",
		Refund: true,
	})
	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
	var result generated.CheckoutResult
	require.NoError(t, env.GetWorkflowResult(&result))
	require.Equal(t, generated.CheckoutResult{
		ReceiptID: "paid-charge",
		Clean:     "county record",
		RefundID:  "ledger-refund",
	}, result)
}

// TestWorkflowOnlyRegisterShape proves packages without activity deps can still
// expose a no-arg generated Register.
func TestWorkflowOnlyRegisterShape(t *testing.T) {
	t.Parallel()
	var suite testsuite.WorkflowTestSuite
	env := suite.NewTestWorkflowEnvironment()
	workflowonly.Register(env)
	env.ExecuteWorkflow(workflowonly.EchoWorkflowName, workflowonly.EchoInput{Value: "ok"})
	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
	var result workflowonly.EchoOutput
	require.NoError(t, env.GetWorkflowResult(&result))
	require.Equal(t, workflowonly.EchoOutput{Value: "ok"}, result)
}

// TestGeneratedRegisterRejectsNilReceivers proves generated nil checks fail
// before Temporal's registration reflection panics.
func TestGeneratedRegisterRejectsNilReceivers(t *testing.T) {
	t.Parallel()
	var suite testsuite.WorkflowTestSuite
	env := suite.NewTestWorkflowEnvironment()
	require.PanicsWithValue(t, "generated.RegisterActivities: *BillingActivities is nil", func() {
		generated.Register(env, nil, &generated.RefundActivities{Ledger: "ledger"})
	})
	require.PanicsWithValue(t, "generated.RegisterActivities: *RefundActivities is nil", func() {
		generated.Register(env, &generated.BillingActivities{Prefix: "paid"}, nil)
	})
}
