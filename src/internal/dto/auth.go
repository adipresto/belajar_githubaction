package dto

import (
	"regexp"
	
	"github.com/arieffadhlan/go-fitbyte/internal/pkg/exceptions"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func (a *AuthRequest) Validate() error {
	if a.Email == "" || a.Password == "" {
		return exceptions.NewBadRequest("email and password are required")
	}

	if !emailRegex.MatchString(a.Email) {
		return exceptions.NewBadRequest("invalid email format")
	}

	passLen := len(a.Password)
	if passLen < 8 || passLen > 32 {
		return exceptions.NewBadRequest("password must be between 8 and 32 characters")
	}
	return nil
}
