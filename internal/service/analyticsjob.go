package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	pb "github.com/ntquang98/go-rkinetics-service/api/app/v1"
	"github.com/ntquang98/go-rkinetics-service/internal/biz"
	"github.com/ntquang98/go-rkinetics-service/internal/domain"
	"github.com/ntquang98/go-rkinetics-service/internal/pkg/common"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AnalyticsJobService struct {
	pb.UnimplementedAnalyticsJobServer

	uc           *biz.AnalyticsJobUsecase
	queueUsecase *biz.QueueUsecase
	log          *log.Helper
}

func NewAnalyticsJobService(uc *biz.AnalyticsJobUsecase, queueUsecase *biz.QueueUsecase, logger log.Logger) *AnalyticsJobService {
	return &AnalyticsJobService{
		uc:           uc,
		queueUsecase: queueUsecase,
		log:          log.NewHelper(logger),
	}
}

func (s *AnalyticsJobService) CreateAnalyticsJob(ctx context.Context, req *pb.CreateAnalyticsJobRequest) (*pb.CreateAnalyticsJobReply, error) {
	result, err := s.uc.CreateAnalyticsJob(ctx, &domain.AnalyticsJob{
		Latitude:  &req.Latitude,
		Longitude: &req.Longitude,
		FileUrl:   req.FileUrl,
		VideoUrl:  req.VideoUrl,
	})

	if err != nil {
		s.log.Errorf("can't create analytics job: %s", err.Error())
		return nil, err
	}

	id := result.ID.Hex()

	err = s.queueUsecase.SendJob(ctx, map[string]string{
		"job_id":       id,
		"file_url":     result.FileUrl,
		"video_url":    result.VideoUrl,
		"callback_url": "http://go-api:8000/v1/analytics-job/result",
	})

	if err != nil {
		s.log.Errorf("can't add job for %s", result.ID.Hex())
	} else {
		s.uc.UpdateAnalyticsJob(ctx, id, &domain.AnalyticsJob{
			Status: domain.AnalyticsJobStatus.AssignedJob,
		})

		result.Status = domain.AnalyticsJobStatus.AssignedJob
	}

	return &pb.CreateAnalyticsJobReply{
		Data: mappingDomainAnalyticsJobToPbAnalyticsJob(result),
	}, nil
}

func (s *AnalyticsJobService) GetAnalyticsJob(ctx context.Context, req *pb.GetAnalyticsJobRequest) (*pb.GetAnalyticsJobReply, error) {
	result, err := s.uc.GetAnalyticsJob(ctx, req.Id)

	if err != nil {
		return nil, err
	}

	return &pb.GetAnalyticsJobReply{
		Data: mappingDomainAnalyticsJobToPbAnalyticsJob(result),
	}, nil
}

func (s *AnalyticsJobService) ListAnalyticsJob(ctx context.Context, req *pb.ListAnalyticsJobRequest) (*pb.ListAnalyticsJobReply, error) {
	result, total, err := s.uc.ListAll(ctx, req.Offset, req.Limit)

	if err != nil {
		return nil, err
	}

	data := make([]*pb.AnalyticsJobModel, 0, len(result))
	for _, item := range result {
		data = append(data, mappingDomainAnalyticsJobToPbAnalyticsJob(item))
	}

	return &pb.ListAnalyticsJobReply{
		Total: total,
		Data:  data,
	}, nil
}

func (s *AnalyticsJobService) CompleteAnalyticsJob(ctx context.Context, req *pb.CompleteAnalyticsJobRequest) (*pb.CompleteAnalyticsJobReply, error) {
	if req.Id == "" {
		return nil, errors.BadRequest(common.ErrorCodeInvalidRequest, "id is required")
	}

	nextStatus := domain.AnalyticsJobStatus.Complete
	jobResult := req.Result
	if req.Result == "" && req.Message != "" {
		nextStatus = domain.AnalyticsJobStatus.Error
		jobResult = req.Message
	}

	_, err := s.uc.UpdateAnalyticsJob(ctx, req.Id, &domain.AnalyticsJob{
		Status: nextStatus,
		Result: jobResult,
	})

	if err != nil {
		s.log.Errorf("can't update analytics job: %s", err.Error())
		return nil, errors.InternalServer(common.ErrorCodeInternalError, "can't save request")
	}

	return &pb.CompleteAnalyticsJobReply{
		Message: "OK",
	}, nil
}

func (s *AnalyticsJobService) RePushJob(ctx context.Context, req *pb.RePushJobRequest) (*pb.RePushJobReply, error) {
	result, err := s.uc.GetAnalyticsJob(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	err = s.queueUsecase.SendJob(ctx, map[string]string{
		"job_id":       req.Id,
		"file_url":     result.FileUrl,
		"video_url":    result.VideoUrl,
		"callback_url": "http://go-api:8000/v1/analytics-job/result",
	})

	if err != nil {
		s.log.Errorf("can't add job for %s", result.ID.Hex())
		return nil, err
	} else {
		s.uc.UpdateAnalyticsJob(ctx, req.Id, &domain.AnalyticsJob{
			Status: domain.AnalyticsJobStatus.AssignedJob,
		})

		result.Status = domain.AnalyticsJobStatus.AssignedJob
	}

	return &pb.RePushJobReply{Message: "OK"}, nil
}

func mappingDomainAnalyticsJobToPbAnalyticsJob(data *domain.AnalyticsJob) *pb.AnalyticsJobModel {
	result := pb.AnalyticsJobModel{
		Id:              data.ID.Hex(),
		CreatedTime:     timestamppb.New(*data.CreatedTime),
		LastUpdatedTime: timestamppb.New(*data.LastUpdatedTime),

		FileUrl:  data.FileUrl,
		VideoUrl: data.VideoUrl,
		Status:   string(data.Status),
		Result:   data.Result,
	}

	if data.Latitude != nil {
		result.Latitude = *data.Latitude
	}

	if data.Longitude != nil {
		result.Longitude = *data.Longitude
	}

	if data.CreatedTime != nil {
		result.CreatedTime = timestamppb.New(*data.CreatedTime)
	}

	if data.LastUpdatedTime != nil {
		result.CreatedTime = timestamppb.New(*data.LastUpdatedTime)
	}

	return &result
}
