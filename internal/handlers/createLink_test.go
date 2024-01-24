package handlers

import (
	"github.com/grafchitaru/shortener/internal/config"
	"github.com/grafchitaru/shortener/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateLink(t *testing.T) {
	mockStorage := &mocks.MockStorage{
		SaveURLError: nil,
		SaveURLID:    123,
	}
	cfg := mocks.NewConfig()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		CreateLink(config.HandlerContext{Config: *cfg, Repos: mockStorage}, w, r)
	})

	req, err := http.NewRequest("POST", "/create", strings.NewReader("http://test.com"))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code, "handler returned wrong status code")
	assert.Equal(t, rr.Body.String()[:len(cfg.BaseShortURL)], cfg.BaseShortURL, "handler returned unexpected body")
}
