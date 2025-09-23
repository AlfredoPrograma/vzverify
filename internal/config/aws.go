package config

import (
	"context"
	"log/slog"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
)

const loadDefaultConfigErr = "load default config failed"

func MustLoadAWSConfig(ctx context.Context, logger *slog.Logger) aws.Config {
	cfg, err := awsConfig.LoadDefaultConfig(ctx)

	if err != nil {
		logger.ErrorContext(ctx, loadDefaultConfigErr, "error", err)
		os.Exit(1)
	}

	return cfg
}
