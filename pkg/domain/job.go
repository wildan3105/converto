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
	JobID        string    `bson:"job_id" json:"job_id"`
	Source       JobSource `bson:"source" json:"source"`
	CreatedAt    time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time `bson:"updated_at" json:"updated_at"`
	ErrorMessage *string   `bson:"error_message,omitempty" json:"error_message,omitempty"`
}
