package codegenhello

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
)

// TestHelloWorkflowUsesGeneratedRegistration proves the public example runs
// through generated registration and activity proxy code.
func TestHelloWorkflowUsesGeneratedRegistration(t *testing.T) {
	t.Parallel()
	var suite testsuite.WorkflowTestSuite
	env := suite.NewTestWorkflowEnvironment()
	Register(env, &GreetingActivities{Prefix: "hello"})
	env.ExecuteWorkflow(HelloWorkflowName, GreetingInput{Name: "Ada"})
	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
	var result GreetingOutput
	require.NoError(t, env.GetWorkflowResult(&result))
	require.Equal(t, GreetingOutput{Message: "hello, Ada"}, result)
}
