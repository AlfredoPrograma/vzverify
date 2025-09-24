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
	textractService services.TextractService
	vzIdService     services.VzIdService
}

func (i *idMatchHandler) Compare(c echo.Context) error {
	key := c.Param("key")

	identityFields, err := i.textractService.ExtractIDContent(c.Request().Context(), key)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	matches, err := i.vzIdService.CompareIdData(c.Request().Context(), identityFields)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"matches": matches,
	})
}

func NewIdMatchHandler(textractService services.TextractService, vzIdService services.VzIdService) IdMatchHandler {
	return &idMatchHandler{
		textractService: textractService,
		vzIdService:     vzIdService,
	}
}
