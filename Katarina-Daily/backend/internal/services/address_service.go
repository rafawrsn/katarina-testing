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

type AddressService struct {
	addresses *repositories.AddressRepository
	db        *gorm.DB
}

func NewAddressService(addresses *repositories.AddressRepository, db *gorm.DB) *AddressService {
	return &AddressService{addresses: addresses, db: db}
}

func (s *AddressService) List(ctx context.Context, userID uuid.UUID) ([]dto.AddressResponse, error) {
	items, err := s.addresses.ListByUser(ctx, userID)
	if err != nil {
		return nil, utils.ErrInternal(err)
	}
	return dto.ToAddressResponseList(items), nil
}

func (s *AddressService) Create(ctx context.Context, userID uuid.UUID, req dto.CreateAddressRequest) (*dto.AddressResponse, error) {
	addr := &models.Address{
		UserID:        userID,
		Label:         req.Label,
		RecipientName: req.RecipientName,
		Phone:         req.Phone,
		AddressLine:   req.Address,
		Province:      req.Province,
		City:          req.City,
		District:      req.District,
		PostalCode:    req.PostalCode,
		IsDefault:     req.IsDefault,
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		txRepo := s.addresses.WithTx(tx)

		// First address for a user is always the default, regardless
		// of what was requested.
		existing, listErr := s.countAddresses(ctx, tx, userID)
		if listErr != nil {
			return listErr
		}
		if existing == 0 {
			addr.IsDefault = true
		}

		if err := tx.WithContext(ctx).Create(addr).Error; err != nil {
			return err
		}
		if addr.IsDefault {
			if err := txRepo.ClearDefaultForUser(ctx, tx, userID, addr.ID); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, utils.ErrInternal(err)
	}

	resp := dto.ToAddressResponse(addr)
	return &resp, nil
}

func (s *AddressService) countAddresses(ctx context.Context, tx *gorm.DB, userID uuid.UUID) (int64, error) {
	var count int64
	err := tx.WithContext(ctx).Model(&models.Address{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

func (s *AddressService) Update(ctx context.Context, userID, addressID uuid.UUID, req dto.UpdateAddressRequest) (*dto.AddressResponse, error) {
	addr, err := s.getOwned(ctx, userID, addressID)
	if err != nil {
		return nil, err
	}

	if req.Label != nil {
		addr.Label = *req.Label
	}
	if req.RecipientName != nil {
		addr.RecipientName = *req.RecipientName
	}
	if req.Phone != nil {
		addr.Phone = *req.Phone
	}
	if req.Address != nil {
		addr.AddressLine = *req.Address
	}
	if req.Province != nil {
		addr.Province = *req.Province
	}
	if req.City != nil {
		addr.City = *req.City
	}
	if req.District != nil {
		addr.District = *req.District
	}
	if req.PostalCode != nil {
		addr.PostalCode = *req.PostalCode
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if req.IsDefault != nil && *req.IsDefault {
			addr.IsDefault = true
			if err := tx.WithContext(ctx).Save(addr).Error; err != nil {
				return err
			}
			return s.addresses.WithTx(tx).ClearDefaultForUser(ctx, tx, userID, addr.ID)
		}
		return tx.WithContext(ctx).Save(addr).Error
	})
	if err != nil {
		return nil, utils.ErrInternal(err)
	}

	resp := dto.ToAddressResponse(addr)
	return &resp, nil
}

func (s *AddressService) Delete(ctx context.Context, userID, addressID uuid.UUID) error {
	addr, err := s.getOwned(ctx, userID, addressID)
	if err != nil {
		return err
	}
	if delErr := s.addresses.Delete(ctx, addr.ID); delErr != nil {
		return utils.ErrInternal(delErr)
	}
	return nil
}

func (s *AddressService) SetDefault(ctx context.Context, userID, addressID uuid.UUID) (*dto.AddressResponse, error) {
	addr, err := s.getOwned(ctx, userID, addressID)
	if err != nil {
		return nil, err
	}

	addr.IsDefault = true
	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Save(addr).Error; err != nil {
			return err
		}
		return s.addresses.WithTx(tx).ClearDefaultForUser(ctx, tx, userID, addr.ID)
	})
	if txErr != nil {
		return nil, utils.ErrInternal(txErr)
	}

	resp := dto.ToAddressResponse(addr)
	return &resp, nil
}

// getOwned loads an address and verifies it belongs to userID, so
// customers can never read/modify another customer's address.
func (s *AddressService) getOwned(ctx context.Context, userID, addressID uuid.UUID) (*models.Address, error) {
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
	return addr, nil
}
