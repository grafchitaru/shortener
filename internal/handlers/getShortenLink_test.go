package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/grafchitaru/shortener/internal/auth"
	"github.com/grafchitaru/shortener/internal/config"
	"github.com/grafchitaru/shortener/internal/mocks"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetShorten(t *testing.T) {
	cfg := mocks.NewConfig()
	token, err := auth.GenerateToken(uuid.New(), cfg.SecretKey)
	require.NoError(t, err)

	mockStorage := &mocks.MockStorage{
		GetAliasError:  nil,
		GetAliasResult: "testalias",
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		GetShorten(config.HandlerContext{Config: *cfg, Repos: mockStorage}, w, r)
	})

	link := Link{
		URL: "http://test.com",
	}
	linkJSON, _ := json.Marshal(link)
	req, err := http.NewRequest("POST", "/api/shorten", bytes.NewBuffer(linkJSON))
	require.NoError(t, err)

	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: token,
		Path:  "/",
	})

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusConflict {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusConflict)
	}

	expected := "{\"result\":\"http://127.0.0.1:8080/testalias\"}"
	if body := rr.Body.String(); body != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			body, expected)
	}
}
