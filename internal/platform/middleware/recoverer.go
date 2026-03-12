package middleware

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

func Recoverer(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {

			if err := recover(); err != nil {

				log.Error().
					Interface("panic", err).
					Msg("panic recovered")

				http.Error(w, "internal server error", http.StatusInternalServerError)
			}

		}()

		next.ServeHTTP(w, r)
	})
}
