package main

import (
	"fmt"
	"net/http"
	"strings"
)

// requestLogger logs all requests
func (a *application) requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log request
		a.infoLogger.Printf("%s::%s {%s} - %s %s", r.RemoteAddr, r.Proto, r.UserAgent(), r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

// secureHeaders inserts custom headers for security
func (a *application) secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// insert security security headers
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")
		next.ServeHTTP(w, r)
	})
}

func (a *application) auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check auth header
		if !isAuthorized(r) {
			// a.landingPage(w, r)
			http.Redirect(w, r, "/auth", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// check for auth
func isAuthorized(r *http.Request) bool {
	authHeader := r.Header.Get("Authorization")
	if strings.Contains(r.URL.Path, "static") || strings.Contains(r.URL.Path, "auth") {
		return true
	}
	return authHeader == "" // temp. To correct and flesh out
}

// panicRecovery recovers from panics
func (a *application) panicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				a.serverError(w, r, fmt.Errorf("internal error: %s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
