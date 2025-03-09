package domain

type FileCategory string

const (
	FileCategoryOriginal  FileCategory = "original"
	FileCategoryConverted FileCategory = "converted"
)

// FileMetadata represents metadata information for files in the conversion process
type FileMetadata struct {
	OriginalName  string `bson:"originalName" json:"original_name"`
	OriginalPath  string `bson:"originalPath" json:"original_path"`
	ConvertedName string `bson:"convertedName" json:"converted_name"`
	ConvertedPath string `bson:"convertedPath" json:"converted_path"`
	SizeInBytes   int64  `bson:"sizeInBytes" json:"size_in_bytes"`
	ID            string `json:"id,omitempty"`
}
