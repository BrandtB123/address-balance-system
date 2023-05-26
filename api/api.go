package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"unit410/models"
)

type API interface {
	GetData() error
}

func HttpRequest[T models.Test](method, URL string, requestBody interface{}) (*T, error) {
	var req *http.Request
	var err error

	if method == "GET" {
		req, err = http.NewRequest(http.MethodGet, URL, nil)
		if err != nil {
			return nil, err
		}
	} else if method == "POST" {
		jsonBody, err := json.Marshal(requestBody)
		if err != nil {
			return nil, err
		}

		req, err = http.NewRequest(http.MethodPost, URL, bytes.NewBuffer(jsonBody))
		if err != nil {
			return nil, err
		}

		req.Header.Set("Content-Type", "application/json")
	} else {
		return nil, fmt.Errorf("unsupported HTTP method: %s", method)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	for resp.StatusCode == 429 {
		time.Sleep(5 * time.Second)
		resp, err = client.Do(req)
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var responseObject T
	err = json.Unmarshal(body, &responseObject)
	if err != nil {
		return nil, err
	}

	return &responseObject, nil
}
