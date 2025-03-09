package worker

import (
	"context"
	"net/http"

	"github.com/wildan3105/converto/pkg/domain"
)

// TriggerIfNeeded triggers a callback if a callback URL is provided
func (w *Worker) TriggerIfNeeded(ctx context.Context, conversion *domain.Conversion) error {
	if *conversion.Conversion.CallbackURL == "" {
		return nil
	}

	req, err := http.NewRequestWithContext(ctx, "POST", *conversion.Conversion.CallbackURL, nil)
	if err != nil {
		log.Warn("Failed to create callback request: %v", err)
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Warn("Failed to trigger callback: %v", err)
		return err
	}

	log.Info("Callback triggered, status code: %d", resp.StatusCode)
	return nil
}
