package domain

// FileMetadata represents metadata information for files in the conversion process
type FileMetadata struct {
	OriginalName  string `bson:"original_name" json:"original_name"`
	OriginalPath  string `bson:"original_path" json:"original_path"`
	ConvertedName string `bson:"converted_name" json:"converted_name"`
	ConvertedPath string `bson:"converted_path" json:"converted_path"`
	SizeInBytes   int64  `bson:"size_in_bytes" json:"size_in_bytes"`
	MimeType      string `bson:"mime_type" json:"mime_type"`
}
