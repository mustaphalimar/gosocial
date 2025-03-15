package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/mustaphalimar/go-social/internal/env"
	"github.com/mustaphalimar/go-social/internal/store"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
	}

	store := store.NewStorage(nil)
	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))

}
