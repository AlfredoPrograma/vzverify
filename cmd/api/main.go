package main

import (
	"context"

	"github.com/alfredoprograma/vzverify/internal/config"
	"github.com/alfredoprograma/vzverify/internal/observability"
	"github.com/alfredoprograma/vzverify/internal/services"
)

func main() {
	env := config.MustLoadEnv()
	logger := observability.NewLogger(env.LogLevel)
	awsCfg := config.MustLoadAWSConfig(context.TODO(), logger)
	services.NewS3Service(env.IdentitiesBucket, awsCfg, logger)
}
