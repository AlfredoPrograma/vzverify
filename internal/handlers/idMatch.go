package handlers

import (
	"net/http"

	"github.com/alfredoprograma/vzverify/internal/services"
	"github.com/labstack/echo/v4"
)

type IdMatchHandler interface {
	Compare(c echo.Context) error
}

type idMatchHandler struct {
	rekognitionService services.RekognitionService
	textractService    services.TextractService
	vzIdService        services.VzIdService
}

func (i *idMatchHandler) Compare(c echo.Context) error {
	idKey := c.QueryParam("idKey")
	faceKey := c.QueryParam("faceKey")

	identityFields, err := i.textractService.ExtractIDContent(c.Request().Context(), idKey)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	matches, err := i.vzIdService.CompareIdData(c.Request().Context(), identityFields)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	i.rekognitionService.CompareFaces(c.Request().Context(), idKey, faceKey)

	return c.JSON(http.StatusOK, echo.Map{
		"matches": matches,
	})
}

func NewIdMatchHandler(textractService services.TextractService, rekognitionService services.RekognitionService, vzIdService services.VzIdService) IdMatchHandler {
	return &idMatchHandler{
		rekognitionService: rekognitionService,
		textractService:    textractService,
		vzIdService:        vzIdService,
	}
}
