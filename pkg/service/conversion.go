package service

import (
	"context"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/wildan3105/converto/pkg/api/schema"
	"github.com/wildan3105/converto/pkg/domain"
	"github.com/wildan3105/converto/pkg/infrastructure/rabbitmq"
	"github.com/wildan3105/converto/pkg/repository"
)

type ConversionService struct {
	repo      repository.ConversionRepository
	publisher *rabbitmq.Publisher
}

func NewConversionService(repo repository.ConversionRepository, publisher *rabbitmq.Publisher) *ConversionService {
	return &ConversionService{
		repo:      repo,
		publisher: publisher,
	}
}

// CreateConversion creates a conversion
func (s *ConversionService) CreateConversion(ctx context.Context, req *schema.CreateConversionRequest) (schema.CreateConversionResponse, error) {
	convertedFileName := strings.TrimSuffix(req.FileName, ".shapr") + req.TargetFormat

	conversionPayload := &domain.Conversion{
		File: domain.FileMetadata{
			OriginalName:  req.FileName,
			ConvertedName: convertedFileName,
			SizeInBytes:   req.FileSize,
		},
		Conversion: domain.ConversionData{
			TargetFormat: req.TargetFormat,
			Progress:     0,
			Status:       domain.ConversionPending,
		},
		Job: domain.ConversionJob{
			JobID:     uuid.NewString(),
			Source:    domain.JobSourceAPI,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	id, err := s.repo.CreateConversion(ctx, conversionPayload)

	if err != nil {
		return schema.CreateConversionResponse{}, err
	}

	publishErr := s.publisher.PublishConversionJob(ctx, conversionPayload.Job, "conversion", "created")

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
			ID:               conversion.ID,
			Status:           conversion.Conversion.Status,
			Progress:         conversion.Conversion.Progress,
			OriginalFileURL:  conversion.File.OriginalPath,
			ConvertedFileURL: conversion.File.ConvertedPath,
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
		ID:               conversion.ID,
		Status:           conversion.Conversion.Status,
		Progress:         conversion.Conversion.Progress,
		OriginalFileURL:  conversion.File.OriginalPath,
		ConvertedFileURL: conversion.File.ConvertedPath,
	}, nil
}
