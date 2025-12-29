package profile

import (
	"context"

	"github.com/arieffadhlan/go-fitbyte/internal/models"
)

type ProfileRepositoryInterface interface {
	GetUserByID(ctx context.Context, userID int) (*models.User, error)
	UpdateProfile(ctx context.Context, userID int, updates map[string]interface{}) error
}
