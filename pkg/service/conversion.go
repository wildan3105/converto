package service

import (
	"context"

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

// ListConversions fetches conversions from the repository
func (s *ConversionService) ListConversions(ctx context.Context, filter bson.M, limit, offset int64) ([]*domain.Conversion, error) {
	return s.repo.ListConversions(ctx, filter, limit, offset)
}
