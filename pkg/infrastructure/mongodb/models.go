package mongodb

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// JobSource represents the source of the job
type JobSource string

const (
	JobSourceAPI JobSource = "API"
	JobSourceCLI JobSource = "CLI"
)

// FileMetadata represents metadata for both original and converted files
type FileMetadata struct {
	OriginalName  string `bson:"originalName"`
	OriginalPath  string `bson:"originalPath"`
	ConvertedName string `bson:"convertedName"`
	ConvertedPath string `bson:"convertedPath"`
	SizeInBytes   int64  `bson:"sizeInBytes"`
	ID            string `bson:"id,omitempty"`
}

// Conversion holds specific conversion job details.
type ConversionData struct {
	TargetFormat string     `bson:"targetFormat"`
	Progress     int        `bson:"progress"`
	Status       string     `bson:"status"`
	ErrorMessage *string    `bson:"errorMessage,omitempty"`
	StartedAt    *time.Time `bson:"startedAt,omitempty"`
	CompletedAt  *time.Time `bson:"completedAt,omitempty"`
}

// JobMetadata holds information about the job source and message queue
type JobMetadata struct {
	Source    JobSource `bson:"source"`
	ID        string    `bson:"id,omitempty"`
	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`
}

// Conversion represents the complete document stored in the database
type Conversion struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	File       FileMetadata       `bson:"file"`
	Conversion ConversionData     `bson:"conversion"`
	Job        JobMetadata        `bson:"job"`
}
