package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/invalid-devs/invalid-oms/pkg/oms/controllers"
	omsMiddleware "github.com/invalid-devs/invalid-oms/pkg/oms/middleware"
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
	r.Use(omsMiddleware.DatabaseMiddleware(db))

	// Public routes
	r.Group(func(r chi.Router) {
		r.Post("/api/v1/auth/sign-in/", controllers.SignIn)
		r.Post("/api/v1/auth/sign-out/", controllers.SignIn)
	})

	// Private routes
	r.Group(func(r chi.Router) {
		r.Use(omsMiddleware.Authorization)

		r.Get("/api/v1/user/{id}/", controllers.GetUser)
		r.Post("/api/v1/user/", controllers.CreateUser)
	})

	http.ListenAndServe("0.0.0.0:3000", r)
}
