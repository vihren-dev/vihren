package invalid

import (
	"context"
	"unsafe"

	"go.temporal.io/sdk/workflow"
)

// ChargeRequest is a valid exported input for invalid-shape tests.
type ChargeRequest struct {
	Amount int
}

// Receipt is a valid exported output for invalid-shape tests.
type Receipt struct {
	ID string
}

// ContainsFunc contains a function field.
type ContainsFunc struct {
	Run func()
}

// ContainsChannel contains a channel field.
type ContainsChannel struct {
	Events chan int
}

// ContainsUnsafePointer contains unsafe.Pointer.
type ContainsUnsafePointer struct {
	Pointer unsafe.Pointer
}

// ContainsNonStringMap contains a map whose key is not a string.
type ContainsNonStringMap struct {
	Values map[int]string
}

// BadField uses a function field.
//
//vihren:activity
func BadField(ctx context.Context, in ContainsFunc) (Receipt, error) {
	_ = ctx
	return Receipt{}, nil
}

// BadChannel uses a channel field.
//
//vihren:activity
func BadChannel(ctx context.Context, in ContainsChannel) (Receipt, error) {
	_ = ctx
	return Receipt{}, nil
}

// BadUnsafe uses unsafe.Pointer.
//
//vihren:activity
func BadUnsafe(ctx context.Context, in ContainsUnsafePointer) (Receipt, error) {
	_ = ctx
	return Receipt{}, nil
}

// BadMap uses a non-string map key.
//
//vihren:activity
func BadMap(ctx context.Context, in ContainsNonStringMap) (Receipt, error) {
	_ = ctx
	return Receipt{}, nil
}

// ActivityWorkflowContext uses workflow.Context in an activity.
//
//vihren:activity
func ActivityWorkflowContext(ctx workflow.Context, in ChargeRequest) (Receipt, error) {
	_ = ctx
	return Receipt{}, nil
}

// DuplicateOne claims a proxy name also used by DuplicateTwo.
//
//vihren:activity proxy=Duplicate
func DuplicateOne(ctx context.Context, in ChargeRequest) (Receipt, error) {
	_ = ctx
	return Receipt{}, nil
}

// DuplicateTwo claims a proxy name also used by DuplicateOne.
//
//vihren:activity proxy=Duplicate
func DuplicateTwo(ctx context.Context, in ChargeRequest) (Receipt, error) {
	_ = ctx
	return Receipt{}, nil
}

// WrongWorkflowContext uses context.Context for a workflow.
//
//vihren:workflow timeout=30s
func WrongWorkflowContext(ctx context.Context, in ChargeRequest) (Receipt, error) {
	_ = ctx
	return Receipt{}, nil
}

// CorrectWorkflow keeps workflow.Context imported and used.
func CorrectWorkflow(ctx workflow.Context) error {
	_ = ctx
	return nil
}
