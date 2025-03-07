package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/wildan3105/converto/pkg/logger"
	"github.com/wildan3105/converto/pkg/service"
)

type HealthHandler struct {
	healthService *service.HealthService
}

func NewHealthHandler(service *service.HealthService) *HealthHandler {
	return &HealthHandler{
		healthService: service,
	}
}

var log = logger.GetInstance()

// GetConversions handles GET /api/v1/conversions
func (h *HealthHandler) Check(c *fiber.Ctx) error {
	err := h.healthService.CheckDependencies()
	if err != nil {
		log.Warn("Dependencies check failed: %v", err)
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"message":   "One of dependencies connection is failed",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":   "ok",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
