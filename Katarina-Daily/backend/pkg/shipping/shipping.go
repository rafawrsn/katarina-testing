// Package shipping defines a provider-agnostic abstraction for
// calculating shipping cost/rates. Checkout business logic depends only
// on the Service interface below, never on a specific courier API, so
// a real Indonesian shipping aggregator (RajaOngkir, Biteship, etc.)
// can be plugged in later without touching the checkout flow.
package shipping

import (
	"context"
	"errors"
)

// Location describes the origin or destination of a shipment for rate
// calculation purposes.
type Location struct {
	Province string
	City     string
	District string
	Postal   string
}

// RateRequest is everything a provider needs to quote shipping cost.
type RateRequest struct {
	Origin      Location
	Destination Location
	// WeightGrams is the total package weight in grams.
	WeightGrams int
	// Courier optionally filters results to a single courier code
	// (e.g. "jne"). Empty means "return all supported couriers".
	Courier string
}

// RateOption is a single courier+service quote.
type RateOption struct {
	Courier     string
	CourierName string
	Service     string
	Description string
	Cost        int64 // in the smallest currency unit (IDR has no
	// subunit in practice, so this is plain Rupiah)
	ETD string // estimated delivery time, e.g. "2-3 days"
}

var ErrNoRatesAvailable = errors.New("no shipping rates available for this destination")

// Service is the abstraction checkout/order code depends on. Swap the
// concrete implementation (see MockProvider) for a real courier
// integration without changing any caller.
type Service interface {
	GetRates(ctx context.Context, req RateRequest) ([]RateOption, error)
}
