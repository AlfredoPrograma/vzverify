package main

import (
	"context"
	"fmt"
	"time"

	"github.com/alfredoprograma/vzverify/internal/config"
	"github.com/alfredoprograma/vzverify/internal/observability"
	"github.com/alfredoprograma/vzverify/internal/services"
)

func main() {
	env := config.MustLoadEnv()
	logger := observability.NewLogger(env.LogLevel)
	awsCfg := config.MustLoadAWSConfig(context.TODO(), logger)
	s3Service := services.NewS3Service(env.IdentitiesBucket, awsCfg, logger)
	textractService := services.NewTextractService(env.IdentitiesBucket, awsCfg, logger)
	rekognitionService := services.NewRekognitionService(env.IdentitiesBucket, env.FaceComparisonTreshold, awsCfg, logger)

	idUploadUrl, idKey, _ := s3Service.GeneratePresignedUpload(context.Background(), services.UploadIdsDir)
	faceUploadUrl, faceKey, _ := s3Service.GeneratePresignedUpload(context.Background(), services.UploadFacesDir)

	fmt.Println(idUploadUrl)
	fmt.Println(faceUploadUrl)

	time.Sleep(time.Second * 10)

	fields, _ := textractService.ExtractIDContent(context.Background(), idKey)
	rekognitionService.CompareFaces(context.Background(), idKey, faceKey)

	fmt.Printf("%#v", fields)
}
