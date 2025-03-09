package domain

import "time"

// JobSource represents the source of the job
type JobSource string

const (
	JobSourceAPI JobSource = "API"
	JobSourceCLI JobSource = "CLI"
)

// ConversionJob represents the job status and metadata for a conversion
type ConversionJob struct {
	JobID        string `json:"job_id"`
	ConversionID string
	Source       JobSource `json:"source"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	ErrorMessage *string   `json:"error_message,omitempty"`
}
