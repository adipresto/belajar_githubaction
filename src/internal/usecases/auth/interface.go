package auth

import (
	"context"

	"github.com/arieffadhlan/go-fitbyte/internal/dto"
)

type AuthUseCaseInterface interface {
	Login(context.Context, *dto.AuthRequest) (*dto.AuthResponse, error)
	Register(context.Context, *dto.AuthRequest) (*dto.AuthResponse, error)
}
