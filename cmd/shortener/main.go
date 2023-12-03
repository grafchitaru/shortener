package main

import (
	"github.com/grafchitaru/shortener/internal/app"
	"github.com/grafchitaru/shortener/internal/storage/sqlite"
	"io"
	"net/http"
	"strings"
)

func mainPage(res http.ResponseWriter, req *http.Request) {
	storage, err := sqlite.New("././internal/storage/storage.db")
	if err != nil {
		panic(err)
	}

	if req.Method == http.MethodPost {
		reqBody, err := io.ReadAll(req.Body)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte(err.Error()))
			return
		}
		alias := app.NewRandomString(6)

		storage.SaveURL(string(reqBody), alias)
		res.WriteHeader(http.StatusCreated)
		res.Write([]byte("http://localhost:8080/" + alias))
		return
	}

	if req.Method == http.MethodGet {
		path := strings.TrimPrefix(req.URL.Path, "/")
		alias, err := storage.GetURL(path)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte("Bad Request"))
		}

		res.Header().Set("Location", alias)
		res.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	res.WriteHeader(http.StatusBadRequest)
	res.Write([]byte("Bad Request"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, mainPage)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
