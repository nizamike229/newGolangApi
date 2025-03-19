package main

import (
	_ "awesomeProject/docs"
	"awesomeProject/internal/controllers/authController"
	"awesomeProject/internal/controllers/taskController"
	customMiddleware2 "awesomeProject/internal/customMiddleware"
	"awesomeProject/internal/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"os"
)

// @title Simple Golang TO-DO
// @version 1.0
// @description This is a simple API for managing tasks and authentication
// @host localhost:8080
// @BasePath /
func main() {
	godotenv.Load("./configs/main.env")
	db, err := gorm.Open(postgres.Open(os.Getenv("DB_DSN")), &gorm.Config{})
	if err != nil {
		logger.Error("Failed to connect to database: " + err.Error())
		os.Exit(1)
	}
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(customMiddleware2.WithDB(db))

	r.Route("/api/task", func(r chi.Router) {
		r.Use(customMiddleware2.AuthMiddleware)
		r.Get("/all", taskController.GetAllPersonalTasks)
		r.Post("/create", taskController.CreateTask)
		r.Delete("/deleteById", taskController.DeleteTask)
		r.Patch("/complete", taskController.CompleteTask)
	})
	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/register", authController.Register)
		r.Post("/login", authController.Login)
	})

	r.Get("/docs/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/swagger.json")
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/docs/swagger.json"),
	))

	logger.Info("Server is running on port 8080")
	launchErr := http.ListenAndServe(":8080", r)
	if launchErr != nil {
		logger.Error(err.Error())
	}
}
