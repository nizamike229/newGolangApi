package main

import (
	"awesomeProject/controllers/authController"
	"awesomeProject/controllers/taskController"
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
	db, err := gorm.Open(postgres.Open(os.Getenv("DB_DSN")), &gorm.Config{})
	if err != nil {
		logger.Error("Failed to connect to database: " + err.Error())
		os.Exit(1)
	}
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(customMiddleware.WithDB(db))

	r.Route("/api/task", func(r chi.Router) {
		r.Use(customMiddleware.AuthMiddleware)
		r.Get("/all", taskController.GetAllPersonalTasks)
		r.Post("/create", taskController.CreateTask)
		r.Delete("/deleteById", taskController.DeleteTask)
	})
	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/register", authController.Register)
		r.Post("/login", authController.Login)
	})

	logger.Info("Server is running on port 8080")
	launchErr := http.ListenAndServe(":8080", r)
	if launchErr != nil {
		logger.Error(err.Error())
	}
}
