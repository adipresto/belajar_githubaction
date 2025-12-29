package auth

import (
	"context"
	"database/sql"

	"github.com/arieffadhlan/go-fitbyte/internal/models"
	"github.com/arieffadhlan/go-fitbyte/internal/pkg/exceptions"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type AuthRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) AuthRepositoryInterface {
	return &AuthRepository{db: db}
}

func (ar *AuthRepository) Create(ctx context.Context, user *models.User) (int, error) {
	var newId int

	query := `
		INSERT INTO users (email, password)
		VALUES ($1, $2)
		RETURNING id`

	err := ar.db.QueryRowContext(ctx, query,
		user.Email,
		user.Password,
	).Scan(&newId)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return 0, exceptions.NewConflict("email already exists")
		} else {
			return 0, exceptions.NewInternal("register user failed")
		}
	}

	return newId, nil
}

func (ar *AuthRepository) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	uData := &models.User{}
	query := `SELECT id, email, password FROM users WHERE email = $1`

	err := ar.db.QueryRowContext(ctx,
		query,
		email,
	).Scan(&uData.ID, &uData.Email, &uData.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exceptions.NewNotFound("user data not found")
		} else {
			return nil, err
		}
	}

	return uData, nil
}
