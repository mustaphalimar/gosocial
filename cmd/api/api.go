package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mustaphalimar/go-social/internal/store"
)

type application struct {
	config config
	store  store.Storage
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type config struct {
	addr string
	db   dbConfig
	env  string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// middlewares
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(60 * time.Second))

	// routers
	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		// v1/posts
		r.Route("/posts", func(r chi.Router) {

			// POST v1/posts
			r.Post("/", app.createPostHandler)

			// v1/posts/someId
			r.Route("/{postId}", func(r chi.Router) {

				r.Get("/", app.getPostHandler) // POST v1/posts/someId
				// r.Patch("", app.updatePostHandler)  // PATCH v1/posts/someId
				r.Delete("/", app.deletePostHandler) // DELETE v1/posts/someId
			})

		})
	})

	return r
}

func (app *application) run(mux http.Handler) error {

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Server has started at %s", app.config.addr)

	return srv.ListenAndServe()
}
