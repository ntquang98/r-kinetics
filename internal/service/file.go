package service

import (
	"context"
	"net/http"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/ntquang98/go-rkinetics-service/internal/biz"
	"github.com/ntquang98/go-rkinetics-service/internal/pkg/common"
)

type FileService struct {
	fileUsecase *biz.FileUsecase
	log         *log.Helper
}

func NewFileService(fileUsecase *biz.FileUsecase, logger log.Logger) *FileService {
	return &FileService{
		fileUsecase: fileUsecase,
		log:         log.NewHelper(logger),
	}
}

func (s *FileService) UploadFile(ctx context.Context, req *http.Request) *common.APIResponse[map[string]string] {
	err := req.ParseMultipartForm(32 << 20) // 32MB max memory
	if err != nil {
		s.log.WithContext(ctx).Errorf("failed to parse multipart form: %v", err)
		return &common.APIResponse[map[string]string]{
			Status:    common.APIStatus.Error,
			Message:   "failed to parse form",
			ErrorCode: string(common.ErrorCodeInvalidRequest),
		}
	}

	file, handler, err := req.FormFile("file")
	if err != nil {
		s.log.Errorf("failed to get file from form: %v", err)
		return &common.APIResponse[map[string]string]{
			Status:    common.APIStatus.Error,
			Message:   "no file uploaded",
			ErrorCode: string(common.ErrorCodeInvalidRequest),
		}
	}
	defer file.Close()

	fileURL, err := s.fileUsecase.UploadFile(ctx, handler)
	if err != nil {
		s.log.Errorf("failed to upload file: %v", err)
		return &common.APIResponse[map[string]string]{
			Status:    common.APIStatus.Error,
			Message:   err.Error(),
			ErrorCode: string(common.ErrorCodeInternalError),
		}
	}

	return &common.APIResponse[map[string]string]{
		Status:  common.APIStatus.Ok,
		Message: "file uploaded successfully",
		Data:    []map[string]string{{"file_url": fileURL}},
	}
}
