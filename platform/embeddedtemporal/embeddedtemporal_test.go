package embeddedtemporal_test

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"

	"github.com/vihren-dev/vihren/platform/embeddedtemporal"
)

const (
	echoWorkflowName = "EmbeddedTemporalEchoWorkflow"
	echoActivityName = "EmbeddedTemporalEchoActivity"
	echoTaskQueue    = "embeddedtemporal-test"
)

// TestStartRunsWorkflowEndToEnd proves the one-line ephemeral server hosts a real
// durable workflow through a registered worker and a connected client, with no
// external Temporal process.
func TestStartRunsWorkflowEndToEnd(t *testing.T) {
	server, err := embeddedtemporal.Start()
	require.NoError(t, err)
	defer server.Close()

	require.Contains(t, server.HostPort(), ":")
	require.Equal(t, embeddedtemporal.DefaultNamespace, server.Namespace())

	registerEcho(t, server)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	var result string
	require.NoError(t, runEcho(ctx, server, "echo-ephemeral", "Ada", &result))
	require.Equal(t, "hello, Ada", result)
}

// TestPersistentStateSurvivesRestart proves the headline desktop capability: run
// a workflow, stop the whole server, start it again on the same database file,
// and recover the completed workflow's result. This exercises both file-backed
// persistence and the startup schema setup/migration path.
func TestPersistentStateSurvivesRestart(t *testing.T) {
	databaseFile := filepath.Join(t.TempDir(), "agent", "state.db")
	const workflowID = "echo-persistent"

	// First run: execute a workflow to completion, then shut everything down.
	first, err := embeddedtemporal.Start(embeddedtemporal.WithDatabaseFile(databaseFile))
	require.NoError(t, err)
	registerEcho(t, first)
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()
	var firstResult string
	require.NoError(t, runEcho(ctx, first, workflowID, "Grace", &firstResult))
	require.Equal(t, "hello, Grace", firstResult)
	first.Close()

	// Second run: a brand new server on the same file. No worker is registered,
	// so the result can only come from persisted history.
	second, err := embeddedtemporal.Start(embeddedtemporal.WithDatabaseFile(databaseFile))
	require.NoError(t, err)
	defer second.Close()

	var recovered string
	require.NoError(t, second.Client().GetWorkflow(ctx, workflowID, "").Get(ctx, &recovered))
	require.Equal(t, "hello, Grace", recovered)
}

// TestPersistentLockRejectsSecondInstance proves the single-writer guard: a
// second server cannot open a database file already held by a running server.
func TestPersistentLockRejectsSecondInstance(t *testing.T) {
	databaseFile := filepath.Join(t.TempDir(), "state.db")

	first, err := embeddedtemporal.Start(embeddedtemporal.WithDatabaseFile(databaseFile))
	require.NoError(t, err)
	defer first.Close()

	_, err = embeddedtemporal.Start(embeddedtemporal.WithDatabaseFile(databaseFile))
	require.Error(t, err)
	require.Contains(t, err.Error(), "locked")

	// After the holder closes, the lock is released and a new server can open it.
	first.Close()
	reopened, err := embeddedtemporal.Start(embeddedtemporal.WithDatabaseFile(databaseFile))
	require.NoError(t, err)
	reopened.Close()
}

// TestStartWorkerRequiresRegisterFunc keeps the helper from silently starting an
// empty worker.
func TestStartWorkerRequiresRegisterFunc(t *testing.T) {
	server, err := embeddedtemporal.Start()
	require.NoError(t, err)
	defer server.Close()

	_, err = server.StartWorker(echoTaskQueue, nil)
	require.EqualError(t, err, "embedded Temporal worker requires a register function")
}

func registerEcho(t *testing.T, server *embeddedtemporal.Server) {
	t.Helper()
	_, err := server.StartWorker(echoTaskQueue, func(registry worker.Registry) {
		registry.RegisterWorkflowWithOptions(echoWorkflow, workflow.RegisterOptions{Name: echoWorkflowName})
		registry.RegisterActivityWithOptions(echoActivity, activity.RegisterOptions{Name: echoActivityName})
	})
	require.NoError(t, err)
}

func runEcho(ctx context.Context, server *embeddedtemporal.Server, workflowID, name string, out *string) error {
	run, err := server.Client().ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{ID: workflowID, TaskQueue: echoTaskQueue},
		echoWorkflowName,
		name,
	)
	if err != nil {
		return err
	}
	return run.Get(ctx, out)
}

func echoWorkflow(ctx workflow.Context, name string) (string, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{StartToCloseTimeout: 5 * time.Second})
	var out string
	if err := workflow.ExecuteActivity(ctx, echoActivityName, name).Get(ctx, &out); err != nil {
		return "", err
	}
	return out, nil
}

func echoActivity(_ context.Context, name string) (string, error) {
	return fmt.Sprintf("hello, %s", name), nil
}
