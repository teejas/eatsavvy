package http

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
)

type Http struct {
	client *http.Client
}

func NewClient() *Http {
	return &Http{
		client: &http.Client{},
	}
}

func (h *Http) Get(url string, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		slog.Error("[http.Get] Failed to create HTTP request", "error", err)
		return nil, err
	}

	body, err := sendRequest(h.client, req, headers)
	if err != nil {
		slog.Error("[http.Get] Failed to send HTTP request", "error", err)
		return nil, err
	}

	return body, nil
}

func (h *Http) Post(url string, reqBody map[string]string, headers map[string]string) ([]byte, error) {
	jsonReqBody, err := json.Marshal(reqBody)
	if err != nil {
		slog.Error("[http.Post] Failed to marshal request body", "error", err)
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonReqBody))
	if err != nil {
		slog.Error("[http.Post] Failed to create HTTP request", "error", err)
		return nil, err
	}

	body, err := sendRequest(h.client, req, headers)
	if err != nil {
		slog.Error("[http.Post] Failed to send HTTP request", "error", err)
		return nil, err
	}

	return body, nil
}

func sendRequest(httpClient *http.Client, req *http.Request, headers map[string]string) ([]byte, error) {
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		slog.Error("[http.sendRequest] Failed to send HTTP request", "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	slog.Info("[http.sendRequest] Response status", "status", resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("[http.sendRequest] Failed to read response body", "error", err)
		return nil, err
	}

	return body, nil
}
