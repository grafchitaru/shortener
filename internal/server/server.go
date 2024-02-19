package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/grafchitaru/shortener/internal/auth"
	"github.com/grafchitaru/shortener/internal/compress"
	"github.com/grafchitaru/shortener/internal/config"
	"github.com/grafchitaru/shortener/internal/handlers"
	"github.com/grafchitaru/shortener/internal/logger"
	"net/http"
)

func New(ctx config.HandlerContext) {

	r := chi.NewRouter()

	r.Use(logger.WithLogging)
	r.Use(compress.WithCompressionResponse)
	r.Use(auth.WithUserCookie(ctx))

	r.Delete("/api/user/urls", func(res http.ResponseWriter, req *http.Request) {
		handlers.DeleteUserUrls(ctx, res, req)
	})

	r.Get("/api/user/urls", func(res http.ResponseWriter, req *http.Request) {
		handlers.GetUserUrls(ctx, res, req)
	})

	r.Post("/api/shorten/batch", func(res http.ResponseWriter, req *http.Request) {
		handlers.CreateLinkBatch(ctx, res, req)
	})

	r.Get("/ping", func(res http.ResponseWriter, req *http.Request) {
		handlers.Ping(ctx, res)
	})

	r.Get("/{id}", func(res http.ResponseWriter, req *http.Request) {
		handlers.GetLink(ctx, res, req)
	})

	r.Post("/", func(res http.ResponseWriter, req *http.Request) {
		handlers.CreateLink(ctx, res, req)
	})

	r.Post("/api/shorten", func(res http.ResponseWriter, req *http.Request) {
		handlers.GetShorten(ctx, res, req)
	})

	err := http.ListenAndServe(ctx.Config.HTTPServerAddress, r)
	if err != nil {
		fmt.Println("Error server: %w", err)
	}
}
