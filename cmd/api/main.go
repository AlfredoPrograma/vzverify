package main

import (
	"context"
	"fmt"

	"github.com/alfredoprograma/vzverify/internal/config"
	"github.com/alfredoprograma/vzverify/internal/observability"
	"github.com/alfredoprograma/vzverify/internal/services"
)

func main() {
	env := config.MustLoadEnv()
	logger := observability.NewLogger(env.LogLevel)
	awsCfg := config.MustLoadAWSConfig(context.TODO(), logger)
	s3Service := services.NewS3Service(env.IdentitiesBucket, awsCfg, logger)

	url, _, _ := s3Service.GeneratePresignedUploadId(context.Background())

	fmt.Println(url)
}
