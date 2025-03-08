package service

import (
	"context"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/wildan3105/converto/pkg/api/schema"
	"github.com/wildan3105/converto/pkg/domain"
	"github.com/wildan3105/converto/pkg/repository"
)

type ConversionService struct {
	repo repository.ConversionRepository
}

func NewConversionService(repo repository.ConversionRepository) *ConversionService {
	return &ConversionService{repo: repo}
}

// CreateConversion creates a conversion
func (s *ConversionService) CreateConversion(ctx context.Context, req *schema.CreateConversionRequest) (schema.CreateConversionResponse, error) {
	convertedFileName := strings.TrimSuffix(req.FileName, ".shapr") + req.TargetFormat

	conversionPayload := &domain.Conversion{
		File: domain.FileMetadata{
			OriginalName:  req.FileName,
			OriginalPath:  "path-original",   // will be fetched from the stored file
			ConvertedName: convertedFileName, // will be fetched later
			ConvertedPath: "path-converted",  // will be fetched later
			SizeInBytes:   req.FileSize,
		},
		Conversion: domain.ConversionData{
			TargetFormat: req.TargetFormat,
			Progress:     0,
			Status:       domain.ConversionPending,
		},
		Job: domain.ConversionJob{
			JobID:     "1",
			Source:    domain.JobSourceAPI,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	res, err := s.repo.CreateConversion(ctx, conversionPayload)

	if err != nil {
		return schema.CreateConversionResponse{}, err
	}

	return schema.CreateConversionResponse{
		ID:      res,
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
