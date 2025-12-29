package profile

import (
	"context"

	"github.com/arieffadhlan/go-fitbyte/internal/dto"
)

type ProfileUseCaseInterface interface {
	GetProfile(ctx context.Context, userID int) (*dto.ProfileResponse, error)
	UpdateProfile(ctx context.Context, userID int, req *dto.ProfileUpdateRequest) (*dto.ProfileResponse, error)
}
