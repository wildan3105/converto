package filestorage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

type FileCategory string

const (
	FileCategoryOriginal  FileCategory = "original"
	FileCategoryConverted FileCategory = "converted"
)

// FileStorage interface to abstract file storage operations
type FileStorage interface {
	SaveFile(file *multipart.FileHeader, fileCategory FileCategory, id string, destPath string) (string, error)
	CopyFile(srcPath, destPath string, progressCb func(progress int)) (string, error)
	GetFullPath(fileCategory FileCategory, id string, fileName string) string
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
func (l *LocalFileStorage) SaveFile(file *multipart.FileHeader, fileCategory FileCategory, id string, destPath string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	if destPath == "" {
		destPath = filepath.Join(l.baseDir, fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename))
	} else {
		destPath = filepath.Join(l.baseDir+string(fileCategory)+"/"+id, destPath)
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

// CopyFile copies a file from srcPath to destPath and updates progress
func (l *LocalFileStorage) CopyFile(srcPath, destPath string, progressCb func(progress int)) (string, error) {
	src, err := os.Open(srcPath)
	if err != nil {
		return "", fmt.Errorf("failed to open source file: %w", err)
	}
	defer src.Close()

	dir := filepath.Dir(destPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	dest, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dest.Close()

	srcInfo, err := src.Stat()
	if err != nil {
		return "", err
	}

	totalBytes := srcInfo.Size()
	buffer := make([]byte, 1024*1024) // 1 MB buffer
	var copiedBytes int64
	lastReportedProgress := 0

	for {
		n, err := src.Read(buffer)
		if n > 0 {
			if _, writeErr := dest.Write(buffer[:n]); writeErr != nil {
				return "", writeErr
			}
			copiedBytes += int64(n)
			progress := int((copiedBytes * 100) / totalBytes)

			if progress >= lastReportedProgress+10 {
				progressCb(progress)
				lastReportedProgress = progress
			}
		}

		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
	}

	progressCb(100)
	return destPath, nil
}

// GetFullPath constructs the full path for a file given its category and name
func (l *LocalFileStorage) GetFullPath(fileCategory FileCategory, id string, fileName string) string {
	return filepath.Join(l.baseDir, string(fileCategory), id, fileName)
}
