package schema

import (
	"mime/multipart"

	"github.com/wildan3105/converto/pkg/domain"
)

// CreateConversionRequest defines the payload for creating a conversion
type CreateConversionRequest struct {
	File         *multipart.FileHeader `form:"file" binding:"required"`                                      // The .shapr file to convert
	TargetFormat string                `form:"target_format" binding:"required,oneof=.step .iges .stl .obj"` // The target format to convert to
	FileSize     int64
	FileName     string
}

type CreateConversionResponse struct {
	ID      string                  `json:"id"`
	Status  domain.ConversionStatus `json:"status"`
	Message string                  `json:"message"`
}

type ListConversionsResponse struct {
	Page  int                  `json:"page"`
	Limit int                  `json:"limit"`
	Data  []ConversionResponse `json:"data"`
}

type ConversionResponse struct {
	ID                string                  `json:"id"`
	Status            domain.ConversionStatus `json:"status"`
	Progress          int                     `json:"progress"`
	OriginalFilePath  string                  `json:"original_file_path"`
	ConvertedFilePath string                  `json:"converted_file_path,omitempty"`
}

type GetFileByConversionId struct {
	Path     string
	FileName string
}
