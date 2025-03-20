package customMiddleware

import (
	"awesomeProject/internal/models"
	"context"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
)

func WithDB(db *gorm.DB) func(http.Handler) http.Handler {
	dbErr := db.Table("tasks").AutoMigrate(&models.Task{})
	usersDbErr := db.Table("users").AutoMigrate(&models.User{})
	if usersDbErr != nil {
		logrus.Error(usersDbErr.Error())
	}
	if dbErr != nil {
		logrus.Error(dbErr.Error())
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "db", db)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
