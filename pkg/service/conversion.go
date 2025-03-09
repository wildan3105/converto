package service

import (
	"context"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/wildan3105/converto/pkg/api/schema"
	"github.com/wildan3105/converto/pkg/domain"
	"github.com/wildan3105/converto/pkg/infrastructure/filestorage"
	"github.com/wildan3105/converto/pkg/infrastructure/rabbitmq"
	"github.com/wildan3105/converto/pkg/logger"
	"github.com/wildan3105/converto/pkg/repository"
)

var log = logger.GetInstance()

type ConversionService struct {
	repo      repository.ConversionRepository
	publisher *rabbitmq.Publisher
	storage   filestorage.FileStorage
}

func NewConversionService(repo repository.ConversionRepository, publisher *rabbitmq.Publisher, storage filestorage.FileStorage) *ConversionService {
	return &ConversionService{
		repo:      repo,
		publisher: publisher,
		storage:   storage,
	}
}

// CreateConversion creates a conversion
func (s *ConversionService) CreateConversion(ctx context.Context, req *schema.CreateConversionRequest) (schema.CreateConversionResponse, error) {
	convertedFileName := strings.TrimSuffix(req.FileName, ".shapr") + req.TargetFormat

	if req.File == nil {
		return schema.CreateConversionResponse{}, fiber.NewError(fiber.StatusBadRequest, "File is required")
	}

	fileID := uuid.NewString()

	originalFilePath, err := s.storage.SaveFile(req.File, domain.FileCategoryOriginal, fileID, req.FileName)
	if err != nil {
		return schema.CreateConversionResponse{}, fiber.NewError(fiber.StatusInternalServerError, "Failed to save original file")
	}

	conversionPayload := &domain.Conversion{
		File: domain.FileMetadata{
			OriginalName:  req.FileName,
			OriginalPath:  originalFilePath,
			ConvertedName: convertedFileName,
			SizeInBytes:   req.FileSize,
			ID:            fileID,
		},
		Conversion: domain.ConversionData{
			TargetFormat: req.TargetFormat,
			Progress:     0,
			Status:       domain.ConversionPending,
			StartedAt:    time.Now(),
		},
		Job: domain.ConversionJob{
			ID:        uuid.NewString(),
			Source:    domain.JobSourceAPI,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	id, err := s.repo.CreateConversion(ctx, conversionPayload)

	if err != nil {
		return schema.CreateConversionResponse{}, err
	}

	event := schema.ConversionEvent{
		JobID:        conversionPayload.Job.ID,
		ConversionID: id,
		Source:       domain.JobSourceAPI,
		CreatedAt:    conversionPayload.Job.CreatedAt,
		UpdatedAt:    conversionPayload.Job.UpdatedAt,
	}

	publishErr := s.publisher.PublishConversionJob(ctx, event, "conversion", "created")

	if publishErr != nil {
		log.Warn("Error when publishing %v", err)
		return schema.CreateConversionResponse{}, err
	}

	return schema.CreateConversionResponse{
		ID:      id,
		Status:  conversionPayload.Conversion.Status,
		Message: "Conversion created successfully",
	}, nil
}

// ListConversions fetches conversions from the repository and maps them to the response schema
func (s *ConversionService) ListConversions(ctx context.Context, status string, page, limit int) (schema.ListConversionsResponse, error) {
	offset := (page - 1) * limit

	conversions, err := s.repo.ListConversions(ctx, status, limit, offset)
	if err != nil {
		return schema.ListConversionsResponse{}, err
	}

	responses := make([]schema.ConversionResponse, len(conversions))
	for i, conversion := range conversions {
		responses[i] = schema.ConversionResponse{
			ID:                conversion.ID,
			Status:            conversion.Conversion.Status,
			Progress:          conversion.Conversion.Progress,
			OriginalFilePath:  conversion.File.OriginalPath,
			ConvertedFilePath: conversion.File.ConvertedPath,
		}
	}

	responseData := schema.ListConversionsResponse{
		Page:  page,
		Limit: limit,
		Data:  responses,
	}

	return responseData, nil
}

// GetConversionByID fetches conversion by ID from the repository
func (s *ConversionService) GetConversionByID(ctx context.Context, id string) (schema.ConversionResponse, error) {
	conversion, err := s.repo.GetConversionByID(ctx, id)
	if err != nil {
		return schema.ConversionResponse{}, err
	}
	if conversion == nil {
		return schema.ConversionResponse{}, fiber.ErrNotFound
	}

	return schema.ConversionResponse{
		ID:                conversion.ID,
		Status:            conversion.Conversion.Status,
		Progress:          conversion.Conversion.Progress,
		OriginalFilePath:  conversion.File.OriginalPath,
		ConvertedFilePath: conversion.File.ConvertedPath,
	}, nil
}

// GetFileByConversionIdAndType returns the file path and name based on conversion ID and file type
func (s *ConversionService) GetFileByConversionIdAndType(ctx context.Context, id string, fileType string) (schema.GetFileByConversionId, error) {
	conversion, err := s.repo.GetConversionByID(ctx, id)
	if err != nil {
		return schema.GetFileByConversionId{
			Path:     "",
			FileName: "",
		}, err
	}
	if conversion == nil {
		return schema.GetFileByConversionId{
			Path:     "",
			FileName: "",
		}, fiber.ErrNotFound
	}

	switch fileType {
	case "original":
		return schema.GetFileByConversionId{
			Path:     conversion.File.OriginalPath,
			FileName: conversion.File.OriginalName,
		}, nil
	case "converted":
		return schema.GetFileByConversionId{
			Path:     conversion.File.ConvertedPath,
			FileName: conversion.File.ConvertedName,
		}, nil
	default:
		return schema.GetFileByConversionId{
			Path:     "",
			FileName: "",
		}, nil
	}
}
