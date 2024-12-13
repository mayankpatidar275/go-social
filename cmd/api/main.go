package main

import (
	"log"

	"github.com/mayankpatidar275/go-social/internal/env"
	"github.com/mayankpatidar275/go-social/internal/store"
)

// This file will mostly have the configurations
func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
	}

	store := store.NewStorage(nil)

	app := &applicaion{
		config: cfg,
		store:  store,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
