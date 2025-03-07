package handler

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/wildan3105/converto/pkg/service"
)

type ConversionHandler struct {
	conversionService *service.ConversionService
}

func NewConversionHandler(service *service.ConversionService) *ConversionHandler {
	return &ConversionHandler{
		conversionService: service,
	}
}

func (h *ConversionHandler) GetConversions(c *fiber.Ctx) error {
	conversions, err := h.conversionService.ListConversions(context.Background(), nil, 5, 5)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch conversions",
		})
	}
	return c.JSON(conversions)
}
