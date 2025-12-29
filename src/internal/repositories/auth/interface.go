package auth

import (
	"context"

	"github.com/arieffadhlan/go-fitbyte/internal/models"
)

type AuthRepositoryInterface interface {
	Create(context.Context, *models.User) (int, error)
	FindUserByEmail(context.Context, string) (*models.User, error)
}
