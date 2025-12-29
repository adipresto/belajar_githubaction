package activity

import (
	"context"
	"database/sql"
	"strings"

	"github.com/arieffadhlan/go-fitbyte/internal/dto"
	"github.com/arieffadhlan/go-fitbyte/internal/models"
	"github.com/arieffadhlan/go-fitbyte/internal/pkg/exceptions"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func NewActivityRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetById(ctx context.Context, id int) (*models.Activity, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	query := `
		SELECT id, user_id, activity_type, done_at, duration_in_minutes, calories_burned, created_at, is_active FROM activities WHERE id = $1 AND is_active <> B'0'
	`

	activity := &models.Activity{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(&activity.ID, &activity.UserId, &activity.ActivityType, &activity.DoneAt, &activity.DurationInMin, &activity.CaloriesBurned, &activity.CreatedAt, &activity.IsActive)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exceptions.ErrNotFound
		}
		return nil, err
	}

	return activity, nil
}

func (r *repository) GetAll(ctx context.Context, queries dto.ActivityQueryParamRequest) ([]*models.Activity, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	baseQuery := `
        SELECT 
            id, user_id, activity_type, done_at, duration_in_minutes, 
            calories_burned, created_at, updated_at
        FROM activities
    `

	params := map[string]interface{}{}
	conditions := []string{}

	conditions = append(conditions, "is_active <> B'0'")

	if queries.UserId != 0 {
		conditions = append(conditions, "user_id = :user_id")
		params["user_id"] = queries.UserId
	}

	if queries.ActivityType != "" {
		conditions = append(conditions, "activity_type = :activity_type")
		params["activity_type"] = queries.ActivityType
	}

	if queries.DoneAtForm != "" {
		conditions = append(conditions, "done_at >= :done_at_from")
		params["done_at_from"] = queries.DoneAtForm
	}

	if queries.DoneAtTo != "" {
		conditions = append(conditions, "done_at <= :done_at_to")
		params["done_at_to"] = queries.DoneAtTo
	}

	if queries.CaloriesBurnedMin > 0 {
		conditions = append(conditions, "calories_burned >= :cal_min")
		params["cal_min"] = queries.CaloriesBurnedMin
	}

	if queries.CaloriesBurnedMax > 0 {
		conditions = append(conditions, "calories_burned <= :cal_max")
		params["cal_max"] = queries.CaloriesBurnedMax
	}

	if len(conditions) > 0 {
		baseQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	baseQuery += " ORDER BY done_at DESC LIMIT :limit OFFSET :offset"
	params["limit"] = queries.Limit
	params["offset"] = queries.Offset

	activities := []*models.Activity{}

	stmt, err := r.db.PrepareNamedContext(ctx, baseQuery)
	if err != nil {
		return nil, err
	}

	if err := stmt.SelectContext(ctx, &activities, params); err != nil {
		return nil, err
	}

	return activities, nil
}

func (r *repository) Post(ctx context.Context, activity *models.Activity) (*models.Activity, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	query := `
		INSERT INTO activities
		(user_id, activity_type, done_at, duration_in_minutes, calories_burned)
		VALUES($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		activity.UserId,
		activity.ActivityType,
		activity.DoneAt,
		activity.DurationInMin,
		activity.CaloriesBurned,
	).Scan(&activity.ID, &activity.CreatedAt, &activity.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return activity, nil
}

func (r *repository) Update(ctx context.Context, updatedActivity *models.Activity) (*models.Activity, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	bitBool := byte(1)
	if updatedActivity.IsActive == false {
		bitBool = byte(0)
	}

	query := `
		UPDATE activities
		SET activity_type = $1, done_at = $2, duration_in_minutes = $3, calories_burned = $4, updated_at = now(), is_active = $7
		WHERE id= $5 and user_id = $6
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx, query, updatedActivity.ActivityType, updatedActivity.DoneAt, updatedActivity.DurationInMin, updatedActivity.CaloriesBurned, updatedActivity.ID, updatedActivity.UserId, bitBool).Scan(&updatedActivity.ID, &updatedActivity.CreatedAt, &updatedActivity.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exceptions.ErrNotFound
		}
		return nil, err
	}

	return updatedActivity, nil
}
