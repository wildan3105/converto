package schema

import (
	"mime/multipart"

	"github.com/wildan3105/converto/pkg/domain"
)

// CreateConversionRequest defines the payload for creating a conversion
type CreateConversionRequest struct {
	File         *multipart.FileHeader `form:"file" binding:"required"`                                      // The .shapr file to convert
	TargetFormat string                `form:"target_format" binding:"required,oneof=.step .iges .stl .obj"` // The target format to convert to
	CallbackURL  string                `form:"callback_url"`
	FileSize     int64
	FileName     string
}

type CreateConversionResponse struct {
	ID      string                  `json:"id"`
	Status  domain.ConversionStatus `json:"status"`
	Message string                  `json:"message"`
}

type ConversionResponse struct {
	ID               string                  `json:"id"`
	Status           domain.ConversionStatus `json:"status"`
	Progress         int                     `json:"progress"`
	OriginalFileURL  string                  `json:"original_file_url"`
	ConvertedFileURL string                  `json:"converted_file_url,omitempty"`
}
