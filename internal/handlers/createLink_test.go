package handlers

import (
	"github.com/grafchitaru/shortener/internal/config"
	"github.com/grafchitaru/shortener/internal/storage/mocks"
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
	cfg := config.NewConfig()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		CreateLink(w, r, mockStorage, cfg)
	})

	req, err := http.NewRequest("POST", "/create", strings.NewReader("http://test.com"))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	expected := cfg.HTTPServerAddress
	if rr.Body.String()[:len(expected)] != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
