package vapi

import (
	"eatsavvy/internal/places"
	"eatsavvy/pkg/http"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
)

type VapiClient struct {
	httpClient *http.Http
}

func NewVapiClient() *VapiClient {
	httpClient := http.NewClient()
	return &VapiClient{
		httpClient: httpClient,
	}
}

type VapiCallResponse struct {
	Id string `json:"id"`
}

func (v *VapiClient) CreateCall(restaurant places.Restaurant) (VapiCallResponse, error) {
	reqBody := getAssistantRequestBody(restaurant)
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + os.Getenv("VAPI_API_KEY"),
	}
	respBody, statusCode, err := v.httpClient.Post("https://api.vapi.ai/call", reqBody, headers)
	if err != nil {
		slog.Error("[vapi.CreateCall] Failed to send HTTP request", "error", err)
		return VapiCallResponse{}, err
	}
	if statusCode >= 400 {
		slog.Error("[vapi.CreateCall] Failed to create Vapi call", "statusCode", statusCode, "responseBody", string(respBody))
		return VapiCallResponse{}, errors.New("failed to create Vapi call")
	}
	var vapiResponse VapiCallResponse
	err = json.Unmarshal(respBody, &vapiResponse)
	if err != nil {
		slog.Error("[vapi.CreateCall] Failed to unmarshal response body", "error", err)
		return VapiCallResponse{}, err
	}
	return vapiResponse, nil
}
