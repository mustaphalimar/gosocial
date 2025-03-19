package main

import (
	"log"

	"github.com/mustaphalimar/go-social/internal/db"
	"github.com/mustaphalimar/go-social/internal/env"
	"github.com/mustaphalimar/go-social/internal/store"
)

func main() {
	addr := env.GetString("DATABASE_URL", "postgresql://postgres:admin@localhost/gosocial?sslmode=disable")

	conn, err := db.New(addr, 30, 30, "15m")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()

	store := store.NewStorage(conn)
	db.Seed(store)
}
