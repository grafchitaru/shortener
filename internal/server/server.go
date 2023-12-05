package server

import "net/http"

func Server() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, MainPage)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
