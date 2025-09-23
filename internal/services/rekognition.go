package services

import (
	"context"
	"errors"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/rekognition/types"
)

const faceComparingErr = "face comparison failed"
const mismatchingFacesErr = "faces don't match"

type RekognitionService interface {
	CompareFaces(ctx context.Context, idSrc string, faceSrc string) (bool, error)
}

type rekognitionService struct {
	bucket                 string
	faceComparisonTreshold float32
	logger                 *slog.Logger
	rekognitionClient      *rekognition.Client
}

func (r *rekognitionService) CompareFaces(ctx context.Context, idSrc string, faceSrc string) (bool, error) {
	output, err := r.rekognitionClient.CompareFaces(ctx, &rekognition.CompareFacesInput{
		SourceImage: &types.Image{
			S3Object: &types.S3Object{
				Bucket: aws.String(r.bucket),
				Name:   aws.String(faceSrc),
			},
		},
		TargetImage: &types.Image{
			S3Object: &types.S3Object{
				Bucket: aws.String(r.bucket),
				Name:   aws.String(idSrc),
			},
		},
	})

	if err != nil {
		r.logger.Error(faceComparingErr, "error", err)
		return false, err
	}

	compared := len(output.FaceMatches) > 0

	if !compared || aws.ToFloat32(output.FaceMatches[0].Similarity) < r.faceComparisonTreshold {
		r.logger.Error(mismatchingFacesErr, "error", faceComparingErr)
		return false, errors.New(mismatchingFacesErr)
	}

	return true, nil
}

func NewRekognitionService(bucket string, faceComparisonTreshold float32, cfg aws.Config, logger *slog.Logger) RekognitionService {
	return &rekognitionService{
		bucket:                 bucket,
		faceComparisonTreshold: faceComparisonTreshold,
		logger:                 logger,
		rekognitionClient:      rekognition.NewFromConfig(cfg),
	}
}
