package main

import (
	"net/http"
)

const (
	ColorRed   = "\033[31m"
	ColorReset = "\033[0m"
)

func (app *application) internalServerResponse(w http.ResponseWriter, r *http.Request, err error) {
	var message = "INTERNAL_SERVER_ERROR"
	app.logger.Errorw(message, "method", r.Method, "path", r.URL.Path, err)
	writeJSONError(w, http.StatusInternalServerError, "The server encountered a problem while procession your request.")
}

func (app *application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	var message = "CONFLICT_ERROR"
	app.logger.Warnf(message, "method", r.Method, "path", r.URL.Path, err)
	writeJSONError(w, http.StatusConflict, err.Error())
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	var message = "BAD_REQUEST_ERROR"
	app.logger.Warnf(message, "method", r.Method, "path", r.URL.Path, err)
	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) forbiddenResponse(w http.ResponseWriter, r *http.Request) {
	var message = "FORBIDDEN_ERROR"
	app.logger.Warnf(message, "method", r.Method, "path", r.URL.Path)
	writeJSONError(w, http.StatusForbidden, "This action is forbidden")
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	var message = "NOT_FOUND_ERROR"
	app.logger.Warnf(message, "method", r.Method, "path", r.URL.Path, err)
	writeJSONError(w, http.StatusNotFound, "Resource not found.")
}

func (app *application) unauthorizedResponse(w http.ResponseWriter, r *http.Request, err error) {
	var message = "UNAUTHORIZED_ERROR"
	app.logger.Errorw(message, "method", r.Method, "path", r.URL.Path, err)

	writeJSONError(w, http.StatusUnauthorized, err.Error())
}
func (app *application) unauthorizedBasicAuthResponse(w http.ResponseWriter, r *http.Request, err error) {
	var message = "UNAUTHORIZED_BASIC_AUTH_ERROR"
	app.logger.Errorw(message, "method", r.Method, "path", r.URL.Path, err)

	w.Header().Set("WWW-Authenticate", `Basic realm="restriced", charset="UTF-8"`)
	writeJSONError(w, http.StatusUnauthorized, err.Error())
}
