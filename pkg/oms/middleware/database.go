package middleware

import (
	"context"
	"gorm.io/gorm"
	"net/http"
)

func DatabaseMiddleware(db *gorm.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// create new context from `r` request context, and assign key `"user"`
			// to value of `"123"`
			ctx := context.WithValue(r.Context(), "db", db)

			// call the next handler in the chain, passing the response writer and
			// the updated request object with the new context value.
			//
			// note: context.Context values are nested, so any previously set
			// values will be accessible as well, and the new `"user"` key
			// will be accessible from this point forward.
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
