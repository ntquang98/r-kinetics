package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/ntquang98/go-rkinetics-service/internal/domain"
)

type AnalyticsJobRepo interface {
	Save(context.Context, *domain.AnalyticsJob) (*domain.AnalyticsJob, error)
	UpdateByID(context.Context, string, *domain.AnalyticsJob) (*domain.AnalyticsJob, error)
	FindByID(context.Context, string) (*domain.AnalyticsJob, error)
	ListAll(context.Context, int64, int64) ([]*domain.AnalyticsJob, int64, error)
}

type AnalyticsJobUsecase struct {
	repo AnalyticsJobRepo
	log  *log.Helper
}

func NewAnalyticsJobUsecase(repo AnalyticsJobRepo, logger log.Logger) *AnalyticsJobUsecase {
	return &AnalyticsJobUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (uc *AnalyticsJobUsecase) CreateAnalyticsJob(ctx context.Context, input *domain.AnalyticsJob) (*domain.AnalyticsJob, error) {
	return uc.repo.Save(ctx, input)
}

func (uc *AnalyticsJobUsecase) UpdateAnalyticsJob(ctx context.Context, id string, input *domain.AnalyticsJob) (*domain.AnalyticsJob, error) {
	return uc.repo.UpdateByID(ctx, id, input)
}

func (uc *AnalyticsJobUsecase) GetAnalyticsJob(ctx context.Context, id string) (*domain.AnalyticsJob, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *AnalyticsJobUsecase) ListAll(ctx context.Context, offset, limit int64) ([]*domain.AnalyticsJob, int64, error) {
	return uc.repo.ListAll(ctx, offset, limit)
}
