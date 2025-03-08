package service

import (
	"context"
	"strings"
	"time"

	"github.com/wildan3105/converto/pkg/api/schema"
	"github.com/wildan3105/converto/pkg/domain"
	"github.com/wildan3105/converto/pkg/repository"
	"go.mongodb.org/mongo-driver/bson"
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

// ListConversions fetches conversions from the repository
func (s *ConversionService) ListConversions(ctx context.Context, filter bson.M, limit, offset int64) ([]*domain.Conversion, error) {
	return s.repo.ListConversions(ctx, filter, limit, offset)
}
