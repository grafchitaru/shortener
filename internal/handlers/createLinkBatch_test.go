package handlers

import (
	"encoding/json"
	"github.com/grafchitaru/shortener/internal/config"
	"github.com/grafchitaru/shortener/internal/mocks"
	"github.com/grafchitaru/shortener/internal/storage"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateLinkBatch(t *testing.T) {
	mockStorage := &mocks.MockStorage{
		SaveURLError: nil,
	}
	cfg := mocks.NewConfig()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		CreateLinkBatch(config.HandlerContext{Config: *cfg, Repos: mockStorage}, w, r)
	})

	batchURLs := `[{"originalURL":"http://test1.com","correlationID":"1"}, {"originalURL":"http://test2.com","correlationID":"2"}]`
	req, err := http.NewRequest("POST", "/api/shorten/batch", strings.NewReader(batchURLs))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	var result []storage.BatchResult
	if err := json.Unmarshal(rr.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}

	for _, r := range result {
		expected := cfg.BaseShortURL

		if !strings.Contains(r.ShortURL, expected) {
			t.Errorf("handler returned unexpected short URL for correlation ID %v: got %v want %v",
				r.CorrelationID, r.ShortURL, expected)
		}
	}
}
