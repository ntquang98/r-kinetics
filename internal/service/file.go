package service

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	pb "github.com/ntquang98/go-rkinetics-service/api/file/v1"
	"github.com/ntquang98/go-rkinetics-service/internal/biz"
	"github.com/ntquang98/go-rkinetics-service/internal/pkg/common"
)

type FileService struct {
	pb.UnimplementedFileServer

	uc  *biz.FileUsecase
	log *log.Helper
}

func NewFileService(uc *biz.FileUsecase, logger log.Logger) *FileService {
	return &FileService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

// gRPC implement
func (s *FileService) UploadFile(ctx context.Context, req *pb.UploadFileRequest) (*pb.UploadFileReply, error) {
	return &pb.UploadFileReply{
		FileUrl: "Not use gRPC to upload",
	}, nil
}

func (s *FileService) UploadFileHTTP(ctx context.Context, req *http.Request) (string, error) {
	file, handler, err := req.FormFile("file")
	if err != nil {
		return "", errors.BadRequest(common.ErrorCodeInvalidRequest, "can not extract file "+err.Error())
	}
	defer file.Close()

	var maxSize int64 = 10 * 1024 * 1024
	if handler.Size > maxSize {
		return "", errors.BadRequest(common.ErrorCodeInvalidRequest, fmt.Sprintf("file size exceeds %d bytes", maxSize))
	}

	filename := filepath.Base(handler.Filename)
	contentType := handler.Header.Get("Content-Type")

	if !strings.HasPrefix(contentType, "image/") && !strings.HasPrefix(contentType, "video/") {
		return "", errors.BadRequest(common.ErrorCodeInvalidRequest, fmt.Sprintf("invalid file type: %s, only images and videos allowed", contentType))
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".mp4":  true,
		".avi":  true,
	}
	if !allowedExts[ext] {
		return "", errors.BadRequest(common.ErrorCodeInvalidRequest, fmt.Sprintf("invalid file extension: %s", ext))
	}

	return s.uc.UploadFile(ctx, filename, contentType, file)
}
