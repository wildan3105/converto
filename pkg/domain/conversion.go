package domain

import "time"

type ConversionStatus string

const (
	ConversionPending    ConversionStatus = "pending"
	ConversionInProgress ConversionStatus = "in-progress"
	ConversionCompleted  ConversionStatus = "completed"
	ConversionFailed     ConversionStatus = "failed"
)

// Conversion represents a conversion task with associated metadata and job status
type Conversion struct {
	ID         string         `bson:"_id,omitempty" json:"id"`
	File       FileMetadata   `bson:"file" json:"file"`
	Conversion ConversionData `bson:"conversion" json:"conversion"`
	Job        ConversionJob  `bson:"job" json:"job"`
}

// ConversionData represents the metadata and status of a conversion task
type ConversionData struct {
	TargetFormat string           `bson:"target_format" json:"target_format"`
	Progress     int              `bson:"progress" json:"progress"`
	Status       ConversionStatus `bson:"status" json:"status"`
	ErrorMessage *string          `bson:"error_message" json:"error_message,omitempty"`
	StartedAt    *time.Time       `bson:"started_at" json:"started_at,omitempty"`
	CompletedAt  *time.Time       `bson:"completed_at" json:"completed_at,omitempty"`
	CallbackURL  *string          `bson:"callback_url" json:"callback_url,omitempty"`
}
