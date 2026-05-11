package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/rearurides/eagle-bank/pkg/token"
)

type contextKey string

const UserIDKey contextKey = "userId"

func Auth(tm *token.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"message":"access token is missing"}`, http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, `{"message":"invalid authorization header format"}`, http.StatusUnauthorized)
				return
			}

			claims, err := tm.Verify(parts[1])
			if err != nil {
				http.Error(w, `{"message":"access token is invalid or expired"}`, http.StatusUnauthorized)
				return
			}

			// store userId in context for handlers to use
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID retrieves the userId from the request context
func GetUserID(r *http.Request) (string, bool) {
	id, ok := r.Context().Value(UserIDKey).(string)
	if !ok || id == "" {
		return "", false
	}
	return id, true
}
