package handlers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/arieffadhlan/go-fitbyte/internal/dto"
	"github.com/arieffadhlan/go-fitbyte/internal/pkg/exceptions"
	internal_validator "github.com/arieffadhlan/go-fitbyte/internal/pkg/validator"
	"github.com/arieffadhlan/go-fitbyte/internal/usecases/activity"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ActivityHandler interface {
	Post(fibCtx *fiber.Ctx) error
	Update(fibCtx *fiber.Ctx) error
	GetAll(fibCtx *fiber.Ctx) error
	GetById(fibCtx *fiber.Ctx) error
}

type activityHandler struct {
	activityUseCase activity.UseCase
}

func NewActivityHandler(activityUseCase activity.UseCase) ActivityHandler {
	return &activityHandler{activityUseCase}
}

func (r *activityHandler) GetAll(fibCtx *fiber.Ctx) error {
	var dto dto.ActivityQueryParamRequest
	_ = fibCtx.QueryParser(&dto)

	if l, err := strconv.Atoi(fibCtx.Query("limit")); err == nil && l >= 0 {
		dto.Limit = l
	} else {
		dto.Limit = 5 // default
	}

	if o, err := strconv.Atoi(fibCtx.Query("offset")); err == nil && o >= 0 {
		dto.Offset = o
	} else {
		dto.Offset = 0 // default
	}

	validTypes := map[string]bool{
		"Walking": true, "Yoga": true, "Stretching": true, "Cycling": true,
		"Swimming": true, "Dancing": true, "Hiking": true, "Running": true,
		"HIIT": true, "JumpRope": true,
	}

	if !validTypes[dto.ActivityType] {
		dto.ActivityType = "" // invalid → ignored
	}

	if raw := fibCtx.Query("doneAtForm"); raw != "" {
		if _, err := time.Parse(time.RFC3339, raw); err == nil {
			dto.DoneAtForm = raw
		} else {
			dto.DoneAtForm = "" // invalid → ignore
		}
	}

	if raw := fibCtx.Query("doneAtTo"); raw != "" {
		if _, err := time.Parse(time.RFC3339, raw); err == nil {
			dto.DoneAtTo = raw
		} else {
			dto.DoneAtTo = "" // invalid → ignore
		}
	}

	if raw := fibCtx.Query("caloriesBurnedMin"); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil && v >= 0 {
			dto.CaloriesBurnedMin = v
		} else {
			dto.CaloriesBurnedMin = 0 // invalid → ignore
		}
	}

	if raw := fibCtx.Query("caloriesBurnedMax"); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil && v >= 0 {
			dto.CaloriesBurnedMax = v
		} else {
			dto.CaloriesBurnedMax = 0 // invalid → ignore
		}
	}

	userId, ok := fibCtx.Locals("id").(int)
	if !ok {
		return fibCtx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user id cant be empty",
		})
	} else {
		dto.UserId = userId
	}

	activities, err := r.activityUseCase.GetAllActivities(fibCtx.Context(), dto)
	if err != nil {
		return fibCtx.Status(exceptions.MapToHttpStatusCode(err)).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return fibCtx.Status(fiber.StatusOK).JSON(activities)
}

func (r *activityHandler) GetById(fibCtx *fiber.Ctx) error {
	activityId := fibCtx.Params("id")
	id, err := strconv.Atoi(activityId)
	if err != nil {
		return fibCtx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "activityId must be an integer",
		})
	}

	activity, err := r.activityUseCase.GetActivityById(fibCtx.Context(), id)
	if err != nil {
		return fibCtx.Status(exceptions.MapToHttpStatusCode(err)).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return fibCtx.Status(fiber.StatusOK).JSON(activity)
}

func (r *activityHandler) Post(fibCtx *fiber.Ctx) error {
	if fibCtx.Get("Content-Type") != "application/json" {
		return fibCtx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid content type",
		})
	}

	var activityRequest dto.ActivityRequest
	if err := fibCtx.BodyParser(&activityRequest); err != nil {
		return fibCtx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var validate = validator.New()
	validate.RegisterValidation("iso8601", internal_validator.ValidateISODate)
	if err := validate.Struct(activityRequest); err != nil {
		return fibCtx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userId, ok := fibCtx.Locals("id").(int)
	if !ok {
		return fibCtx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user id cant be empty",
		})
	}

	activityResponse, err := r.activityUseCase.PostActivity(fibCtx.Context(), &activityRequest, userId)

	if err != nil {
		return fibCtx.Status(exceptions.MapToHttpStatusCode(err)).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return fibCtx.Status(fiber.StatusCreated).JSON(activityResponse)
}

func (r *activityHandler) Update(fibCtx *fiber.Ctx) error {

	activityId := fibCtx.Params("id")
	if activityId == "" {
		fmt.Println("activityId is required")
		return fibCtx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "activityId is required",
		})
	}

	userId, ok := fibCtx.Locals("id").(int)
	if !ok {
		fmt.Println("user id can not be empty")
		return fibCtx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user id cant be empty",
		})
	}

	if fibCtx.Method() == fiber.MethodPatch {
		if fibCtx.Get("Content-Type") != "application/json" {
			fmt.Println("invalid content type")
			return fibCtx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid content type",
			})
		}

		var activityRequest dto.ActivityUpdateRequest
		if err := fibCtx.BodyParser(&activityRequest); err != nil {
			fmt.Println(err)
			return fibCtx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		var validate = validator.New()
		validate.RegisterValidation("iso8601", internal_validator.ValidateISODate)
		if err := validate.Struct(activityRequest); err != nil {
			fmt.Println(err)
			return fibCtx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		activityResponse, err := r.activityUseCase.UpdateActivity(fibCtx.Context(), &activityRequest, userId, activityId)

		if err != nil {
			fmt.Println(err)
			if err == exceptions.ErrNotFound {
				return fibCtx.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "activityId is not found",
				})
			}
			return fibCtx.Status(exceptions.MapToHttpStatusCode(err)).JSON(fiber.Map{

				"error": err.Error(),
			})
		}

		return fibCtx.Status(fiber.StatusCreated).JSON(activityResponse)
	}

	err := r.activityUseCase.DisableActivity(fibCtx.Context(), userId, activityId)

	if err != nil {
		fmt.Println(err)
		if err == exceptions.ErrNotFound {
			return fibCtx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "activityId is not found",
			})
		}
		return fibCtx.Status(exceptions.MapToHttpStatusCode(err)).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return fibCtx.Status(fiber.StatusOK).JSON("Deleted")
}
