package taskController

import (
	"awesomeProject/logger"
	"awesomeProject/models"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"net/http"
	"os"
)

var jwtSecret = getJwt()

func GetAllPersonalTasks(w http.ResponseWriter, r *http.Request) {
	var db = r.Context().Value("db").(*gorm.DB)
	username, err := getUsernameFromCookie(r)
	userId, err := getUserIdFromUsername(r, db, username)
	if err != nil {
		http.Error(w, "Unauthorized:"+err.Error(), http.StatusInternalServerError)
		return
	}
	var tasks []models.Task
	if err := db.Where("user_id = ?", userId).Order("created_at desc").Find(&tasks).Error; err.Error != nil {
		logger.Error("Failed to fetch tasks")
		http.Error(w, "Failed to fetch tasks", http.StatusInternalServerError)
		return
	}

	tasksJson, err := json.Marshal(tasks)
	if err != nil {
		logger.Error("Failed to marshal tasks: " + err.Error())
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(tasksJson)
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var db = r.Context().Value("db").(*gorm.DB)
	var taskRequest models.TaskRequest
	err := json.NewDecoder(r.Body).Decode(&taskRequest)
	if err != nil {
		logger.Error("Failed to unmarshal request: " + err.Error())
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	username, err := getUsernameFromCookie(r)
	if err != nil {
		http.Error(w, "Unauthorized:"+err.Error(), http.StatusInternalServerError)
		logger.Error("Failed to get username: " + err.Error())
		return
	}

	userId, err := getUserIdFromUsername(r, db, username)
	if err != nil {
		http.Error(w, "Unauthorized:"+err.Error(), http.StatusInternalServerError)
		logger.Error("Failed to get userId: " + err.Error())
		return
	}

	task := models.Task{
		Title:     taskRequest.Title,
		Priority:  taskRequest.Priority,
		Deadline:  taskRequest.Deadline,
		Completed: taskRequest.Completed,
		UserId:    userId,
	}

	if err := db.Table("tasks").Create(&task).Error; err != nil {
		logger.Error("Failed to create task: " + err.Error())
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info("Created task: " + task.Title)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Task was created successfully!"))
}

func DeletePerson(w http.ResponseWriter, r *http.Request) {
	//TODO
}

func getUsernameFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("authToken")
	if err != nil {
		logger.Error("No auth token found: " + err.Error())
		return "", err
	}

	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		logger.Error("Invalid token: " + err.Error())
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		logger.Error("Invalid claims")
		return "", fmt.Errorf("invalid token claims")
	}

	username, ok := claims["sub"].(string)
	if !ok {
		logger.Error("Username not found in claims")
		return "", fmt.Errorf("username not found in token")
	}

	return username, nil
}
func getUserIdFromUsername(r *http.Request, db *gorm.DB, username string) (uuid.UUID, error) {
	username, err := getUsernameFromCookie(r)
	if err != nil {
		logger.Error("No auth token found: " + err.Error())
		return uuid.UUID{}, err
	}
	var user models.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		logger.Error("User not found:" + err.Error())
		return uuid.UUID{}, err
	}
	return user.ID, nil
}

var getJwt = func() []byte {
	godotenv.Load("./main.env")
	return []byte(os.Getenv("JWT_SECRET"))
}
