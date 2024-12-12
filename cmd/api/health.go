package main

import "net/http"

func (app *applicaion) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}