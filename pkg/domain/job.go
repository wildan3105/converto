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
	ID           string    `bson:"id" json:"id"`
	Source       JobSource `bson:"source" json:"source"`
	CreatedAt    time.Time `bson:"createdAt" json:"created_at"`
	UpdatedAt    time.Time `bson:"updatedAt" json:"updated_at"`
	ErrorMessage *string   `bson:"errorMessage" json:"error_message,omitempty"`
}
