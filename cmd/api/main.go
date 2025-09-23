package main

import (
	"context"
	"fmt"
	"strings"
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
	services.NewRekognitionService(env.IdentitiesBucket, env.FaceComparisonTreshold, awsCfg, logger)
	vzIdService := services.NewVzIdService(env.VZIdApiUrl, env.VZIdApp, env.VZIdToken, logger)

	idUploadUrl, idKey, _ := s3Service.GeneratePresignedUpload(context.Background(), services.UploadIdsDir)
	// faceUploadUrl, faceKey, _ := s3Service.GeneratePresignedUpload(context.Background(), services.UploadFacesDir)

	fmt.Println(idUploadUrl)
	time.Sleep(time.Second * 5)

	fields, _ := textractService.ExtractIDContent(context.Background(), idKey)

	idData := services.IdData{
		Nationality: "v",
		IdNumber:    strings.ReplaceAll(fields["v"], ".", ""),
		Names:       fields["nombres"],
		LastNames:   fields["apellidos"],
	}

	success, _ := vzIdService.CompareIdData(context.Background(), idData)
	fmt.Println("soy arrechisimo: ", success)
	// rekognitionService.CompareFaces(context.Background(), idKey, faceKey)
	// fmt.Printf("%#v", fields)
}
