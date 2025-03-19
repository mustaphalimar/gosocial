package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/mustaphalimar/go-social/internal/db"
	"github.com/mustaphalimar/go-social/internal/env"
	"github.com/mustaphalimar/go-social/internal/store"
)

const version = "0.0.1"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr: env.GetString("DATABASE_URL", "postgresql://postgres:admin@localhost/gosocial?sslmode=disable"),
			// limit number of open connection to the db from our API connection pool
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
	}

	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		log.Panic(err.Error())
	}
	defer db.Close()
	log.Println("Database connection pool established.")

	store := store.NewStorage(db)
	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))

}
