package main

import (
	"log"
	"time"

	"todo-api/internal/models"
	"todo-api/internal/storage"
)

func main() {
	storage := storage.NewMemoryStorage()

	tasks := []models.CreateTaskRequest{
		{Title: "Изучить основы Go", Description: "Изучить синтаксис и основные концепции языка Go"},
		{Title: "Написать REST API", Description: "Создать RESTful API с использованием Gin框架"},
		{Title: "Изучить работу с базой данных", Description: "Разобраться с подключением и запросами к БД"},
		{Title: "Написать тесты", Description: "Создать unit-тесты для всех компонентов приложения"},
		{Title: "Документировать API", Description: "Создать Swagger документацию для API"},
		{Title: "Оптимизировать производительность", Description: "Провести профилирование и оптимизацию кода"},
		{Title: "Настроить CI/CD", Description: "Настроить автоматическую сборку и деплой"},
		{Title: "Изучить Docker", Description: "Научиться контейнеризации приложения"},
		{Title: "Создать фронтенд", Description: "Разработать интерфейс для работы с API"},
		{Title: "Протестировать безопасность", Description: "Провести security testing API"},
		{Title: "Завершенная задача 1", Description: "Эта задача уже выполнена"},
		{Title: "Завершенная задача 2", Description: "Еще одна выполненная задача"},
	}

	for i, taskReq := range tasks {
		task := models.Task{
			Title:       taskReq.Title,
			Description: taskReq.Description,
			Completed:   i >= 10, // Последние 2 задачи выполнены
			CreatedAt:   time.Now().Add(-time.Duration(i) * time.Hour),
			UpdatedAt:   time.Now().Add(-time.Duration(i) * time.Hour),
		}

		_, err := storage.Create(task)
		if err != nil {
			log.Printf("Failed to create task %d: %v", i+1, err)
		} else {
			log.Printf("Created task: %s", task.Title)
		}
	}

	log.Printf("Successfully seeded %d tasks", len(tasks))
}
