package main

import (
	"awesomeProject/controllers/authController"
	"awesomeProject/controllers/personController"
	"awesomeProject/customMiddleware"
	"awesomeProject/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"os"
)

func main() {
	godotenv.Load("./main.env")
	var db, _ = gorm.Open(postgres.Open(os.Getenv("DB_DSN")))
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(customMiddleware.WithDB(db))

	r.Route("/api/persons", func(r chi.Router) {
		r.Use(customMiddleware.AuthMiddleware)
		r.Get("/all", personController.GetAllPersons)
		r.Post("/create", personController.CreatePerson)
		r.Delete("/deleteById", personController.DeletePerson)
	})
	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/register", authController.Register)
		r.Post("/login", authController.Login)
	})

	logger.Info("Server is running on port 8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		logger.Error(err.Error())
	}
}
