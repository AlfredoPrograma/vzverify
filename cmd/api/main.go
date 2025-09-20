package main

import (
	"context"

	"github.com/alfredoprograma/vzverify/internal/config"
	"github.com/alfredoprograma/vzverify/internal/observability"
)

func main() {
	logger := observability.NewLogger("debug")
	config.NewAWSConfig(context.TODO(), logger)

}
