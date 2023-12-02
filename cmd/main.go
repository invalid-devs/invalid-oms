package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/invalid-devs/invalid-oms/pkg/oms/controllers"
	"github.com/invalid-devs/invalid-oms/pkg/oms/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
)

func main() {
	dsn := "host=localhost user=root password=123 dbname=invalid_oms port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("database connection fail")
	}

	db.AutoMigrate(&models.User{})

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(DatabaseMiddleware(db))

	r.Get("/v1/user/{id}", controllers.GetUser)
	r.Post("/v1/user/", controllers.CreateUser)
	http.ListenAndServe("0.0.0.0:3000", r)
}

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
