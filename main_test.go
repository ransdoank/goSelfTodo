package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	// Load the environment variables from the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// setupRouter initializes the Gin router with the routes defined in main.go
func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/tasks", CreateTask)
	r.GET("/tasks", GetTasks)
	r.PUT("/tasks/:id", UpdateTask)
	r.DELETE("/tasks/:id", DeleteTask)
	return r
}

// initTestDB initializes a test database connection
func initTestDB() {
	// Get the environment variables for database connection
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database")
	}

	// Migrate the schema for testing
	DB.AutoMigrate(&Task{})
}

// TestCreateTask tests the CreateTask endpoint
func TestCreateTask(t *testing.T) {
	initTestDB() // Initialize test DB
	router := setupRouter()

	task := Task{
		Title:       "This task from TestCreateTask",
		Description: "This is a test task",
		Status:      "Waiting List",
		DueDate:     time.Now().Add(24 * time.Hour),
	}
	jsonValue, _ := json.Marshal(task)
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var createdTask Task
	err := json.Unmarshal(w.Body.Bytes(), &createdTask)
	assert.Nil(t, err)
	assert.Equal(t, task.Title, createdTask.Title)
	assert.Equal(t, task.Description, createdTask.Description)
	assert.Equal(t, task.Status, createdTask.Status)
}

// TestGetTasks tests the GetTasks endpoint
func TestGetTasks(t *testing.T) {
	initTestDB() // Initialize test DB
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/tasks", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var tasks []Task
	err := json.Unmarshal(w.Body.Bytes(), &tasks)
	assert.Nil(t, err)
}

// TestUpdateTask tests the UpdateTask endpoint
func TestUpdateTask(t *testing.T) {
	initTestDB() // Initialize test DB
	router := setupRouter()

	// Create a task first
	task := Task{
		Title:       "This task from TestUpdateTask",
		Description: "This is a test task",
		Status:      "Waiting List",
		DueDate:     time.Now().Add(24 * time.Hour),
	}
	DB.Create(&task)

	// Update the task
	updatedTask := Task{
		Title:       "This task from TestUpdateTask - Updated Task",
		Description: "This is an updated test task",
		Status:      "In Progress",
		DueDate:     time.Now().Add(24 * time.Hour),
	}
	jsonValue, _ := json.Marshal(updatedTask)
	req, _ := http.NewRequest("PUT", "/tasks/"+fmt.Sprint(task.ID), bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var returnedTask Task
	err := json.Unmarshal(w.Body.Bytes(), &returnedTask)
	assert.Nil(t, err)
	assert.Equal(t, updatedTask.Title, returnedTask.Title)
	assert.Equal(t, updatedTask.Description, returnedTask.Description)
	assert.Equal(t, updatedTask.Status, returnedTask.Status)
}

// TestDeleteTask tests the DeleteTask endpoint
func TestDeleteTask(t *testing.T) {
	initTestDB() // Initialize test DB
	router := setupRouter()

	// Create a task first
	task := Task{
		Title:       "This task from TestDeleteTask",
		Description: "This is a test task",
		Status:      "Waiting List",
		DueDate:     time.Now(), //.Add(24 * time.Hour),
	}
	DB.Create(&task)

	req, _ := http.NewRequest("DELETE", "/tasks/"+fmt.Sprint(task.ID), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Check if the task is soft-deleted
	var deletedTask Task
	err := DB.Unscoped().Where("id = ?", task.ID).First(&deletedTask).Error
	assert.Nil(t, err)
	assert.NotNil(t, deletedTask.DeletedAt)
}

/** Optional For Check Invalid Response uncomment this session ** /
// TestInvalidCreateTask tests the CreateTask endpoint with invalid input
func TestInvalidCreateTask(t *testing.T) {
	initTestDB() // Initialize test DB
	router := setupRouter()

	// Missing required fields
	task := Task{
		Description: "This task has no title",
	}
	jsonValue, _ := json.Marshal(task)
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestUpdateNonexistentTask tests the UpdateTask endpoint with a nonexistent task ID
func TestUpdateNonexistentTask(t *testing.T) {
	initTestDB() // Initialize test DB
	router := setupRouter()

	updatedTask := Task{
		Title:       "Nonexistent Task",
		Description: "This task does not exist",
		Status:      "In Progress",
	}
	jsonValue, _ := json.Marshal(updatedTask)
	req, _ := http.NewRequest("PUT", "/tasks/9999", bytes.NewBuffer(jsonValue)) // Assume 9999 does not exist
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestDeleteNonexistentTask tests the DeleteTask endpoint with a nonexistent task ID
func TestDeleteNonexistentTask(t *testing.T) {
	initTestDB() // Initialize test DB
	router := setupRouter()

	req, _ := http.NewRequest("DELETE", "/tasks/9999", nil) // Assume 9999 does not exist
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
/****/
