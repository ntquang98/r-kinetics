package data

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/ntquang98/go-rkinetics-service/internal/biz"
	"github.com/ntquang98/go-rkinetics-service/internal/conf"
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
func (r *s3FileRepo) UploadFile(ctx context.Context, file *multipart.FileHeader, contentType string) (string, error) {
	f, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	// Generate a unique key (e.g., uploads/12345-filename.jpg)
	key := fmt.Sprintf("uploads/%s-%s", generateUniqueID(), filepath.Base(file.Filename))
	uploader := manager.NewUploader(r.client)
	_, err = uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(r.bucket),
		Key:         aws.String(key),
		Body:        f,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	fileURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", r.bucket, r.region, key)
	r.log.Infof("Uploaded file %s to S3 with URL %s", file.Filename, fileURL)
	return fileURL, nil
}

func generateUniqueID() string {
	return primitive.NewObjectID().Hex()
}
