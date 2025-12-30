package jwt

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
)

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

func Middleware(secret string, db *sqlx.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing authorization header",
			})
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || (strings.ToLower(parts[0]) != "bearer" && strings.ToLower(parts[0]) != "Bearer") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid authorization header",
			})
		}

		token := parts[1]
		claim, err := ParseToken(token, secret)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid or expired token",
			})
		}

		// var userID int
		// query := `SELECT id FROM users WHERE id = $1`
		// err = db.QueryRowContext(context.Background(), query, claim.ID).Scan(&userID)
		// if err != nil {
		// 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		// 		"error": "user not found or not registered",
		// 	})
		// }

		c.Locals("id", claim.ID)
		c.Locals("email", claim.Email)

		return c.Next()
	}
}
