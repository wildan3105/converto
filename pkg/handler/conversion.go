package handler

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/wildan3105/converto/pkg/api/schema"
	"github.com/wildan3105/converto/pkg/domain"
	"github.com/wildan3105/converto/pkg/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConversionHandler struct {
	conversionService service.ConversionService
}

func NewConversionHandler(service service.ConversionService) *ConversionHandler {
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

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse multipart form",
		})
	}

	files := form.File["file"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "File is required"})
	}
	if len(files) > 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Only one file is allowed"})
	}

	file := files[0]
	req.File = file

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

	req.FileSize = file.Size

	conversion, err := h.conversionService.CreateConversion(context.Background(), req)
	if err != nil {
		if err.Error() == "service temporarily unavailable" {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error": "Service temporarily unavailable. Please try again later.",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create conversion",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(conversion)
}

func (h *ConversionHandler) GetConversions(c *fiber.Ctx) error {
	status := c.Query("status")

	if status != "" && !isValidConversionStatus(status) {
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

func (h *ConversionHandler) GetFileByConversionId(c *fiber.Ctx) error {
	id := c.Params("id")
	fileType := c.Query("type")

	if fileType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File type query is required. Must be one of 'type=original', 'type=converted'",
		})
	}

	if !isValidFileType(fileType) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid file type. Must be one of 'original', 'converted'",
		})
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format. Must be a valid MongoDB ObjectID",
		})
	}

	fileDetails, err := h.conversionService.GetFileByConversionIdAndType(context.Background(), objectID.Hex(), fileType)
	if err != nil {
		if errors.Is(err, fiber.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "File not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch file",
		})
	}

	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileDetails.FileName))
	c.Set("Content-Type", "application/octet-stream")

	return c.SendFile(fileDetails.Path, false)
}

// isValidConversionStatus checks if the provided status is a valid ConversionStatus
func isValidConversionStatus(status string) bool {
	switch domain.ConversionStatus(status) {
	case domain.ConversionPending, domain.ConversionInProgress, domain.ConversionCompleted, domain.ConversionFailed:
		return true
	default:
		return false
	}
}

// isValidFileType checks if the provided file type is a valid FileCategory
func isValidFileType(fileType string) bool {
	switch domain.FileCategory(fileType) {
	case domain.FileCategoryConverted, domain.FileCategoryOriginal:
		return true
	default:
		return false
	}
}
