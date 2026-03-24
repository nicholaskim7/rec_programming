package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/nicholaskim7/rec-programming/internal/auth"
)

type ContextKey string

const UserIDKey ContextKey = "userID"

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tokenString string

		// try to get the cookie
		cookie, err := r.Cookie("auth_token")
		if err == nil {
			// cookie found
			tokenString = cookie.Value
		}
		if tokenString == "" {
			// if no token try the auth header
			authHeader := r.Header.Get("Authorization")
			if len(authHeader) > 7 && strings.ToUpper(authHeader[:7]) == "BEARER " {
				tokenString = authHeader[7:]
			}
		}
		if tokenString == "" {
			http.Error(w, "Unauthorized: Please log in", http.StatusUnauthorized)
			return
		}
		// validate token
		userID, err := auth.ValidateToken(tokenString)
		if err != nil {
			// token invalid
			http.Error(w, "Unauthorized: Invalid session", http.StatusUnauthorized)
			return
		}
		// store userID in the request context
		// tells next handler who this is
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}