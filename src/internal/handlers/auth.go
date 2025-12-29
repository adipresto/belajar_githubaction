package handlers

import (
	"net/http"

	"github.com/arieffadhlan/go-fitbyte/internal/dto"
	"github.com/arieffadhlan/go-fitbyte/internal/pkg/exceptions"
	"github.com/arieffadhlan/go-fitbyte/internal/usecases/auth"
	"github.com/gofiber/fiber/v2"
)

type AuthHandlerInterface interface {
	Login(ctx *fiber.Ctx) error
	Register(ctx *fiber.Ctx) error
}

type AuthHandler struct {
	uc auth.AuthUseCaseInterface
}

func NewAuthHandler(uc auth.AuthUseCaseInterface) AuthHandlerInterface {
	return &AuthHandler{uc: uc}
}

func (ah *AuthHandler) Register(ctx *fiber.Ctx) error {
	var req dto.AuthRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "invalid request payload"})
	}

	r, err := ah.uc.Register(ctx.Context(), &req)
	if err != nil {
		if appErr, ok := err.(*exceptions.AppError); ok {
			return ctx.Status(appErr.Code).JSON(appErr)
		} else {
			return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
	}

	return ctx.Status(http.StatusCreated).JSON(r)
}

func (ah *AuthHandler) Login(ctx *fiber.Ctx) error {
	var req dto.AuthRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "invalid request payload"})
	}

	r, err := ah.uc.Login(ctx.Context(), &req)
	if err != nil {
		if appErr, ok := err.(*exceptions.AppError); ok {
			return ctx.Status(appErr.Code).JSON(appErr)
		} else {
			return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": err.Error()})
		}
	}

	return ctx.JSON(r)
}
