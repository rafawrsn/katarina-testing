package shipping

import (
	"context"
	"errors"
	"hash/fnv"
	"math"
	"strings"
)

// MockProvider computes dynamic, formula-based shipping quotes from
// origin, destination, weight, courier and service - satisfying the
// same interface a real integration would.
//
// It is NOT a real courier integration. It exists so the rest of the
// system (checkout, order snapshots, admin shipment tooling) has a
// fully working shipping cost pipeline to build and test against today.
// To go live, implement Service against a real Indonesian shipping
// aggregator (e.g. RajaOngkir, Biteship, Shipper) using SHIPPING_API_KEY
// and swap the constructor in main.go - no other code changes needed.
type MockProvider struct{}

func NewMockProvider() *MockProvider { return &MockProvider{} }

type courierService struct {
	Courier     string
	CourierName string
	Service     string
	Description string
	BaseFee     float64 // flat starting fee, IDR
	PerKg       float64 // cost per kg (rounded up), IDR
	ETD         string
}

var catalog = []courierService{
	{"jne", "JNE", "REG", "Layanan Reguler", 9000, 8000, "2-3 hari"},
	{"jne", "JNE", "YES", "Yakin Esok Sampai", 18000, 15000, "1 hari"},
	{"jnt", "J&T Express", "EZ", "Ekonomis", 8000, 7500, "2-4 hari"},
	{"jnt", "J&T Express", "REG", "Reguler", 10000, 9000, "1-3 hari"},
	{"sicepat", "SiCepat", "REG", "Reguler", 9000, 8500, "2-3 hari"},
	{"sicepat", "SiCepat", "BEST", "Besok Sampai Tujuan", 20000, 16000, "1 hari"},
	{"anteraja", "AnterAja", "REG", "Reguler", 8500, 8000, "2-3 hari"},
}

func (m *MockProvider) GetRates(ctx context.Context, req RateRequest) ([]RateOption, error) {
	if req.WeightGrams <= 0 {
		return nil, errors.New("weight must be greater than zero")
	}

	weightKg := math.Ceil(float64(req.WeightGrams) / 1000.0)
	if weightKg < 1 {
		weightKg = 1
	}
	zone := zoneMultiplier(req.Origin, req.Destination)

	out := make([]RateOption, 0, len(catalog))
	for _, cs := range catalog {
		if req.Courier != "" && !strings.EqualFold(req.Courier, cs.Courier) {
			continue
		}
		cost := (cs.BaseFee + cs.PerKg*weightKg) * zone
		roundedCost := int64(math.Round(cost/500) * 500) // round to nearest Rp 500

		out = append(out, RateOption{
			Courier:     cs.Courier,
			CourierName: cs.CourierName,
			Service:     cs.Service,
			Description: cs.Description,
			Cost:        roundedCost,
			ETD:         cs.ETD,
		})
	}

	if len(out) == 0 {
		return nil, ErrNoRatesAvailable
	}
	return out, nil
}

// zoneMultiplier gives same-city shipments the cheapest rate,
// same-province a moderate rate, and other provinces a higher rate that
// varies deterministically by destination so quotes aren't identical
// for every out-of-province address (a stand-in for real distance).
func zoneMultiplier(origin, destination Location) float64 {
	switch {
	case strings.EqualFold(strings.TrimSpace(origin.City), strings.TrimSpace(destination.City)):
		return 1.0
	case strings.EqualFold(strings.TrimSpace(origin.Province), strings.TrimSpace(destination.Province)):
		return 1.3
	default:
		return 1.8 + pseudoDistanceFactor(destination.Province)
	}
}

// pseudoDistanceFactor deterministically maps a province name to a
// value in [0.0, 0.6) so different provinces get different (but
// stable) inter-province multipliers.
func pseudoDistanceFactor(province string) float64 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(strings.ToLower(strings.TrimSpace(province))))
	return float64(h.Sum32()%600) / 1000.0
}
