package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mustaphalimar/go-social/internal/store"
)

const POST_KEY string = "post"

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

// createPostHandler godoc
//
//	@Summary		Creates a new post
//	@Description	Creates a new post with title, content, and tags
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			post	body		CreatePostPayload	true	"Post payload"
//	@Success		201		{object}	store.Post
//	@Failure		400		{object}	error	"Invalid input"
//	@Failure		500		{object}	error	"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/posts [post]
func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var postPayload CreatePostPayload

	if err := readJSON(w, r, &postPayload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(postPayload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := getUserFromContext(r)
	post := &store.Post{
		Title:   postPayload.Title,
		Content: postPayload.Content,
		Tags:    postPayload.Tags,
		UserID:  user.ID,
	}

	ctx := r.Context()

	if err := app.store.Posts.Create(ctx, post); err != nil {
		app.internalServerResponse(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, post); err != nil {
		app.internalServerResponse(w, r, err)
		return
	}
}

// getPostHandler godoc
//
//	@Summary		Get a post by ID
//	@Description	Retrieves a post by its ID, including its comments
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Post ID"
//	@Success		200	{object}	store.Post
//	@Failure		500	{object}	error	"Internal server error"
//	@Router			/posts/{id} [get]
func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	comments, err := app.store.Comments.GetByPostId(r.Context(), post.ID)
	if err != nil {
		app.internalServerResponse(w, r, err)
		return
	}

	post.Comments = comments

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerResponse(w, r, err)
	}
}

type UpdatePostPayload struct {
	Title   string `json:"title" validate:"omitempty,max=100"`
	Content string `json:"content" validate:"omitempty,max=1000"`
}

// updatePostHandler godoc
//
//	@Summary		Update a post
//	@Description	Updates a post's title and/or content by ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"Post ID"
//	@Param			post	body		UpdatePostPayload	true	"Updated post payload"
//	@Success		200		{object}	store.Post
//	@Failure		400		{object}	error	"Invalid input"
//	@Failure		500		{object}	error	"Internal server error"
//	@Router			/posts/{id} [put]
func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	var payload UpdatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if payload.Content != "" {
		post.Content = payload.Content
	}

	if payload.Title != "" {
		post.Title = payload.Title
	}

	err := app.store.Posts.Update(r.Context(), post)
	if err != nil {

		app.internalServerResponse(w, r, err)
		return
	}

	if err = app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerResponse(w, r, err)
		return
	}
}

// deletePostHandler godoc
//
//	@Summary		Delete a post
//	@Description	Deletes a post by its ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postId	path	int	true	"Post ID"
//	@Success		204		"No Content"
//	@Failure		400		{object}	error	"Invalid post ID"
//	@Failure		404		{object}	error	"Post not found"
//	@Failure		500		{object}	error	"Internal server error"
//	@Router			/posts/{postId} [delete]
func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postId")
	postId, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	ctx := r.Context()

	if err := app.store.Posts.Delete(ctx, postId); err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerResponse(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) postsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "postId")
		postId, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}
		ctx := r.Context()

		post, err := app.store.Posts.GetById(ctx, int64(postId))
		if err != nil {
			switch {
			case errors.Is(err, store.ErrorNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerResponse(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, POST_KEY, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromCtx(r *http.Request) *store.Post {
	post, _ := r.Context().Value(POST_KEY).(*store.Post)
	return post
}
