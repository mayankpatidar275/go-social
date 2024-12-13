package main

import (
	"log"

	"github.com/mayankpatidar275/go-social/internal/env"
)

// This file will mostly have the configurations
func main(){
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
	}

	app := &applicaion{
		config: cfg,
	}

	mux := app.mount();
	log.Fatal(app.run(mux));
}