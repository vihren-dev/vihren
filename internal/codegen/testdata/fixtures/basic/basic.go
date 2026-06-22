package basic

import (
	"context"

	"go.temporal.io/sdk/workflow"
)

// Activities holds dependencies for activity marker tests.
type Activities struct{}

// ChargeRequest is a serializable activity input.
type ChargeRequest struct {
	Amount int
}

// Receipt is a serializable activity output.
type Receipt struct {
	ID string
}

// Raw is a serializable raw input.
type Raw struct {
	Value string
}

// Clean is a serializable clean output.
type Clean struct {
	Value string
}

// CustomerID is a primitive alias used by a multi-argument activity.
type CustomerID string

// CheckoutRequest is a serializable workflow input.
type CheckoutRequest struct {
	Amount int
}

// CheckoutResult is a serializable workflow output.
type CheckoutResult struct {
	Receipt Receipt
}

// ChargeCard charges a card.
//
//vihren:activity proxy=Charge
func (activities *Activities) ChargeCard(ctx context.Context, in ChargeRequest) (Receipt, error) {
	_ = activities
	_ = ctx
	_ = in
	return Receipt{}, nil
}

// Normalize cleans a raw value.
//
//vihren:activity name=external.normalize
func Normalize(in Raw) (Clean, error) {
	_ = in
	return Clean{}, nil
}

// Price computes a value from multiple business parameters.
//
//vihren:activity
func Price(customer CustomerID, cents int) (Receipt, error) {
	_ = customer
	_ = cents
	return Receipt{}, nil
}

// Ping has no input and no output.
//
//vihren:activity
func Ping() error {
	return nil
}

// Checkout is a marked workflow.
//
//vihren:workflow
func Checkout(ctx workflow.Context, in CheckoutRequest) (CheckoutResult, error) {
	_ = ctx
	_ = in
	return CheckoutResult{}, nil
}
