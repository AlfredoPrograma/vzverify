package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

const idVerificationRequestErr = "id verification request failed"
const idMismatchErr = "collected data from id img and from api mismatchs"

type IdData struct {
	Nationality string
	Names       string
	LastNames   string
	IdNumber    string
}

type IdDataResponse struct {
	Nationality    string `json:"nacionalidad"`
	FirstName      string `json:"primer_nombre"`
	MiddleName     string `json:"segundo_nombre"`
	FirstLastName  string `json:"primer_apellido"`
	SecondLastName string `json:"segundo_apellido"`
	IdNumber       int    `json:"cedula"`
}

type IdResponse struct {
	ErrorStr any            `json:"error_str"`
	Error    bool           `json:"error"`
	Data     IdDataResponse `json:"data"`
}

type VzIdService interface {
	CompareIdData(ctx context.Context, idData IdData) (bool, error)
}

type vzIdService struct {
	appId      string
	apiToken   string
	apiUrl     string
	logger     *slog.Logger
	httpClient *http.Client
}

func (v *vzIdService) compare(collectedFromImg IdData, collectedFromApi IdDataResponse) bool {
	collectedNames := strings.Join([]string{collectedFromApi.FirstName, collectedFromApi.MiddleName}, " ")
	collectedLastNames := strings.Join([]string{collectedFromApi.FirstLastName, collectedFromApi.SecondLastName}, " ")

	hasSameNames := strings.EqualFold(collectedFromImg.Names, collectedNames)
	hasSameLastNames := strings.EqualFold(collectedFromImg.LastNames, collectedLastNames)
	hasSameNationality := strings.EqualFold(collectedFromImg.Nationality, collectedFromApi.Nationality)
	hasSameId := collectedFromImg.IdNumber == fmt.Sprintf("%d", collectedFromApi.IdNumber)

	return hasSameNames && hasSameLastNames && hasSameId && hasSameNationality
}

func (v *vzIdService) CompareIdData(ctx context.Context, idData IdData) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, v.apiUrl, nil)

	if err != nil {
		v.logger.Debug("cannot build request", "error", err)
		return false, err
	}

	qs := req.URL.Query()
	qs.Add("app_id", v.appId)
	qs.Add("token", v.apiToken)
	qs.Add("nacionalidad", idData.Nationality)
	qs.Add("cedula", idData.IdNumber)
	req.URL.RawQuery = qs.Encode()

	resp, err := v.httpClient.Do(req)

	if err != nil {
		v.logger.Error("cannot request id data", "error", err)
		return false, err
	}

	var data IdResponse

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		v.logger.Error("cannot decode response body", "error", err)
		return false, err
	}

	if data.Error {
		v.logger.Error(idVerificationRequestErr, "error", data.ErrorStr)
		return false, errors.New(idVerificationRequestErr)
	}

	areSame := v.compare(idData, data.Data)

	if !areSame {
		v.logger.Error(idMismatchErr, "error", idMismatchErr, "dataFromImg", idData, "dataFromApi", data)
		return false, nil
	}

	return true, nil
}

func NewVzIdService(apiUrl string, appId string, apiToken string, logger *slog.Logger) VzIdService {
	return &vzIdService{
		appId:      appId,
		apiToken:   apiToken,
		apiUrl:     apiUrl,
		logger:     logger,
		httpClient: &http.Client{},
	}
}
