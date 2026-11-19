package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"fashion-ecommerce-backend/internal/config"
	"fashion-ecommerce-backend/internal/dto"
	"fashion-ecommerce-backend/internal/repositories"
	"fashion-ecommerce-backend/internal/utils"
	"fashion-ecommerce-backend/pkg/shipping"
)

type ShippingService struct {
	addresses *repositories.AddressRepository
	carts     *repositories.CartRepository
	provider  shipping.Service
	cfg       *config.Config
}

func NewShippingService(addresses *repositories.AddressRepository, carts *repositories.CartRepository, provider shipping.Service, cfg *config.Config) *ShippingService {
	return &ShippingService{addresses: addresses, carts: carts, provider: provider, cfg: cfg}
}

// GetRatesForCart quotes shipping cost for the customer's current cart
// contents delivered to the given address. Weight is always computed
// server-side from live variant weights - never trusted from a client.
func (s *ShippingService) GetRatesForCart(ctx context.Context, userID, addressID uuid.UUID, courier string) ([]dto.ShippingRateOption, error) {
	addr, err := s.addresses.FindByID(ctx, addressID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("address")
		}
		return nil, utils.ErrInternal(err)
	}
	if addr.UserID != userID {
		return nil, utils.ErrForbidden
	}

	cart, err := s.carts.GetOrCreateCart(ctx, userID)
	if err != nil {
		return nil, utils.ErrInternal(err)
	}
	full, err := s.carts.GetCartWithItems(ctx, cart.ID)
	if err != nil {
		return nil, utils.ErrInternal(err)
	}
	if len(full.Items) == 0 {
		return nil, utils.ErrValidation("cart is empty")
	}

	totalWeight := 0
	for _, item := range full.Items {
		totalWeight += item.ProductVariant.WeightGrams * item.Quantity
	}

	rates, err := s.provider.GetRates(ctx, shipping.RateRequest{
		Origin: shipping.Location{
			Province: s.cfg.ShippingOriginProv,
			City:     s.cfg.ShippingOriginCity,
		},
		Destination: shipping.Location{
			Province: addr.Province,
			City:     addr.City,
			District: addr.District,
			Postal:   addr.PostalCode,
		},
		WeightGrams: totalWeight,
		Courier:     courier,
	})
	if err != nil {
		if errors.Is(err, shipping.ErrNoRatesAvailable) {
			return nil, utils.NewAppError(404, "NO_SHIPPING_RATES", "no shipping rates available for this destination")
		}
		return nil, utils.ErrInternal(err)
	}

	out := make([]dto.ShippingRateOption, 0, len(rates))
	for _, r := range rates {
		out = append(out, dto.ShippingRateOption{
			Courier:     r.Courier,
			Service:     r.Service,
			Description: r.Description,
			Cost:        decimal.NewFromInt(r.Cost),
			ETD:         r.ETD,
		})
	}
	return out, nil
}

// QuoteForCheckout re-quotes a single courier+service for the checkout
// service, which never trusts a client-supplied shipping cost.
func (s *ShippingService) QuoteForCheckout(ctx context.Context, destProvince, destCity, destDistrict, destPostal string, weightGrams int, courier, service string) (decimal.Decimal, error) {
	rates, err := s.provider.GetRates(ctx, shipping.RateRequest{
		Origin: shipping.Location{
			Province: s.cfg.ShippingOriginProv,
			City:     s.cfg.ShippingOriginCity,
		},
		Destination: shipping.Location{
			Province: destProvince,
			City:     destCity,
			District: destDistrict,
			Postal:   destPostal,
		},
		WeightGrams: weightGrams,
		Courier:     courier,
	})
	if err != nil {
		return decimal.Zero, err
	}
	for _, r := range rates {
		if r.Service == service {
			return decimal.NewFromInt(r.Cost), nil
		}
	}
	return decimal.Zero, utils.ErrValidation("selected courier/service is not available for this destination")
}
