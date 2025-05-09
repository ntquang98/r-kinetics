package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	pb "github.com/ntquang98/go-rkinetics-service/api/app/v1"
	"github.com/ntquang98/go-rkinetics-service/internal/biz"
	"github.com/ntquang98/go-rkinetics-service/internal/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AnalyticsJobService struct {
	pb.UnimplementedAnalyticsJobServer

	uc  *biz.AnalyticsJobUsecase
	log *log.Helper
}

func NewAnalyticsJobService(uc *biz.AnalyticsJobUsecase, logger log.Logger) *AnalyticsJobService {
	return &AnalyticsJobService{
		uc:  uc,
		log: log.NewHelper(logger),
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
		return nil, err
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
