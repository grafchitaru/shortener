package handlers

import (
	"github.com/grafchitaru/shortener/internal/config"
	"github.com/grafchitaru/shortener/internal/storage/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetLink(t *testing.T) {
	mockStorage := &mocks.MockStorage{
		GetURLError:  nil,
		GetURLResult: "http://test.com",
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		GetLink(config.HandlerContext{Repos: mockStorage}, w, r)
	})

	req, err := http.NewRequest("GET", "/testalias", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusTemporaryRedirect {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusTemporaryRedirect)
	}

	expected := "http://test.com"
	if location := rr.Header().Get("Location"); location != expected {
		t.Errorf("handler returned unexpected location: got %v want %v",
			location, expected)
	}
}
