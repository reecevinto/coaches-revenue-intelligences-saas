package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/reecevinto/coaches-revenue-intelligences-saas/internal/platform/jwt"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func Auth(jwtService *jwt.Service) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				http.Error(w, "missing token", 401)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			claims, err := jwtService.ValidateToken(tokenStr)
			if err != nil {
				http.Error(w, "invalid token", 401)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
