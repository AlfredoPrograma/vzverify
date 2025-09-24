package handlers

import (
	"net/http"

	"github.com/alfredoprograma/vzverify/internal/services"
	"github.com/labstack/echo/v4"
)

type UploadHandler interface {
	GeneratePresignedUpload(c echo.Context) error
}

type uploadHandler struct {
	s3Service services.S3Service
}

func (u *uploadHandler) GeneratePresignedUpload(c echo.Context) error {
	uploadDir := services.UploadDir(c.Param("dir"))

	if err := uploadDir.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	url, key, err := u.s3Service.GeneratePresignedUpload(c.Request().Context(), uploadDir)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"url": url,
		"key": key,
	})
}

func NewUploadHandler(s3Service services.S3Service) UploadHandler {
	return &uploadHandler{
		s3Service: s3Service,
	}
}
