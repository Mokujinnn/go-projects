package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"todo-api/internal/handlers"
	"todo-api/internal/models"
	"todo-api/internal/service"
	"todo-api/internal/storage"

	_ "todo-api/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	port = flag.Int("p", 8080, "port for the server")
)

// @title           Todo List API
// @version         1.0
// @description     REST API для управления задачами (To-Do List)

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth
func main() {
	flag.Parse()

	storage := storage.NewMemoryStorage()
	todoService := service.NewTodoService(storage)
	todoHandler := handlers.NewTodoHandler(todoService)

	seedData(storage)

	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		tasks := v1.Group("/tasks")
		{
			tasks.GET("", todoHandler.GetTasks)
			tasks.POST("", todoHandler.CreateTask)
			tasks.GET("/:id", todoHandler.GetTask)
			tasks.PUT("/:id", todoHandler.UpdateTask)
			tasks.DELETE("/:id", todoHandler.DeleteTask)
			tasks.PATCH("/:id/complete", todoHandler.CompleteTask)
		}
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Server starting on port %d", *port)
	log.Fatal(router.Run(addr))
}

func seedData(storage *storage.MemoryStorage) {
	tasks := []struct {
		title       string
		description string
		completed   bool
	}{
		{"Задача 1", "Задача 1", false},
		{"Задача 2", "Задача 2", false},
		{"Задача 3", "Задача 3", true},
		{"Задача 4", "Задача 4", false},
		{"Задача 5", "Задача 5", false},
	}

	for _, task := range tasks {
		newTask := models.Task{
			Title:       task.title,
			Description: task.description,
			Completed:   task.completed,
		}
		_, err := storage.Create(newTask)
		if err != nil {
			log.Printf("Failed to create task: %v", err)
		}
	}

	log.Printf("Created %d test tasks", len(tasks))
}
