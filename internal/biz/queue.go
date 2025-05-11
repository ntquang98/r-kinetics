package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type QueueRepo interface {
	SendJob(context.Context, interface{}) error
}

type QueueUsecase struct {
	repo QueueRepo
	log  *log.Helper
}

func NewQueueUsecase(repo QueueRepo, logger log.Logger) *QueueUsecase {
	return &QueueUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (uc *QueueUsecase) SendJob(ctx context.Context, msg interface{}) error {
	return uc.repo.SendJob(ctx, msg)
}
