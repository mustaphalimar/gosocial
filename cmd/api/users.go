package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mustaphalimar/go-social/internal/store"
)

type userKey string

const userCtx userKey = "user"

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
	}
}

type FollowUser struct {
	UserId int64 `json:"user_id"`
}

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	// extracting the user from the params, /users/{userId}/follow
	userToFollow := getUserFromContext(r)

	// TODO: reverts back to userId from ctx when jwt auth is implemented
	var payload FollowUser
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	ctx := r.Context()

	if err := app.store.Followers.Follow(ctx, userToFollow.ID, payload.UserId); err != nil {
		switch err {
		case store.ErrConflict:
			app.conflictError(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	// extracting the user from the params, /users/{userId}/follow
	userToUnfollow := getUserFromContext(r)

	// TODO: reverts back to auth userId from ctx when jwt auth is implemented
	var payload FollowUser
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
	}
	ctx := r.Context()

	if err := app.store.Followers.Unfollow(ctx, userToUnfollow.ID, payload.UserId); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, err := strconv.ParseInt(chi.URLParam(r, "userId"), 10, 64)
		if err != nil {
			app.badRequestError(w, r, err)
			return
		}

		ctx := r.Context()

		user, err := app.store.Users.GetById(ctx, userId)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrorNotFound):
				app.notFoundError(w, r, err)
				return
			default:
				app.internalServerError(w, r, err)
				return
			}
		}
		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// extracting the user from the params, /users/{userId}/follow
func getUserFromContext(r *http.Request) *store.User {
	user, _ := r.Context().Value(userCtx).(*store.User)
	return user
}
