package generated

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// BillingActivities holds dependencies for billing activities.
type BillingActivities struct {
	Prefix string
}

// RefundActivities holds dependencies for refund activities.
type RefundActivities struct {
	Ledger string
}

// ChargeRequest is the generated-fixture charge input.
type ChargeRequest struct {
	Amount int
}

// Receipt is the generated-fixture charge output.
type Receipt struct {
	ID string
}

// Raw is the generated-fixture normalization input.
type Raw struct {
	Value string
}

// Clean is the generated-fixture normalization output.
type Clean struct {
	Value string
}

// CustomerID is a primitive alias used by a multi-argument activity.
type CustomerID string

// RefundRequest is the generated-fixture refund input.
type RefundRequest struct {
	Amount int
}

// RefundRecord is the generated-fixture refund output.
type RefundRecord struct {
	ID string
}

// CheckoutRequest is the generated-fixture workflow input.
type CheckoutRequest struct {
	Amount int
	Raw    string
	Refund bool
}

// CheckoutResult is the generated-fixture workflow output.
type CheckoutResult struct {
	ReceiptID string
	Clean     string
	RefundID  string
}

// Normalize trims a raw value without dependencies.
//
//vihren:activity
func Normalize(ctx context.Context, in Raw) (Clean, error) {
	_ = ctx
	return Clean{Value: strings.TrimSpace(in.Value)}, nil
}

// Price computes a value from multiple business parameters and no activity context.
//
//vihren:activity
func Price(customer CustomerID, cents int) (Receipt, error) {
	return Receipt{ID: fmt.Sprintf("%s-%d", customer, cents)}, nil
}

// Ping proves error-only activities with no inputs are part of the generated
// proxy shape.
//
//vihren:activity
func Ping() error {
	return nil
}

// ChargeCard charges the requested amount.
//
//vihren:activity
func (activities *BillingActivities) ChargeCard(ctx context.Context, in ChargeRequest) (Receipt, error) {
	_ = ctx
	_ = in
	return Receipt{ID: activities.Prefix + "-charge"}, nil
}

// Refund refunds the requested amount.
//
//vihren:activity
func (activities *RefundActivities) Refund(ctx context.Context, in RefundRequest) (RefundRecord, error) {
	_ = ctx
	_ = in
	return RefundRecord{ID: activities.Ledger + "-refund"}, nil
}

// Checkout calls generated activity proxies from a workflow.
//
//vihren:workflow
func Checkout(ctx workflow.Context, in CheckoutRequest) (CheckoutResult, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Second,
		RetryPolicy:         &temporal.RetryPolicy{MaximumAttempts: 1},
	})
	receipt, err := Activity.ChargeCard(ctx, ChargeRequest{Amount: in.Amount})
	if err != nil {
		return CheckoutResult{}, err
	}
	clean, err := Activity.Normalize(ctx, Raw{Value: in.Raw})
	if err != nil {
		return CheckoutResult{}, err
	}
	if _, err := Activity.Price(ctx, CustomerID("customer"), in.Amount); err != nil {
		return CheckoutResult{}, err
	}
	if err := Activity.Ping(ctx); err != nil {
		return CheckoutResult{}, err
	}
	result := CheckoutResult{ReceiptID: receipt.ID, Clean: clean.Value}
	if in.Refund {
		refund, err := Activity.Refund(ctx, RefundRequest{Amount: in.Amount})
		if err != nil {
			return CheckoutResult{}, err
		}
		result.RefundID = refund.ID
	}
	return result, nil
}
