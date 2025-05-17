package main

import (
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/mustaphalimar/go-social/internal/db"
	"github.com/mustaphalimar/go-social/internal/env"
	"github.com/mustaphalimar/go-social/internal/mailer"
	"github.com/mustaphalimar/go-social/internal/store"
	"go.uber.org/zap"
)

const version = "0.0.1"

//	@title			Swagger Example API
//	@description	Go-Social Docs
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath					/v1
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config{
		addr:      env.GetString("ADDR", ":8080"),
		apiURL:    env.GetString("EXTERNAL_URL", "localhost:8080"),
		clientURL: env.GetString("CLIENT_URL", "http://localhost:5173"),
		db: dbConfig{
			addr: env.GetString("DATABASE_URL", "postgresql://postgres:admin@localhost/gosocial?sslmode=disable"),
			// limit number of open connection to the db from our API connection pool
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
		mail: mailConfig{
			exp:       time.Hour * 24 * 3, // days
			apiKey:    env.GetString("SENDGRID_API_KEY", ""),
			fromEmail: env.GetString("FROM_EMAIL", ""),
		},
		auth: authConfig{
			basic: basicConfig{
				username: env.GetString("BASIC_AUTH_USERNAME", ""),
				password: env.GetString("BASIC_AUTH_PASSWORD", ""),
			},
		},
	}
	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// Database
	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	logger.Info("Database connection pool established.")

	store := store.NewStorage(db)

	mailer := mailer.NewSendgrid(cfg.mail.apiKey, cfg.mail.fromEmail)

	app := &application{
		config: cfg,
		store:  store,
		logger: logger,
		mailer: mailer,
	}

	mux := app.mount()
	logger.Fatal(app.run(mux))

}
