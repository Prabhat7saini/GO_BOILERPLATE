package HttpService

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

type HttpService interface {
	Get(url string, headers map[string]string) ([]byte, error)
	Post(url string, headers map[string]string, body []byte) ([]byte, error)
	Patch(url string, headers map[string]string, body []byte) ([]byte, error)
}

type httpServiceImpl struct{}

func NewHttpService() HttpService {
	return &httpServiceImpl{}
}

func (h *httpServiceImpl) Get(url string, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Printf("Failed to create GET request: %v", err)
		return nil, err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("GET request failed: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (h *httpServiceImpl) Post(url string, headers map[string]string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Failed to create POST request: %v", err)
		return nil, err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("POST request failed: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (h *httpServiceImpl) Patch(url string, headers map[string]string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Failed to create PATCH request: %v", err)
		return nil, err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("PATCH request failed: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
