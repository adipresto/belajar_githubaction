package activity

import (
	"context"
	"strconv"

	"github.com/arieffadhlan/go-fitbyte/internal/dto"
	"github.com/arieffadhlan/go-fitbyte/internal/models"
	"github.com/arieffadhlan/go-fitbyte/internal/pkg/exceptions"
	"github.com/arieffadhlan/go-fitbyte/internal/repositories/activity"
)

type useCase struct {
	activityRepository activity.Repository
}

func NewUseCase(activityRepository activity.Repository) UseCase {
	return &useCase{
		activityRepository,
	}
}

func (u *useCase) GetAllActivities(ctx context.Context, queries dto.ActivityQueryParamRequest) ([]*dto.ActivityResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	activities, err := u.activityRepository.GetAll(ctx, queries)
	if err != nil {
		return nil, err
	}

	// convert activities to dto.ActivityResponse
	activitiesResponse := make([]*dto.ActivityResponse, len(activities))
	for i, activity := range activities {
		activitiesResponse[i] = &dto.ActivityResponse{
			ActivityId:        strconv.Itoa(activity.ID),
			ActivityType:      activity.ActivityType,
			DoneAt:            activity.DoneAt,
			DurationInMinutes: activity.DurationInMin,
			CaloriesBurned:    activity.CaloriesBurned,
			CreatedAt:         activity.CreatedAt,
			UpdatedAt:         activity.UpdatedAt,
		}
	}
	return activitiesResponse, nil
}

func (u *useCase) GetActivityById(ctx context.Context, id int) (*dto.ActivityResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	activity, err := u.activityRepository.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dto.ActivityResponse{
		ActivityId:        strconv.Itoa(activity.ID),
		ActivityType:      activity.ActivityType,
		DoneAt:            activity.DoneAt,
		DurationInMinutes: activity.DurationInMin,
		CaloriesBurned:    activity.CaloriesBurned,
	}, nil
}

func (u *useCase) PostActivity(ctx context.Context, activity *dto.ActivityRequest, userId int) (*dto.ActivityResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	newActivity := &models.Activity{
		UserId:         userId,
		ActivityType:   activity.ActivityType,
		DoneAt:         activity.DoneAt,
		DurationInMin:  activity.DurationInMinutes,
		CaloriesBurned: activity.ActivityType.GetTotalCalories(activity.DurationInMinutes),
	}

	res, err := u.activityRepository.Post(ctx, newActivity)
	if err != nil {
		return nil, err
	}

	activityId := strconv.Itoa(res.ID)

	return &dto.ActivityResponse{
		ActivityId:        activityId,
		ActivityType:      res.ActivityType,
		DoneAt:            res.DoneAt,
		DurationInMinutes: res.DurationInMin,
		CaloriesBurned:    res.CaloriesBurned,
		CreatedAt:         res.CreatedAt,
		UpdatedAt:         res.UpdatedAt,
	}, nil
}

func (u *useCase) DisableActivity(ctx context.Context, userID int, activityID string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if activityID == "" {
		return exceptions.ErrNotFound
	}

	activityIdInt, err := strconv.Atoi(activityID)
	if err != nil {
		return err
	}

	// Fetch only if we need fallback values
	var existing *models.Activity
	existing, err = u.activityRepository.GetById(ctx, activityIdInt)
	if err != nil {
		if err == exceptions.ErrNotFound {
			return exceptions.ErrNotFound
		}
		return err
	}

	updatePayloadActivity := &models.Activity{
		ActivityType:   existing.ActivityType,
		DoneAt:         existing.DoneAt,
		DurationInMin:  existing.DurationInMin,
		CaloriesBurned: existing.ActivityType.GetTotalCalories(existing.DurationInMin),
		ID:             activityIdInt,
		UserId:         userID,
		IsActive:       true,
	}

	_, err = u.activityRepository.Update(ctx, updatePayloadActivity)
	if err != nil {
		if err == exceptions.ErrNotFound {
			return exceptions.ErrNotFound
		}
		return err
	}

	return nil
}

func (u *useCase) UpdateActivity(ctx context.Context, updatedActivity *dto.ActivityUpdateRequest, userID int, activityID string) (*dto.ActivityResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if activityID == "" {
		return nil, exceptions.ErrNotFound
	}

	activityIdInt, err := strconv.Atoi(activityID)
	if err != nil {
		return nil, err
	}

	// Fetch only if we need fallback values
	var existing *models.Activity
	if updatedActivity.ActivityType == nil || updatedActivity.DurationInMinutes == nil || updatedActivity.DoneAt == nil {
		existing, err = u.activityRepository.GetById(ctx, activityIdInt)
		if err != nil {
			if err == exceptions.ErrNotFound {
				return nil, exceptions.ErrNotFound
			}
			return nil, err
		}
	}

	// Merge final values
	var finalActivityType models.ActivityType
	if updatedActivity.ActivityType != nil {
		finalActivityType = *updatedActivity.ActivityType
	} else {
		finalActivityType = existing.ActivityType
	}

	var finalDuration int
	if updatedActivity.DurationInMinutes != nil {
		finalDuration = *updatedActivity.DurationInMinutes
	} else {
		finalDuration = existing.DurationInMin
	}

	var finalDoneAt string
	if updatedActivity.DoneAt != nil {
		finalDoneAt = *updatedActivity.DoneAt
	} else {
		finalDoneAt = existing.DoneAt
	}

	updatePayloadActivity := &models.Activity{
		ActivityType:   finalActivityType,
		DoneAt:         finalDoneAt,
		DurationInMin:  finalDuration,
		CaloriesBurned: finalActivityType.GetTotalCalories(finalDuration),
		ID:             activityIdInt,
		UserId:         userID,
		IsActive:       true,
	}

	res, err := u.activityRepository.Update(ctx, updatePayloadActivity)
	if err != nil {
		if err == exceptions.ErrNotFound {
			return nil, exceptions.ErrNotFound
		}
		return nil, err
	}

	activityId := strconv.Itoa(res.ID)

	return &dto.ActivityResponse{
		ActivityId:        activityId,
		ActivityType:      res.ActivityType,
		DoneAt:            res.DoneAt,
		DurationInMinutes: res.DurationInMin,
		CaloriesBurned:    res.CaloriesBurned,
		CreatedAt:         res.CreatedAt,
		UpdatedAt:         res.UpdatedAt,
	}, nil
}
