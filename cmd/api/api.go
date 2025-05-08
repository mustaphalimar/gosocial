package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mustaphalimar/go-social/docs"
	"github.com/mustaphalimar/go-social/internal/store"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

type application struct {
	config config
	store  store.Storage
	logger *zap.SugaredLogger
}

type mailConfig struct {
	exp time.Duration
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type config struct {
	addr   string
	db     dbConfig
	env    string
	apiURL string
	mail   mailConfig
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

		// Swagger Doc http://localhost:8080/swagger/doc.json
		docsUrl := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL(docsUrl), //The url pointing to API definition
		))

		// v1/posts
		r.Route("/posts", func(r chi.Router) {

			// POST v1/posts
			r.Post("/", app.createPostHandler)

			// v1/posts/someId
			r.Route("/{postId}", func(r chi.Router) {
				r.Use(app.postsContextMiddleware) // fetches the post and add it to the context of the request

				r.Get("/", app.getPostHandler)       // POST v1/posts/someId
				r.Patch("/", app.updatePostHandler)  // PATCH v1/posts/someId
				r.Delete("/", app.deletePostHandler) // DELETE v1/posts/someId
			})

		})

		r.Route("/users", func(r chi.Router) {
			r.Route("/{userId}", func(r chi.Router) {
				r.Use(app.userContextMiddleware)
				r.Get("/", app.getUserHandler)

				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unfollowUserHandler)
			})

			r.Group(func(r chi.Router) {
				r.Get("/feed", app.getUserFeedHandler)
			})
		})

	})

	return r
}

func (app *application) run(mux http.Handler) error {
	// Docs
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = app.config.apiURL
	docs.SwaggerInfo.BasePath = "/v1"

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	app.logger.Infow("Server has started", "addr", app.config.addr, "env", app.config.env)

	return srv.ListenAndServe()
}
