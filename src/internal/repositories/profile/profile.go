package profile

import (
	"context"
	"fmt"
	"strings"

	"github.com/arieffadhlan/go-fitbyte/internal/models"
	"github.com/jmoiron/sqlx"
)

type ProfileRepository struct {
	db *sqlx.DB
}

func NewProfileRepository(db *sqlx.DB) ProfileRepositoryInterface {
	return &ProfileRepository{
		db: db,
	}
}

func (r *ProfileRepository) GetUserByID(ctx context.Context, userID int) (*models.User, error) {
	query := `
		SELECT id, email, password, name, preference, weight_unit, height_unit, 
		       weight, height, image_uri, created_at, updated_at
		FROM users 
		WHERE id = $1
	`

	var user models.User
	err := r.db.GetContext(ctx, &user, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return &user, nil
}

func (r *ProfileRepository) UpdateProfile(ctx context.Context, userID int, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return fmt.Errorf("no fields to update")
	}

	setParts := make([]string, 0, len(updates))
	args := make([]interface{}, 0, len(updates)+1)
	argIndex := 1

	for field, value := range updates {
		setParts = append(setParts, fmt.Sprintf("%s = $%d", field, argIndex))
		args = append(args, value)
		argIndex++
	}

	setParts = append(setParts, fmt.Sprintf("updated_at = CURRENT_TIMESTAMP"))

	query := fmt.Sprintf(`
		UPDATE users 
		SET %s 
		WHERE id = $%d
	`, strings.Join(setParts, ", "), argIndex)

	args = append(args, userID)

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}

	return nil
}
