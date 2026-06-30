package codegen

import (
	"strings"
	"testing"
)

// TestRenderPackagePreservesActivityShapes proves generated proxies keep the
// ADR's SDK-compatible activity argument and return shapes.
func TestRenderPackagePreservesActivityShapes(t *testing.T) {
	t.Parallel()
	packages, diagnostics, err := Discover(
		DiscoverConfig{Dir: moduleRoot(t)},
		"./internal/codegen/testdata/fixtures/basic",
	)
	if err != nil {
		t.Fatalf("discover packages: %v", err)
	}
	if len(diagnostics) != 0 {
		t.Fatalf("diagnostics = %#v, want none", diagnostics)
	}
	source, err := RenderPackage(packages[0])
	if err != nil {
		t.Fatalf("render package: %v", err)
	}
	text := string(source)
	for _, want := range []string{
		"type activityRegistry interface",
		"func RegisterActivities(r activityRegistry, activities *Activities)",
		"r.RegisterActivityWithOptions(activities.ChargeCard, activity.RegisterOptions{Name: ChargeActivityName})",
		"func (activityProxy) Price(ctx workflow.Context, customer CustomerID, cents int) (Receipt, error)",
		"workflow.ExecuteActivity(ctx, PriceActivityName, customer, cents).Get(ctx, &out)",
		"func (activityProxy) Ping(ctx workflow.Context) error",
		"workflow.ExecuteActivity(ctx, PingActivityName).Get(ctx, nil)",
		"type workflowRegistry interface",
		"func RegisterWorkflows(r workflowRegistry)",
		"type CheckoutRun struct",
		"func (run CheckoutRun) Get(ctx context.Context) (CheckoutResult, error)",
		"func (cl Client) Checkout(ctx context.Context, opts client.StartWorkflowOptions, in CheckoutRequest) (CheckoutResult, error)",
		"func (cl Client) CheckoutAsync(ctx context.Context, opts client.StartWorkflowOptions, in CheckoutRequest) (CheckoutRun, error)",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("generated source missing %q:\n%s", want, text)
		}
	}
}
