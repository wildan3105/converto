package worker

import (
	"context"
	"fmt"

	"github.com/wildan3105/converto/pkg/domain"
	"github.com/wildan3105/converto/pkg/infrastructure/filestorage"
	"go.mongodb.org/mongo-driver/bson"
)

// Handle processes a conversion job
func (w *Worker) Handle(ctx context.Context, job domain.ConversionJob) error {
	conversion, err := w.repo.GetConversionByID(ctx, job.ConversionID)
	if err != nil {
		return fmt.Errorf("failed to fetch conversion: %w", err)
	}

	log.Info("Processing conversion: %s", conversion.ID)

	originalPath := conversion.File.OriginalPath
	convertedName := conversion.File.ConvertedName

	// Use GetFullPath method to generate the full path for the converted file
	convertedPath := w.storage.GetFullPath(filestorage.FileCategoryConverted, convertedName)

	// Simulate conversion process with progress updates
	if err := w.UpdateProgress(ctx, conversion, 10); err != nil {
		return err
	}

	// Emulate file conversion by copying the original file to the converted path
	convertedFilePath, err := w.storage.CopyFile(originalPath, convertedPath)
	if err != nil {
		return fmt.Errorf("failed to emulate file conversion: %w", err)
	}
	log.Info("Converted file stored at: %s", convertedFilePath)

	// Combine update of converted path, status, and final progress to 100%
	updateData := bson.M{
		"file.converted_path": convertedFilePath,
		"conversion.status":   "completed",
		"conversion.progress": 100,
	}

	fmt.Println("updating the status and progress to 100 percent file of conversionID > ", conversion.ID)

	if err := w.repo.UpdateConversion(ctx, conversion.ID, updateData); err != nil {
		return fmt.Errorf("failed to update conversion status: %w", err)
	}

	return nil
}
