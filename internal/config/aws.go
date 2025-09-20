package config

import (
	"context"
	"log/slog"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
)

func NewAWSConfig(ctx context.Context, logger *slog.Logger) awsConfig.Config {
	cfg, err := awsConfig.LoadDefaultConfig(ctx)

	if err != nil {
		logger.Info("cannot load default config")
		logger.Debug("error loading AWS config", "error", err)
	}

	return cfg
}
