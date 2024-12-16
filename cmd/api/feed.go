package main

import "net/http"

func (app *applicaion) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	// pagination, filter

	ctx := r.Context()

	feed, err := app.store.Posts.GetUserFeed(ctx, int64(66))
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
	}
}
