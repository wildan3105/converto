package mongodb

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FileMetadata represents metadata for both original and converted files
type FileMetadata struct {
	OriginalName  string `bson:"original_name" json:"original_name"`
	OriginalPath  string `bson:"original_path" json:"original_path"`
	ConvertedName string `bson:"converted_name" json:"converted_name"`
	ConvertedPath string `bson:"converted_path" json:"converted_path"`
	SizeInBytes   int64  `bson:"size_in_bytes" json:"size_in_bytes"`
	MimeType      string `bson:"mime_type" json:"mime_type"`
}

// Conversion holds specific conversion job details.
type ConversionData struct {
	TargetFormat string     `bson:"target_format" json:"target_format"`
	Progress     int        `bson:"progress" json:"progress"`
	Status       string     `bson:"status" json:"status"`
	ErrorMessage *string    `bson:"error_message" json:"error_message,omitempty"`
	StartedAt    *time.Time `bson:"started_at" json:"started_at,omitempty"`
	CompletedAt  *time.Time `bson:"completed_at" json:"completed_at,omitempty"`
	CallbackURL  *string    `bson:"callback_url" json:"callback_url,omitempty"`
}

// JobMetadata holds information about the job source and message queue
type JobMetadata struct {
	Source    string    `bson:"source" json:"source"`
	QueueID   string    `bson:"queue_id" json:"queue_id,omitempty"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// Conversion represents the complete document stored in the database
type Conversion struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ConversionID string             `bson:"conversion_id" json:"conversion_id"`
	File         FileMetadata       `bson:"file" json:"file"`
	Conversion   ConversionData     `bson:"conversion" json:"conversion"`
	Job          JobMetadata        `bson:"job" json:"job"`
}
