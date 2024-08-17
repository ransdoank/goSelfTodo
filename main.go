package main

import (
	"log"
	"os"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/joho/godotenv"
)

var DB *gorm.DB

func init() {
	// Load the environment variables from the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	// Get the environment variables for database connection
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Construct the DSN (Data Source Name) for MySQL
	dsn := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB = db

	// Migrate the schema
	db.AutoMigrate(&Task{})

	r := gin.Default()

	r.POST("/tasks", CreateTask)
	r.GET("/tasks", GetTasks)
	r.PUT("/tasks/:id", UpdateTask)
	r.DELETE("/tasks/:id", DeleteTask)

	r.Run() // default on :8080
}

type Task struct {
	ID          uint `gorm:"primaryKey"`
	Title       string
	Description string
	Status      string
	DueDate     time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func CreateTask(c *gin.Context) {
	var task Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	DB.Create(&task)
	c.JSON(http.StatusOK, task)
}

func GetTasks(c *gin.Context) {
	var tasks []Task
	query := DB

	// Filter by status if the query parameter is provided
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if dueDate := c.Query("due_date"); dueDate != "" {
		// Parse the date from the query parameter
		parsedDate, err := time.Parse("2006-01-02", dueDate) // assuming date is provided in YYYY-MM-DD format
		if err == nil {
			// Convert the parsedDate to the string format expected by MySQL
			formattedDate := parsedDate.Format("2006-01-02")
			query = query.Where("DATE(due_date) = ?", formattedDate)
		}
	}

	query.Find(&tasks)
	c.JSON(http.StatusOK, tasks)
}

func UpdateTask(c *gin.Context) {
	var task Task
	if err := DB.Where("id = ?", c.Param("id")).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	DB.Save(&task)
	c.JSON(http.StatusOK, task)
}

func DeleteTask(c *gin.Context) {
	var task Task
	if err := DB.Where("id = ?", c.Param("id")).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	/**
		hard delete action
	**/
	// DB.Delete(&task)
	/** **/

	/**
		soft delete action
		- Set the DeletedAt field to the current time
		- Save the updated task with the soft delete timestamp
	**/
	task.DeletedAt = gorm.DeletedAt{
		Time:  time.Now(),
		Valid: true,
	}
	DB.Save(&task)

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}
