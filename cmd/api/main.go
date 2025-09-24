package main

import (
	"context"

	"github.com/alfredoprograma/vzverify/internal/config"
	"github.com/alfredoprograma/vzverify/internal/handlers"
	"github.com/alfredoprograma/vzverify/internal/observability"
	"github.com/alfredoprograma/vzverify/internal/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	env := config.MustLoadEnv()
	logger := observability.NewLogger(env.LogLevel)
	awsCfg := config.MustLoadAWSConfig(context.Background(), logger)

	// Initialize services
	s3Service := services.NewS3Service(env.IdentitiesBucket, awsCfg, logger)
	textractService := services.NewTextractService(env.IdentitiesBucket, awsCfg, logger)
	vzIdService := services.NewVzIdService(env.VZIdApiUrl, env.VZIdApp, env.VZIdToken, logger)
	rekognitionService := services.NewRekognitionService(env.IdentitiesBucket, env.FaceComparisonTreshold, awsCfg, logger)

	// Initialize handlers
	uploadsHandler := handlers.NewUploadHandler(s3Service)
	idMatchHandler := handlers.NewIdMatchHandler(textractService, rekognitionService, vzIdService)

	// Initialize http listening and routes
	srv := echo.New()

	srv.Use(middleware.Logger())
	srv.Use(middleware.AddTrailingSlash())

	apiGroup := srv.Group("/api/v1")

	uploadsGroup := apiGroup.Group("/uploads")
	uploadsGroup.GET("/:dir", uploadsHandler.GeneratePresignedUpload)

	idMatchGroup := apiGroup.Group("/match")
	idMatchGroup.GET("", idMatchHandler.Compare)

	srv.Logger.Fatal(srv.Start(":8080"))

	/* 	idUploadUrl, idKey, _ := s3Service.GeneratePresignedUpload(context.Background(), services.UploadIdsDir)
	   	faceUploadUrl, faceKey, _ := s3Service.GeneratePresignedUpload(context.Background(), services.UploadFacesDir)

	   	fmt.Println(idUploadUrl)
	   	time.Sleep(time.Second * 5)

	   	idData, _ := textractService.ExtractIDContent(context.Background(), idKey)
	   	success, _ := vzIdService.CompareIdData(context.Background(), idData)

	   	rekognitionService.CompareFaces(context.Background(), idKey, faceKey) */
	// fmt.Printf("%#v", fields)
}
