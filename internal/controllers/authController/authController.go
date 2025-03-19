package authController

import (
	"awesomeProject/internal/logger"
	"awesomeProject/internal/models"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"net/http"
	"os"
	"time"
)

var jwtSecret = getJwt()

// Register регистрирует нового пользователя
// @Summary Register a new user
// @Description Register a new user with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.UserRequest true "User data" { "username": "string", "password": "string" }
// @Success 201 {object} string "User registered"
// @Failure 400 {object} string "Bad request"
// @Router /api/auth/register [post]
func Register(w http.ResponseWriter, r *http.Request) {
	var db = r.Context().Value("db").(*gorm.DB)
	var request models.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		logger.Error("Validation Error" + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var existingUser models.User
	err := db.Where("username = ?", request.Username).First(&existingUser).Error
	if err == nil {
		logger.Error("Username already taken: " + request.Username)
		http.Error(w, "Username already taken", http.StatusConflict)
		return
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error("Database error: " + err.Error())
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	user := models.User{Username: request.Username}
	if err := user.HashPassword(request.Password); err != nil {
		http.Error(w, "Failed to hash password", http.StatusBadRequest)
		return
	}

	if err := db.Table("users").Create(&user).Error; err != nil {
		logger.Error("Failed to create user")
		http.Error(w, "Failed to create user", http.StatusBadRequest)
		return
	}

	logger.Info("User registered: @" + user.Username)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully!"})
}

// Login аутентифицирует пользователя
// @Summary Login a user
// @Description Authenticate a user and return a token
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.UserRequest true "User data" { "username": "string", "password": "string" }
// @Success 200 {object} string "Token"
// @Failure 401 {object} string "Unauthorized"
// @Router /api/auth/login [post]
func Login(w http.ResponseWriter, r *http.Request) {
	var db = r.Context().Value("db").(*gorm.DB)
	var request models.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.Error(err.Error())
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := db.Table("users").Where("username = ?", request.Username).First(&user).Error; err != nil {
		logger.Error(err.Error())
		http.Error(w, "Invalid credentials", http.StatusNotFound)
		return
	}

	if !user.CheckPassword(request.Password) {
		http.Error(w, "Invalid credentials", http.StatusNotFound)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":    user.Username,
		"userId": user.ID.String(),
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "Failed to sign token", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "authToken",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	})

	logger.Info("User logged in: @" + user.Username)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "Logged in successfully!"})
}

var getJwt = func() []byte {
	godotenv.Load("./main.env")
	return []byte(os.Getenv("JWT_SECRET"))
}
