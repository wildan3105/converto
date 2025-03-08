package worker

import (
	"context"
	"fmt"

	"github.com/wildan3105/converto/pkg/domain"
)

// Handle processes a conversion job
func (w *Worker) Handle(ctx context.Context, job domain.ConversionJob) error {
	conversion, err := w.repo.GetConversionByID(ctx, job.JobID)
	if err != nil {
		return fmt.Errorf("failed to fetch conversion: %w", err)
	}

	fmt.Println("processing conversion", conversion)

	// if err := w.UpdateProgress(ctx, conversion, 10); err != nil {
	// 	return err
	// }

	// if err := w.TriggerIfNeeded(ctx, conversion); err != nil {
	// 	return err
	// }

	return nil
}
