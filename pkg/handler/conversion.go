package handler

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/wildan3105/converto/pkg/api/schema"
	"github.com/wildan3105/converto/pkg/domain"
	"github.com/wildan3105/converto/pkg/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConversionHandler struct {
	conversionService *service.ConversionService
}

func NewConversionHandler(service *service.ConversionService) *ConversionHandler {
	return &ConversionHandler{
		conversionService: service,
	}
}

func (h *ConversionHandler) CreateConversion(c *fiber.Ctx) error {
	req := new(schema.CreateConversionRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "File is required"})
	}

	fileName := file.Filename
	if filepath.Ext(fileName) != ".shapr" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid file type. Only .shapr files are allowed",
		})
	}

	req.FileName = fileName

	targetFormat := c.FormValue("target_format")
	if targetFormat == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Target format is required",
		})
	}

	allowedFormats := map[string]bool{
		".step": true,
		".iges": true,
		".stl":  true,
		".obj":  true,
	}

	if !allowedFormats[targetFormat] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid target format. Allowed formats are: .step, .iges, .stl, .obj",
		})
	}

	fileSize := file.Size
	req.FileSize = fileSize

	conversion, err := h.conversionService.CreateConversion(context.Background(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create conversion",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(conversion)
}

func (h *ConversionHandler) GetConversions(c *fiber.Ctx) error {
	status := c.Query("status")

	// Validate the status if provided
	if status != "" && !IsValidConversionStatus(status) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid status. Must be one of 'pending', 'in_progress', 'completed', 'failed'",
		})
	}

	page := c.QueryInt("page", 1)
	if page < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid page value. Must at least be 1",
		})
	}

	limit := c.QueryInt("limit", 10)
	if limit < 1 || limit > 20 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid limit value. Must be below 21 and above 0",
		})
	}
	conversions, err := h.conversionService.ListConversions(context.Background(), status, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch conversions",
		})
	}
	return c.JSON(conversions)
}

func (h *ConversionHandler) GetConversionByID(c *fiber.Ctx) error {
	id := c.Params("id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format. Must be a valid MongoDB ObjectID",
		})
	}

	conversion, err := h.conversionService.GetConversionByID(context.Background(), objectID.Hex())

	if err != nil {
		if errors.Is(err, fiber.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Conversion not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch conversion",
		})
	}

	return c.JSON(conversion)
}

// IsValidConversionStatus checks if the provided status is a valid ConversionStatus
func IsValidConversionStatus(status string) bool {
	switch domain.ConversionStatus(status) {
	case domain.ConversionPending, domain.ConversionInProgress, domain.ConversionCompleted, domain.ConversionFailed:
		return true
	default:
		return false
	}
}
