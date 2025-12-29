package dto

import (
	"time"

	"github.com/arieffadhlan/go-fitbyte/internal/models"
)

type ActivityRequest struct {
	ActivityType      models.ActivityType `json:"activityType" validate:"required,oneof=Walking Yoga Stretching Cycling Swimming Dancing Hiking Running HIIT JumpRope"`
	DoneAt            string              `json:"doneAt" validate:"required,iso8601"`
	DurationInMinutes int                 `json:"durationInMinutes" validate:"required,gte=1"`
}

type ActivityUpdateRequest struct {
	ActivityType      *models.ActivityType `json:"activityType,omitempty" validate:"omitempty,oneof=Walking Yoga Stretching Cycling Swimming Dancing Hiking Running HIIT JumpRope"`
	DoneAt            *string              `json:"doneAt,omitempty" validate:"omitempty,iso8601"`
	DurationInMinutes *int                 `json:"durationInMinutes,omitempty" validate:"omitempty,gte=1"`
}

type ActivityQueryParamRequest struct {
	Limit             int    `form:"limit" binding:"gte=0"`
	Offset            int    `form:"offset" binding:"gte=0"`
	ActivityType      string `form:"activityType" binding:"omitempty,oneof=Walking Yoga Stretching Cycling Swimming Dancing Hiking Running HIIT JumpRope"`
	DoneAtForm        string `form:"doneAtForm"`
	DoneAtTo          string `form:"doneAtTo"`
	CaloriesBurnedMin int    `form:"caloriesBurnedMin" binding:"gte=0"`
	CaloriesBurnedMax int    `form:"caloriesBurnedMax" binding:"gte=0"`
	UserId            int    `form:"userId"`
}

type ActivityPayload struct {
	Limit             int
	Offset            int
	ActivityType      string
	DoneAtForm        time.Time
	DoneAtTo          time.Time
	CaloriesBurnedMin int
	CaloriesBurnedMax int
}

type ActivityResponse struct {
	ActivityId        string              `json:"activityId"`
	ActivityType      models.ActivityType `json:"activityType"`
	DoneAt            string              `json:"doneAt"`
	DurationInMinutes int                 `json:"durationInMinutes"`
	CaloriesBurned    int                 `json:"caloriesBurned"`
	CreatedAt         string              `json:"createdAt"`
	UpdatedAt         string              `json:"updatedAt"`
}
