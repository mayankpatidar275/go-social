package main

import "net/http"

func (app *applicaion) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))

	// This is how we are using the repository pattern
	// app.store.Posts.Create(r.Context());
}
