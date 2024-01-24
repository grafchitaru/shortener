package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/grafchitaru/shortener/internal/config"
	"github.com/grafchitaru/shortener/internal/mocks"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetShorten(t *testing.T) {
	mockStorage := &mocks.MockStorage{
		GetAliasError:  nil,
		GetAliasResult: "testalias",
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		GetShorten(config.HandlerContext{Repos: mockStorage}, w, r)
	})

	link := Link{
		URL: "http://test.com",
	}
	linkJSON, _ := json.Marshal(link)
	req, err := http.NewRequest("POST", "/api/shorten", bytes.NewBuffer(linkJSON))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusConflict {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusConflict)
	}

	expected := "{\"result\":\"/testalias\"}"
	if body := rr.Body.String(); body != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			body, expected)
	}
}
