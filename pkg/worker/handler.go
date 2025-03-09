package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/wildan3105/converto/pkg/api/schema"
	"github.com/wildan3105/converto/pkg/domain"
	"go.mongodb.org/mongo-driver/bson"
)

// Handle processes a conversion job
func (w *Worker) Handle(ctx context.Context, event schema.ConversionEvent) error {
	conversion, err := w.repo.GetConversionByID(ctx, event.ConversionID)
	if err != nil {
		return fmt.Errorf("failed to fetch conversion: %w", err)
	}

	log.Info("Processing conversion: %s", conversion.ID)

	originalPath := conversion.File.OriginalPath
	convertedName := conversion.File.ConvertedName
	fileID := conversion.File.ID

	convertedPath := w.storage.GetFullPath(domain.FileCategoryConverted, fileID, convertedName)

	progressCb := func(progress int) {
		status := domain.ConversionInProgress
		if progress == 100 {
			status = domain.ConversionCompleted
		}

		updateData := bson.M{
			"conversion.progress": progress,
			"conversion.status":   status,
		}

		if progress == 100 {
			updateData["file.convertedPath"] = convertedPath
			updateData["conversion.completedAt"] = time.Now()
		}

		log.Info("Conversion progress: %d%% for conversion ID: %s", progress, conversion.ID)

		if err := w.repo.UpdateConversion(ctx, conversion.ID, updateData); err != nil {
			log.Warn("Failed to update progress to %d%%: %v", progress, err)
		}
	}

	convertedFilePath, err := w.storage.CopyFile(originalPath, convertedPath, progressCb)
	if err != nil {
		return fmt.Errorf("failed to emulate file conversion: %w", err)
	}

	log.Info("Converted file stored at: %s", convertedFilePath)

	return nil
}
