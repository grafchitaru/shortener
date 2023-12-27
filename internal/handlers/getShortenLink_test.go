package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/grafchitaru/shortener/internal/config"
	"github.com/grafchitaru/shortener/internal/storage/mocks"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func TestGetShorten(t *testing.T) {
	mockStorage := &mocks.MockStorage{}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		GetShorten(config.HandlerContext{Repos: mockStorage}, w, r)
	})

	link := Link{URL: "http://test.com"}
	linkJSON, _ := json.Marshal(link)
	req, err := http.NewRequest("POST", "/shorten", bytes.NewBuffer(linkJSON))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("handler returned unexpected content type: got %v want %v",
			contentType, "application/json")
	}

	expectedPattern := regexp.MustCompile(`\{\"result\"\:\"\/\w+\"\}`)
	if !expectedPattern.MatchString(rr.Body.String()) {
		t.Errorf("handler returned unexpected result: got %v", rr.Body.String())
	}
}
