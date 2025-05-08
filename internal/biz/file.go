package biz

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
)

type FileRepo interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader, contentType string) (string, error)
}

type FileUsecase struct {
	repo FileRepo
	log  *log.Helper
}

func NewFileUsecase(repo FileRepo, logger log.Logger) *FileUsecase {
	return &FileUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

// UploadFile validates and uploads a file
func (uc *FileUsecase) UploadFile(ctx context.Context, file *multipart.FileHeader) (string, error) {
	// Validate file size (max 10MB)
	const maxSize = 10 * 1024 * 1024 // 10MB
	if file.Size > maxSize {
		return "", fmt.Errorf("file size exceeds %d bytes", maxSize)
	}

	// Validate file type (image or video)
	contentType := file.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") && !strings.HasPrefix(contentType, "video/") {
		return "", fmt.Errorf("invalid file type: %s, only images and videos allowed", contentType)
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".mp4":  true,
		".avi":  true,
	}
	if !allowedExts[ext] {
		return "", fmt.Errorf("invalid file extension: %s", ext)
	}

	// Upload file
	fileURL, err := uc.repo.UploadFile(ctx, file, contentType)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	uc.log.Infof("Successfully processed file %s", file.Filename)
	return fileURL, nil
}
