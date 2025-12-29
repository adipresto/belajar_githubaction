package activity

import (
	"context"

	"github.com/arieffadhlan/go-fitbyte/internal/dto"
	"github.com/arieffadhlan/go-fitbyte/internal/models"
)

type Repository interface {
	GetAll(c context.Context, queries dto.ActivityQueryParamRequest) ([]*models.Activity, error)
	GetById(c context.Context, id int) (*models.Activity, error)
	Post(c context.Context, activity *models.Activity) (*models.Activity, error)
	Update(c context.Context, updatedActivity *models.Activity) (*models.Activity, error)
}
