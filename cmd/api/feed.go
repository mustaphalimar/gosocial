package main

import (
	"net/http"

	"github.com/mustaphalimar/go-social/internal/store"
)

// getUserFeedHandler godoc
//
//	@Summary		Get user feed
//	@Description	Retrieves a paginated, filtered, and sorted feed of posts for the user
//	@Tags			feed
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int		false	"Limit the number of posts"
//	@Param			offset	query		int		false	"Offset for pagination"
//	@Param			sort	query		string	false	"Sort order (asc or desc)"
//	@Param			tags	query		string	false	"Tags"
//	@Param			search	query		string	false	"Search"
//	@Success		200		{array}		store.Post
//	@Failure		400		{object}	error	"Invalid query parameters"
//	@Failure		500		{object}	error	"Internal server error"
//	@Router			/users/feed [get]
func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	// pagination, filters, sort
	fq := store.PaginatedFeedQuery{
		Limit:  10,
		Offset: 0,
		Sort:   "desc",
	}

	fq, err := fq.Parse(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(fq); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	feed, err := app.store.Posts.GetUserFeed(ctx, int64(100), fq)

	if err != nil {
		app.internalServerResponse(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerResponse(w, r, err)
	}

}
