package service

import (
	"errors"

	"todo-api/internal/models"
	"todo-api/internal/storage"
)

type TodoService struct {
	storage *storage.MemoryStorage
}

func NewTodoService(storage *storage.MemoryStorage) *TodoService {
	return &TodoService{
		storage: storage,
	}
}

func (s *TodoService) CreateTask(req models.CreateTaskRequest) (models.Task, error) {
	task := models.Task{
		Title:       req.Title,
		Description: req.Description,
		Completed:   false,
	}

	return s.storage.Create(task)
}

func (s *TodoService) GetTask(id string) (models.Task, error) {
	if err := validateUUID(id); err != nil {
		return models.Task{}, err
	}

	return s.storage.GetByID(id)
}

func (s *TodoService) GetAllTasks(query models.TaskQuery) (models.TasksResponse, error) {
	if query.Limit <= 0 {
		query.Limit = 10
	}
	if query.Offset < 0 {
		query.Offset = 0
	}

	tasks, total, err := s.storage.GetAll(query)
	if err != nil {
		return models.TasksResponse{}, err
	}

	return models.TasksResponse{
		Tasks:  tasks,
		Total:  total,
		Limit:  query.Limit,
		Offset: query.Offset,
	}, nil
}

func (s *TodoService) UpdateTask(id string, req models.UpdateTaskRequest) (models.Task, error) {
	if err := validateUUID(id); err != nil {
		return models.Task{}, err
	}

	existing, err := s.storage.GetByID(id)
	if err != nil {
		return models.Task{}, err
	}

	// Обновляем только переданные поля
	if req.Title != "" {
		existing.Title = req.Title
	}
	if req.Description != "" {
		existing.Description = req.Description
	}
	if req.Completed != nil {
		existing.Completed = *req.Completed
	}

	return s.storage.Update(id, existing)
}

func (s *TodoService) DeleteTask(id string) error {
	if err := validateUUID(id); err != nil {
		return err
	}

	return s.storage.Delete(id)
}

func (s *TodoService) CompleteTask(id string) (models.Task, error) {
	if err := validateUUID(id); err != nil {
		return models.Task{}, err
	}

	return s.storage.CompleteTask(id)
}

func validateUUID(id string) error {
	if len(id) != 36 { // UUID v4 длина
		return errors.New("invalid UUID format")
	}
	return nil
}
