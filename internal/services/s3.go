package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

const presignPutObjectErr = "presign put object failed"

type UploadDir string

const (
	UploadIdsDir   UploadDir = "ids"
	UploadFacesDir UploadDir = "faces"
)

type S3Service interface {
	GeneratePresignedUpload(ctx context.Context, uploadDir UploadDir) (string, string, error)
}

type s3Service struct {
	bucket        string
	presignClient *s3.PresignClient
	logger        *slog.Logger
}

func (s *s3Service) GeneratePresignedUpload(ctx context.Context, uploadDir UploadDir) (string, string, error) {
	fileId := uuid.NewString()
	keyPath := fmt.Sprintf("%s/%s", uploadDir, fileId)

	presignedUrl, err := s.presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(keyPath),
	})

	if err != nil {
		s.logger.ErrorContext(ctx, presignPutObjectErr, "error", err)
		return "", "", err
	}

	return presignedUrl.URL, keyPath, nil
}

func NewS3Service(bucket string, cfg aws.Config, logger *slog.Logger) S3Service {
	s3Client := s3.NewFromConfig(cfg)

	return &s3Service{
		bucket:        bucket,
		presignClient: s3.NewPresignClient(s3Client),
		logger:        logger,
	}
}
