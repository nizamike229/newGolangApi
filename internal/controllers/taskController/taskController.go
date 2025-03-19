package taskController

import (
	"awesomeProject/internal/logger"
	"awesomeProject/internal/models"
	"encoding/json"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
)

// GetAllPersonalTasks возвращает все задачи пользователя
// @Summary Get all personal tasks
// @Description Retrieve a list of all tasks for the authenticated user
// @Tags tasks
// @Accept json
// @Produce json
// @Success 200 {array} string "List of tasks"
// @Failure 401 {object} string "Unauthorized"
// @Router /api/task/all [get]
func GetAllPersonalTasks(w http.ResponseWriter, r *http.Request) {
	var db = r.Context().Value("db").(*gorm.DB)
	userId := r.Context().Value("userID")
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

// CreateTask создаёт новую задачу
// @Summary Create a new task
// @Description Create a new task for the authenticated user
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body models.TaskRequest true "Task description"
// @Success 201 {object} string "Task created"
// @Failure 400 {object} string "Bad request"
// @Failure 401 {object} string "Unauthorized"
// @Router /api/task/create [post]
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
		Completed: false,
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

// DeleteTask удаляет задачу по ID
// @Summary Delete a task
// @Description Delete a task by its ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id query int true "Task ID"
// @Success 200 {object} string "Task deleted"
// @Failure 400 {object} string "Bad request"
// @Failure 401 {object} string "Unauthorized"
// @Router /api/task/deleteById [delete]
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

// CompleteTask помечает задачу как завершённую
// @Summary Complete a task
// @Description Mark a task as completed
// @Tags tasks
// @Accept json
// @Produce json
// @Param id body int true "Task ID"
// @Success 200 {object} string "Task completed"
// @Failure 400 {object} string "Bad request"
// @Failure 401 {object} string "Unauthorized"
// @Router /api/task/complete [patch]
func CompleteTask(w http.ResponseWriter, r *http.Request) {
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
		logger.Error("User(" + userId.String() + ") is not authorized to complete this task")
		http.Error(w, "You are not authorized to complete this task", http.StatusUnauthorized)
		return
	}
	task.Completed = true
	if err := db.Save(&task).Error; err != nil {
		logger.Error("Failed to complete task: " + err.Error())
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Task was completed successfully!"))
}
