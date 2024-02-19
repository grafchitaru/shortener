package handlers

import (
	"bytes"
	"errors"
	"github.com/google/uuid"
	"github.com/grafchitaru/shortener/internal/auth"
	"github.com/grafchitaru/shortener/internal/mocks"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/grafchitaru/shortener/internal/config"
)

func TestGetUserUrls(t *testing.T) {
	cfg := mocks.NewConfig()

	token, err := auth.GenerateToken(uuid.New(), cfg.SecretKey)
	require.NoError(t, err)

	tests := []struct {
		name           string
		mockStorage    *mocks.MockStorage
		expectedStatus int
	}{
		{
			name: "Error when getting URLs",
			mockStorage: &mocks.MockStorage{
				GetURLError: errors.New("some error"),
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name: "Successful operation",
			mockStorage: &mocks.MockStorage{
				GetURLError: nil,
				GetURLResult: `
					[
						{
							"short_url": "http://example.com/1",
							"original_url": "http://example.com"
						},
						{
							"short_url": "http://example.com/2",
							"original_url": "http://example.com/page2"
						}
					]
				`,
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/user/urls", bytes.NewBufferString(""))
			req.AddCookie(&http.Cookie{
				Name:  "token",
				Value: token,
				Path:  "/",
			})

			ctx := config.HandlerContext{Config: *cfg, Repos: tt.mockStorage}
			GetUserUrls(ctx, r, req)
			if status := r.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}
