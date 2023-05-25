package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"unit410/models"
)

type API interface {
	GetData() ([]*models.Bal, error)
}

func aasciiCodesToString(asciiCodes []int) string {
	var runes []rune

	for _, code := range asciiCodes {
		runes = append(runes, rune(code))
	}

	return string(runes)
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
		return nil, fmt.Errorf("Unsupported HTTP method: %s", method)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	for resp.Status == "429" {
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
	for strings.Contains(string(body), "Just a moment...") {
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(body, &responseObject)
	}
	for strings.Contains(string(body), "error code: 1015") {
		time.Sleep(3 * time.Second)
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(body, &responseObject)
	}

	return &responseObject, nil
}
