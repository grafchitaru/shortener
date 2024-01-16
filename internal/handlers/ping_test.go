package handlers

import (
	"errors"
	"github.com/grafchitaru/shortener/internal/config"
	"github.com/grafchitaru/shortener/internal/storage/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {
	tests := []struct {
		name           string
		mockStorage    *mocks.MockStorage
		expectedStatus int
	}{
		{
			name: "Error when pinging database",
			mockStorage: &mocks.MockStorage{
				PingError: errors.New("some error"),
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Successful operation",
			mockStorage: &mocks.MockStorage{
				PingError: nil,
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRecorder()
			http.NewRequest("GET", "/ping", nil)
			ctx := config.HandlerContext{Repos: tt.mockStorage}
			Ping(ctx, r)
			if status := r.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}
