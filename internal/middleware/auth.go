package middleware

import (
	"context"
	"nds-go-starter/internal/core/auth"
	"nds-go-starter/internal/json"
	"net/http"
	"strings"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func Auth(jwtManager *auth.JWTManager, enabled bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !enabled {
				next.ServeHTTP(w, r)
				return
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				json.Write(w, r, http.StatusUnauthorized, "missing authorization header")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				json.Write(w, r, http.StatusUnauthorized, "invalid authorization header format")
				return
			}

			accessToken := parts[1]
			claims, err := jwtManager.ValidateToken(accessToken)
			if err != nil {
				json.Write(w, r, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
