package taskController

import (
	"awesomeProject/logger"
	"awesomeProject/models"
	"encoding/json"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
)

func GetAllPersonalTasks(w http.ResponseWriter, r *http.Request) {
	var db = r.Context().Value("db").(*gorm.DB)
	userId := r.Context().Value("userID").(uuid.UUID)
	var tasks []models.Task
	if err := db.Find(&tasks); err.Error != nil {
		logger.Error("Failed to fetch tasks")
		http.Error(w, "Failed to fetch tasks", http.StatusInternalServerError)
		return
	}

	var personalTasks []models.Task
	for i := range tasks {
		if tasks[i].UserId == userId {
			personalTasks = append(personalTasks, tasks[i])
		}
	}
	tasksJson, err := json.Marshal(personalTasks)
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
		logger.Error("Invalid request: " + err.Error())
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	userId := r.Context().Value("userID").(uuid.UUID)

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

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	var db = r.Context().Value("db").(*gorm.DB)
	userId := r.Context().Value("userID").(uuid.UUID)
	var taskId int
	err := json.NewDecoder(r.Body).Decode(&taskId)
	if err != nil {
		logger.Error("Invalid request: " + err.Error())
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var task models.Task
	if err := db.First(&task, taskId).Error; err != nil {
		logger.Error("Failed to fetch task: " + err.Error())
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if task.UserId != userId {
		logger.Error("User(" + userId.String() + ") is not authorized to delete this task")
		http.Error(w, "You are not authorized to delete this task", http.StatusUnauthorized)
		return
	}
	if err := db.Delete(&task).Error; err != nil {
		logger.Error("Failed to delete task: " + err.Error())
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	logger.Info("Deleted task: " + task.Title)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Task was deleted successfully!"))
}
