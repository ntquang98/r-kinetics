package data

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/ntquang98/go-rkinetics-service/internal/biz"
	"github.com/ntquang98/go-rkinetics-service/internal/domain"
	"github.com/ntquang98/go-rkinetics-service/internal/pkg/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type analyticsJobRepo struct {
	data     *Data
	instance *Instance[*domain.AnalyticsJob]
	log      *log.Helper
}

func NewAnalyticsJobRepo(data *Data, logger log.Logger) biz.AnalyticsJobRepo {
	instance := NewDBInstance[*domain.AnalyticsJob]("analytics_job")
	instance.ApplyDatabase(data.mongo.Database())

	return &analyticsJobRepo{
		data:     data,
		instance: instance,
		log:      log.NewHelper(logger),
	}
}

func (r *analyticsJobRepo) Save(ctx context.Context, input *domain.AnalyticsJob) (*domain.AnalyticsJob, error) {
	input.Status = domain.AnalyticsJobStatus.Draft
	result, err := r.instance.Create(ctx, input)

	if err != nil {
		return nil, err
	}

	return result[0], nil
}
func (r *analyticsJobRepo) UpdateByID(ctx context.Context, id string, input *domain.AnalyticsJob) (*domain.AnalyticsJob, error) {
	ID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, errors.BadRequest(common.ErrorCodeInvalidRequest, fmt.Sprintf("%s is not valid ID", id))
	}

	results, err := r.instance.UpdateOne(ctx, domain.AnalyticsJob{ID: &ID}, input)

	if err != nil {
		return nil, err
	}

	return results[0], nil
}

func (r *analyticsJobRepo) FindByID(ctx context.Context, id string) (*domain.AnalyticsJob, error) {
	ID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, errors.BadRequest(common.ErrorCodeInvalidRequest, fmt.Sprintf("%s is not valid ID", id))
	}

	results, err := r.instance.Query(ctx, domain.AnalyticsJob{ID: &ID}, 0, 0, nil)

	if err != nil {
		return nil, err
	}

	return results[0], nil
}
func (r *analyticsJobRepo) ListAll(ctx context.Context, offset, limit int64) ([]*domain.AnalyticsJob, int64, error) {
	results, err := r.instance.Query(ctx, nil, offset, limit, &primitive.D{{"_id", -1}})

	if err != nil {
		return nil, 0, err
	}

	count, err := r.instance.Count(ctx, nil)

	if err != nil {
		return nil, 0, err
	}

	return results, count, nil
}
