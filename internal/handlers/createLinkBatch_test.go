package handlers

import (
	"encoding/json"
	"github.com/grafchitaru/shortener/internal/config"
	"github.com/grafchitaru/shortener/internal/mocks"
	"github.com/grafchitaru/shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code, "handler returned wrong status code")
	assert.Equal(t, rr.Header().Get("Content-Type"), "application/json", "handler returned wrong header content type")

	var result []storage.BatchResult
	if err := json.Unmarshal(rr.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}

	for _, r := range result {
		assert.Contains(t, r.ShortURL, cfg.BaseShortURL, "handler returned unexpected short URL for correlation")
	}
}
