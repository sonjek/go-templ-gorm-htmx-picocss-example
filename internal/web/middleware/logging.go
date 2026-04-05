package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v3"
)

func LoggingMiddleware(c fiber.Ctx) error {
	start := time.Now()

	err := c.Next()

	statusCode := c.Response().StatusCode()
	if statusCode == 0 {
		statusCode = http.StatusOK
	}

	// #nosec G706 - input is sanitized via slog's internal escaping
	slog.Info("Request",
		"status", statusCode,
		"method", c.Method(),
		"path", c.Path(),
		"duration", time.Since(start),
	)

	return err
}
