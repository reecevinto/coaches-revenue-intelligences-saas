package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/reecevinto/coaches-revenue-intelligences-saas/internal/platform/jwt"
)

// unique context key type to avoid collisions
type userContextKey string

// key used to store user id in request context
const UserIDKey userContextKey = "user_id"

// Auth middleware validates JWT tokens and injects the user ID into request context
func Auth(jwtService *jwt.Service) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			// Expect: Authorization: Bearer <token>
			parts := strings.Split(authHeader, " ")

			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
				return
			}

			tokenStr := parts[1]

			claims, err := jwtService.ValidateToken(tokenStr)
			if err != nil {
				http.Error(w, "invalid or expired token", http.StatusUnauthorized)
				return
			}

			// attach user ID to request context
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
