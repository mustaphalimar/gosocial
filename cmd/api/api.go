package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/mustaphalimar/go-social/docs"
	"github.com/mustaphalimar/go-social/internal/auth"
	"github.com/mustaphalimar/go-social/internal/mailer"
	"github.com/mustaphalimar/go-social/internal/store"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

type application struct {
	config        config
	store         store.Storage
	logger        *zap.SugaredLogger
	mailer        mailer.Client
	authenticator auth.Authenticator
}

type mailConfig struct {
	exp       time.Duration
	apiKey    string
	fromEmail string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type config struct {
	addr      string
	db        dbConfig
	env       string
	apiURL    string
	mail      mailConfig
	clientURL string
	auth      authConfig
}

type authConfig struct {
	basic basicConfig
	jwt   jwtConfig
}

type jwtConfig struct {
	secret string
	exp    time.Duration
	iss    string
}

type basicConfig struct {
	username string
	password string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// middlewares
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(60 * time.Second))

	// routers
	r.Route("/v1", func(r chi.Router) {
		r.With(app.BasicAuthMiddleware()).Get("/health", app.healthCheckHandler)

		// Swagger Doc http://localhost:8080/swagger/doc.json
		docsUrl := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL(docsUrl), //The url pointing to API definition
		))

		// v1/posts
		r.Route("/posts", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			// POST v1/posts
			r.Post("/", app.createPostHandler)

			// v1/posts/someId
			r.Route("/{postId}", func(r chi.Router) {
				r.Use(app.postsContextMiddleware) // fetches the post and add it to the context of the request

				// POST v1/posts/someId
				r.Get("/", app.getPostHandler)
				// PATCH v1/posts/someId
				r.Patch("/", app.checkPostsOwnership("moderator", app.updatePostHandler))
				// DELETE v1/posts/someId
				r.Delete("/", app.checkPostsOwnership("admin", app.deletePostHandler))
			})

		})

		r.Route("/users", func(r chi.Router) {
			r.Put("/activate/{token}", app.activateUserHandler)

			r.Route("/{userId}", func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)
				r.Get("/", app.getUserHandler)

				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unfollowUserHandler)
			})

			r.Group(func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)
				r.Get("/feed", app.getUserFeedHandler)
			})
		})

		// auth routes
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", app.registerUserHandler)
			r.Post("/token", app.createTokenHandler)
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
