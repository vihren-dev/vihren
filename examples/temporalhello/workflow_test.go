package temporalhello

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/testsuite"
	"go.temporal.io/sdk/workflow"
)

// TestHelloWorkflowUsesManualTemporalRegistration proves the public example
// runs through explicit Temporal SDK registration and activity execution.
func TestHelloWorkflowUsesManualTemporalRegistration(t *testing.T) {
	t.Parallel()
	var suite testsuite.WorkflowTestSuite
	env := suite.NewTestWorkflowEnvironment()
	activities := &GreetingActivities{Prefix: "hello"}
	env.RegisterActivityWithOptions(
		activities.ComposeGreeting,
		activity.RegisterOptions{Name: ComposeGreetingActivityName},
	)
	env.RegisterWorkflowWithOptions(
		HelloWorkflow,
		workflow.RegisterOptions{Name: HelloWorkflowName},
	)

	env.ExecuteWorkflow(HelloWorkflowName, GreetingInput{Name: "Ada"})

	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
	var result GreetingOutput
	require.NoError(t, env.GetWorkflowResult(&result))
	require.Equal(t, GreetingOutput{Message: "hello, Ada"}, result)
}
