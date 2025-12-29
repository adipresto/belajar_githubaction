package activity

import (
	"context"

	"github.com/arieffadhlan/go-fitbyte/internal/dto"
)

type UseCase interface {
	GetAllActivities(context.Context, dto.ActivityQueryParamRequest) ([]*dto.ActivityResponse, error)
	GetActivityById(context.Context, int) (*dto.ActivityResponse, error)
	PostActivity(context.Context, *dto.ActivityRequest, int) (*dto.ActivityResponse, error)
	UpdateActivity(ctx context.Context, updatedActivity *dto.ActivityUpdateRequest, userID int, activityID string) (*dto.ActivityResponse, error)
	DisableActivity(ctx context.Context, userID int, activityID string) error
}
