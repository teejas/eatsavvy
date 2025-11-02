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

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		slog.Error("[http.Post] Failed to send HTTP request", "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	slog.Info("[http.Post] Response status", "status", resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("[http.Post] Failed to read response body", "error", err)
		return nil, err
	}

	slog.Info("[http.Post] Response body", "body", string(body))

	return body, nil
}
