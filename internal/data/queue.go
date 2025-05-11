package data

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/ntquang98/go-rkinetics-service/internal/biz"
	"github.com/ntquang98/go-rkinetics-service/internal/conf"
	"github.com/ntquang98/go-rkinetics-service/internal/pkg/common"
)

type sqsRepo struct {
	client   *sqs.Client
	queueURL string
	log      *log.Helper
}

func NewSqsRepo(conf *conf.Data, logger log.Logger) biz.QueueRepo {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(conf.S3.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			conf.S3.Access, conf.S3.Secret, "")),
	)

	if err != nil {
		log.NewHelper(logger).Errorf("failed to load AWS config: %v", err)
		panic(err)
	}

	queueURL := fmt.Sprintf("https://sqs.%s.amazonaws.com/%s/%s", conf.S3.Bucket, conf.Sqs.Account, conf.Sqs.Qname)

	client := sqs.NewFromConfig(cfg)

	return &sqsRepo{
		client:   client,
		queueURL: queueURL,
		log:      log.NewHelper(logger),
	}
}

func (s *sqsRepo) SendJob(ctx context.Context, msg interface{}) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return errors.InternalServer(common.ErrorConversion, fmt.Sprintf("failed to marshal message: %s", err.Error()))
	}

	_, err = s.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    &s.queueURL,
		MessageBody: aws.String(string(body)),
	})

	if err != nil {
		return errors.InternalServer(common.ErrorCodeInternalError, fmt.Sprintf("failed to send message to SQS: %s", err.Error()))
	}

	return nil
}
