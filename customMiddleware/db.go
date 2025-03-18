package customMiddleware

import (
	"awesomeProject/logger"
	"awesomeProject/models"
	"context"
	"gorm.io/gorm"
	"net/http"
)

func WithDB(db *gorm.DB) func(http.Handler) http.Handler {
	dbErr := db.Table("persons").AutoMigrate(&models.Person{})
	usersDbErr := db.Table("users").AutoMigrate(&models.User{})
	if usersDbErr != nil {
		logger.Error(usersDbErr.Error())
	}
	if dbErr != nil {
		logger.Error(dbErr.Error())
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "db", db)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
