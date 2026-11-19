package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"fashion-ecommerce-backend/internal/dto"
	"fashion-ecommerce-backend/internal/models"
	"fashion-ecommerce-backend/internal/repositories"
	"fashion-ecommerce-backend/internal/utils"
)

type CartService struct {
	carts    *repositories.CartRepository
	variants *repositories.VariantRepository
}

func NewCartService(carts *repositories.CartRepository, variants *repositories.VariantRepository) *CartService {
	return &CartService{carts: carts, variants: variants}
}

func (s *CartService) GetCart(ctx context.Context, userID uuid.UUID) (*dto.CartResponse, error) {
	cart, err := s.carts.GetOrCreateCart(ctx, userID)
	if err != nil {
		return nil, utils.ErrInternal(err)
	}
	full, err := s.carts.GetCartWithItems(ctx, cart.ID)
	if err != nil {
		return nil, utils.ErrInternal(err)
	}
	resp := dto.ToCartResponse(full)
	return &resp, nil
}

// AddItem implements the exact validation order from the spec:
// 1. variant exists, 2. variant active, 3. quantity > 0 (enforced by
// binding), 4. enough stock, 5. increment if already in cart,
// 6. never allow quantity above available stock.
func (s *CartService) AddItem(ctx context.Context, userID uuid.UUID, req dto.AddCartItemRequest) (*dto.CartResponse, error) {
	variant, err := s.variants.FindByID(ctx, req.ProductVariantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("product_variant")
		}
		return nil, utils.ErrInternal(err)
	}
	if !variant.IsActive() {
		return nil, utils.ErrValidation("this product variant is not currently available")
	}

	cart, err := s.carts.GetOrCreateCart(ctx, userID)
	if err != nil {
		return nil, utils.ErrInternal(err)
	}

	existing, err := s.carts.FindItemByVariant(ctx, cart.ID, variant.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, utils.ErrInternal(err)
	}

	desiredQty := req.Quantity
	if existing != nil {
		desiredQty += existing.Quantity
	}
	if desiredQty > variant.Stock {
		return nil, utils.ErrValidation("requested quantity exceeds available stock")
	}

	if existing != nil {
		existing.Quantity = desiredQty
		if err := s.carts.UpdateItem(ctx, existing); err != nil {
			return nil, utils.ErrInternal(err)
		}
	} else {
		item := &models.CartItem{
			CartID:           cart.ID,
			ProductVariantID: variant.ID,
			Quantity:         req.Quantity,
		}
		if err := s.carts.CreateItem(ctx, item); err != nil {
			return nil, utils.ErrInternal(err)
		}
	}

	return s.GetCart(ctx, userID)
}

func (s *CartService) UpdateItem(ctx context.Context, userID, itemID uuid.UUID, req dto.UpdateCartItemRequest) (*dto.CartResponse, error) {
	item, err := s.getOwnedItem(ctx, userID, itemID)
	if err != nil {
		return nil, err
	}

	variant, err := s.variants.FindByID(ctx, item.ProductVariantID)
	if err != nil {
		return nil, utils.ErrInternal(err)
	}
	if req.Quantity > variant.Stock {
		return nil, utils.ErrValidation("requested quantity exceeds available stock")
	}

	item.Quantity = req.Quantity
	if err := s.carts.UpdateItem(ctx, item); err != nil {
		return nil, utils.ErrInternal(err)
	}

	return s.GetCart(ctx, userID)
}

func (s *CartService) DeleteItem(ctx context.Context, userID, itemID uuid.UUID) error {
	item, err := s.getOwnedItem(ctx, userID, itemID)
	if err != nil {
		return err
	}
	if err := s.carts.DeleteItem(ctx, item.ID); err != nil {
		return utils.ErrInternal(err)
	}
	return nil
}

func (s *CartService) ClearCart(ctx context.Context, userID uuid.UUID) error {
	cart, err := s.carts.GetOrCreateCart(ctx, userID)
	if err != nil {
		return utils.ErrInternal(err)
	}
	if err := s.carts.ClearItems(ctx, cart.ID); err != nil {
		return utils.ErrInternal(err)
	}
	return nil
}

// getOwnedItem loads a cart item and verifies it belongs to a cart
// owned by userID, so customers can never mutate another customer's
// cart via a guessed item id.
// Uses FindCartByUserID (read-only) instead of GetOrCreateCart to avoid
// a superfluous cart-creation side-effect during an ownership check.
func (s *CartService) getOwnedItem(ctx context.Context, userID, itemID uuid.UUID) (*models.CartItem, error) {
	item, err := s.carts.FindItemByID(ctx, itemID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("cart_item")
		}
		return nil, utils.ErrInternal(err)
	}

	cart, err := s.carts.FindCartByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// User has no cart at all → item can't belong to them.
			return nil, utils.ErrForbidden
		}
		return nil, utils.ErrInternal(err)
	}
	if item.CartID != cart.ID {
		return nil, utils.ErrForbidden
	}
	return item, nil
}
