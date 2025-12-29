package app

import (
	"time"

	"github.com/arieffadhlan/go-fitbyte/internal/config"
	"github.com/arieffadhlan/go-fitbyte/internal/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jmoiron/sqlx"
)

func NewServer(cfg *config.Config, db *sqlx.DB) *fiber.App {
	app := fiber.New(fiber.Config{
		IdleTimeout:  600 * time.Second,
		ReadTimeout:  600 * time.Second,
		WriteTimeout: 600 * time.Second,
	})

	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path} - ${error}\n",
	}))

	handlers.SetupRouter(cfg, db, app)

	return app
}
