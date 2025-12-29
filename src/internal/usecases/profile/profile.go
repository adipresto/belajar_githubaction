package profile

import (
	"context"
	"fmt"

	"github.com/arieffadhlan/go-fitbyte/internal/dto"
	"github.com/arieffadhlan/go-fitbyte/internal/models"
	"github.com/arieffadhlan/go-fitbyte/internal/repositories/profile"
)

type ProfileUseCase struct {
	profileRepo profile.ProfileRepositoryInterface
}

func NewProfileUseCase(profileRepo profile.ProfileRepositoryInterface) ProfileUseCaseInterface {
	return &ProfileUseCase{
		profileRepo: profileRepo,
	}
}

func (uc *ProfileUseCase) GetProfile(ctx context.Context, userID int) (*dto.ProfileResponse, error) {
	user, err := uc.profileRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	var response dto.ProfileResponse
	response.FromUser(user)

	return &response, nil
}

func (uc *ProfileUseCase) UpdateProfile(ctx context.Context, userID int, req *dto.ProfileUpdateRequest) (*dto.ProfileResponse, error) {

	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	updates := make(map[string]interface{})

	if req.Preference != nil {
		updates["preference"] = models.PreferenceType(*req.Preference)
	}

	if req.WeightUnit != nil {
		updates["weight_unit"] = models.WeightUnitType(*req.WeightUnit)
	}

	if req.HeightUnit != nil {
		updates["height_unit"] = models.HeightUnitType(*req.HeightUnit)
	}

	if req.Weight != nil {
		updates["weight"] = *req.Weight
	}

	if req.Height != nil {
		updates["height"] = *req.Height
	}

	if req.Name != nil {
		updates["name"] = *req.Name
	}

	if req.ImageURI != nil {
		updates["image_uri"] = *req.ImageURI
	}

	if len(updates) > 0 {
		err := uc.profileRepo.UpdateProfile(ctx, userID, updates)
		if err != nil {
			return nil, fmt.Errorf("failed to update profile: %w", err)
		}
	}

	return uc.GetProfile(ctx, userID)
}
