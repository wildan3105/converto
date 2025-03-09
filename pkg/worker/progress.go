package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/wildan3105/converto/pkg/domain"
	"go.mongodb.org/mongo-driver/bson"
)

// UpdateProgress simulates progress and updates the database in 10% increments
func (w *Worker) UpdateProgress(ctx context.Context, conversion *domain.Conversion, startProgress int) error {
	for i := startProgress; i < 100; i += 10 {
		time.Sleep(1 * time.Second)
		conversion.Conversion.Progress = i

		updateData := bson.M{
			"conversion.progress": i,
			"conversion.status":   domain.ConversionInProgress,
		}

		fmt.Println("updating the progress file of conversionID > ", conversion.ID)

		if err := w.repo.UpdateConversion(ctx, conversion.ID, updateData); err != nil {
			log.Warn("Failed to update progress to %d%%: %v", i, err)
			return err
		}
		log.Info("Progress updated to %d%%", i)
	}
	return nil
}
