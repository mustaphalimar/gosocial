package main

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
)

func (app *application) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// read the auth header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.unauthorizedBasicAuthResponse(w, r, errors.New("Authorization header is missing."))
				return
			}
			// parse it -> get the base64
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Basic" {
				app.unauthorizedBasicAuthResponse(w, r, errors.New("Authorization header malformed."))
				return
			}

			// decode the base64
			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				app.unauthorizedBasicAuthResponse(w, r, err)
				return
			}

			// check the credentials
			username := app.config.auth.basic.username
			password := app.config.auth.basic.password

			creds := strings.SplitN(string(decoded), ":", 2)
			if len(creds) != 2 || creds[0] != username || creds[1] != password {
				app.unauthorizedBasicAuthResponse(w, r, errors.New("Invalid credentials."))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
