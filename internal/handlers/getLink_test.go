package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/grafchitaru/shortener/internal/config"
	"github.com/grafchitaru/shortener/internal/storage/mocks"
)

func TestGetLink(t *testing.T) {
	tests := []struct {
		name           string
		mockStorage    *mocks.MockStorage
		expectedStatus int
	}{
		{
			name: "Error when getting URL",
			mockStorage: &mocks.MockStorage{
				SaveURLError:   nil,
				SaveURLID:      1,
				GetURLError:    errors.New("some error"),
				GetURLResult:   "",
				GetAliasResult: "",
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "Successful operation",
			mockStorage: &mocks.MockStorage{
				SaveURLError:   nil,
				SaveURLID:      1,
				GetURLError:    nil,
				GetURLResult:   "http://example.com",
				GetAliasResult: "tUaOlJ",
			},
			expectedStatus: http.StatusTemporaryRedirect,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/tUaOlJ", bytes.NewBufferString(""))
			ctx := config.HandlerContext{Repos: tt.mockStorage}
			GetLink(ctx, r, req)
			if status := r.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}
