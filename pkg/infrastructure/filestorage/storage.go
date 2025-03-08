package filestorage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

// FileStorage interface to abstract file storage operations
type FileStorage interface {
	SaveFile(file *multipart.FileHeader, destPath string) (string, error)
}

// LocalFileStorage is an implementation of FileStorage using local filesystem
type LocalFileStorage struct {
	baseDir string
}

// NewLocalFileStorage creates a new instance of LocalFileStorage
func NewLocalFileStorage(baseDir string) *LocalFileStorage {
	return &LocalFileStorage{baseDir: baseDir}
}

// SaveFile saves the uploaded file to the specified path or default directory
func (l *LocalFileStorage) SaveFile(file *multipart.FileHeader, destPath string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	if destPath == "" {
		destPath = filepath.Join(l.baseDir, fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename))
	} else {
		destPath = filepath.Join(l.baseDir, destPath)
	}

	dir := filepath.Dir(destPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	dest, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dest.Close()

	if _, err := io.Copy(dest, src); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return destPath, nil
}
