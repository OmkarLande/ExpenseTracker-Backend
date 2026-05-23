package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"ExpenseTracker-Backend/internal/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
)

type S3Storage struct {
	client     *s3.Client
	presigner  *s3.PresignClient
	bucketName string
}

func NewS3Storage(cfg *config.Config) (*S3Storage, error) {
	if cfg.AWSRegion == "" || cfg.AWSS3Bucket == "" {
		return nil, errors.New("missing S3 configurations (AWS_REGION or AWS_S3_BUCKET)")
	}

	var awsCfg aws.Config
	var err error

	if cfg.AWSAccessKeyID != "" && cfg.AWSSecretAccessKey != "" {
		awsCfg, err = awsconfig.LoadDefaultConfig(context.TODO(),
			awsconfig.WithRegion(cfg.AWSRegion),
			awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
				cfg.AWSAccessKeyID,
				cfg.AWSSecretAccessKey,
				"",
			)),
		)
	} else {
		awsCfg, err = awsconfig.LoadDefaultConfig(context.TODO(),
			awsconfig.WithRegion(cfg.AWSRegion),
		)
	}

	if err != nil {
		return nil, fmt.Errorf("unable to load AWS configuration: %w", err)
	}

	s3Client := s3.NewFromConfig(awsCfg)
	presignClient := s3.NewPresignClient(s3Client)

	return &S3Storage{
		client:     s3Client,
		presigner:  presignClient,
		bucketName: cfg.AWSS3Bucket,
	}, nil
}

func (s *S3Storage) GenerateUploadURL(ctx context.Context, key string, contentType string) (string, error) {
	input := &s3.PutObjectInput{
		Bucket:        aws.String(s.bucketName),
		Key:           aws.String(key),
		ContentType:   aws.String(contentType),
	}

	presignedReq, err := s.presigner.PresignPutObject(ctx, input, func(opts *s3.PresignOptions) {
		opts.Expires = 15 * time.Minute
	})
	if err != nil {
		return "", fmt.Errorf("failed to presign put object request: %w", err)
	}

	return presignedReq.URL, nil
}

func (s *S3Storage) ObjectExists(ctx context.Context, key string) (bool, int64, error) {
	input := &s3.HeadObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	}

	resp, err := s.client.HeadObject(ctx, input)
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			if apiErr.ErrorCode() == "NotFound" || apiErr.ErrorCode() == "NoSuchKey" {
				return false, 0, nil
			}
		}
		return false, 0, err
	}

	size := aws.ToInt64(resp.ContentLength)
	return true, size, nil
}

func (s *S3Storage) DeleteObject(ctx context.Context, key string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	}

	_, err := s.client.DeleteObject(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete object from S3: %w", err)
	}

	return nil
}
