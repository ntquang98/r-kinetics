package data

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/ntquang98/go-rkinetics-service/internal/biz"
	"github.com/ntquang98/go-rkinetics-service/internal/conf"
	"github.com/ntquang98/go-rkinetics-service/internal/pkg/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// s3FileRepo implements FileRepo for AWS S3
type s3FileRepo struct {
	client *s3.Client
	bucket string
	region string
	log    *log.Helper
}

// NewS3FileRepo creates a new FileRepo for S3
func NewS3FileRepo(conf *conf.Data, logger log.Logger) biz.FileRepo {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(conf.S3.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			conf.S3.Access, conf.S3.Secret, "")),
	)
	if err != nil {
		log.NewHelper(logger).Errorf("failed to load AWS config: %v", err)
		panic(err)
	}
	client := s3.NewFromConfig(cfg)
	return &s3FileRepo{
		client: client,
		bucket: conf.S3.Bucket,
		region: conf.S3.Region,
		log:    log.NewHelper(logger),
	}
}

// UploadFile uploads a file to S3 and returns its public URL
func (r *s3FileRepo) UploadFile(ctx context.Context, filename string, contentType string, fileReader io.Reader) (string, error) {
	key := fmt.Sprintf("uploads/%s-%s", generateUniqueID(), filename)
	uploader := manager.NewUploader(r.client)
	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(r.bucket),
		Key:         aws.String(key),
		Body:        fileReader,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", errors.InternalServer(common.ErrorCodeInternalError, fmt.Sprintf("failed to upload file to S3: %s", err.Error()))
	}

	fileURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", r.bucket, r.region, key)
	return fileURL, nil
}

func generateUniqueID() string {
	return primitive.NewObjectID().Hex()
}
