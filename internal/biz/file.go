package biz

import (
	"context"
	"io"

	"github.com/go-kratos/kratos/v2/log"
)

type FileRepo interface {
	UploadFile(ctx context.Context, filename string, contentType string, fileReader io.Reader) (string, error)
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
func (uc *FileUsecase) UploadFile(ctx context.Context, filename, contentType string, fileReader io.Reader) (string, error) {
	return uc.repo.UploadFile(ctx, filename, contentType, fileReader)
}
