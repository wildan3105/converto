package worker

import (
	"context"
	"time"

	"github.com/wildan3105/converto/pkg/domain"
	"go.mongodb.org/mongo-driver/bson"
)

// UpdateProgress simulates progress and updates the database in 10% increments
func (w *Worker) UpdateProgress(ctx context.Context, conversion *domain.Conversion, progress int) error {
	for i := progress; i <= 100; i += 10 {
		time.Sleep(1 * time.Second)
		conversion.Conversion.Progress = i

		filter := bson.M{
			"conversion.progress": conversion.Conversion.Progress,
		}

		if err := w.repo.UpdateConversion(ctx, conversion.ID, filter); err != nil {
			log.Warn("Failed to update progress to %d%%: %v", i, err)
			return err
		}
		log.Info("Progress updated to %d%%", i)
	}
	return nil
}
