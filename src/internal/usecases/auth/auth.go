package auth

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"github.com/arieffadhlan/go-fitbyte/internal/config"
	"github.com/arieffadhlan/go-fitbyte/internal/dto"
	"github.com/arieffadhlan/go-fitbyte/internal/models"
	"github.com/arieffadhlan/go-fitbyte/internal/pkg/exceptions"
	"github.com/arieffadhlan/go-fitbyte/internal/pkg/security"
	authRepo "github.com/arieffadhlan/go-fitbyte/internal/repositories/auth"
)

// =======================================================================
//  SESSION STORE (IN-MEMORY)
// =======================================================================

type SessionStore interface {
	Set(key, value string)
	Get(key string) (string, bool)
	Delete(key string)
}

type InMemoryStore struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		data: make(map[string]string),
	}
}

func (s *InMemoryStore) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

func (s *InMemoryStore) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[key]
	return val, ok
}

func (s *InMemoryStore) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
}

// =======================================================================
//  JWT SECTION
// =======================================================================

type Claims struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateToken(id int, email, secret string, ttl time.Duration) (string, error) {
	claims := &Claims{
		ID:    id,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseToken(tokenStr, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if clm, ok := token.Claims.(*Claims); ok && token.Valid {
		return clm, nil
	}

	return nil, jwt.ErrTokenInvalidClaims
}

// =======================================================================
//  JWT MIDDLEWARE — CEK SESSION LOGIN
// =======================================================================

func JWTMiddleware(secret string, store SessionStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return unauth(c, "missing authorization header")
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return unauth(c, "invalid authorization header")
		}

		tokenStr := parts[1]

		claims, err := ParseToken(tokenStr, secret)
		if err != nil {
			return unauth(c, "invalid or expired token")
		}

		sessionKey := fmt.Sprintf("session:%d", claims.ID)
		savedToken, ok := store.Get(sessionKey)

		if !ok || savedToken != tokenStr {
			return unauth(c, "user not logged in")
		}

		c.Locals("id", claims.ID)
		c.Locals("email", claims.Email)

		return c.Next()
	}
}

func unauth(c *fiber.Ctx, msg string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": msg,
	})
}

// =======================================================================
//  AUTH USE CASE (LOGIN & REGISTER)
// =======================================================================

type AuthUseCase struct {
	authRepository authRepo.AuthRepositoryInterface
	cfg            *config.Config
	session        SessionStore
}

func NewAuthUseCase(
	authRepository authRepo.AuthRepositoryInterface,
	cfg *config.Config,
	session SessionStore,
) *AuthUseCase {
	return &AuthUseCase{
		authRepository: authRepository,
		cfg:            cfg,
		session:        session,
	}
}

func (uc *AuthUseCase) Register(ctx context.Context, req *dto.AuthRequest) (*dto.AuthResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	hashedPassword, err := security.HashingPassword(req.Password)
	if err != nil {
		return nil, exceptions.NewInternal("failed to hash password")
	}

	user := &models.User{
		Email:    req.Email,
		Password: hashedPassword,
	}

	id, err := uc.authRepository.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	token, err := GenerateToken(id, req.Email, uc.cfg.JwtSecret, 10*time.Second)
	if err != nil {
		return nil, exceptions.NewInternal("failed to generate tokens")
	}

	uc.session.Set(fmt.Sprintf("session:%d", id), token)

	return &dto.AuthResponse{
		Email: req.Email,
		Token: token,
	}, nil
}

func (uc *AuthUseCase) Login(ctx context.Context, req *dto.AuthRequest) (*dto.AuthResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	usr, err := uc.authRepository.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, exceptions.NewNotFound("user not found")
	}

	if usr == nil {
		return nil, exceptions.NewNotFound("invalid email or password")
	}

	if err := security.ComparePassword(req.Password, usr.Password); err != nil {
		return nil, exceptions.NewNotFound("invalid email or password")
	}

	token, err := GenerateToken(usr.ID, usr.Email, uc.cfg.JwtSecret, 10*time.Second)
	if err != nil {
		return nil, exceptions.NewInternal("failed to generate tokens")
	}

	uc.session.Set(fmt.Sprintf("session:%d", usr.ID), token)

	return &dto.AuthResponse{
		Email: usr.Email,
		Token: token,
	}, nil
}
