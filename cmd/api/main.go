package main

import "log"

// This file will mostly have the configurations
func main(){
	cfg := config{
		addr: ":8080",
	}

	app := &applicaion{
		config: cfg,
	}

	mux := app.mount();
	log.Fatal(app.run(mux));
}