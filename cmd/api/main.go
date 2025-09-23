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

	url, key, _ := s3Service.GeneratePresignedUpload(context.Background(), services.UploadFacesDir)

	fmt.Println(url)
	time.Sleep(time.Second * 6)

	fields, _ := textractService.ExtractIDContent(context.Background(), key)

	fmt.Printf("%#v", fields)
}
