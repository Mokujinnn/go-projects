package storage

import (
	"errors"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"todo-api/internal/models"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

type MemoryStorage struct {
	mu    sync.RWMutex
	tasks map[string]models.Task
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		tasks: make(map[string]models.Task),
	}
}

func (s *MemoryStorage) Create(task models.Task) (models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task.ID = uuid.New().String()
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	s.tasks[task.ID] = task
	return task, nil
}

func (s *MemoryStorage) GetByID(id string) (models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, exists := s.tasks[id]
	if !exists {
		return models.Task{}, ErrTaskNotFound
	}

	return task, nil
}

func (s *MemoryStorage) GetAll(query models.TaskQuery) ([]models.Task, int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var tasks []models.Task
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}

	// Фильтрация
	tasks = s.filterTasks(tasks, query)

	// Сортировка
	tasks = s.sortTasks(tasks, query)

	// Пагинация
	total := len(tasks)
	start := query.Offset
	if start > len(tasks) {
		start = len(tasks)
	}
	end := start + query.Limit
	if end > len(tasks) {
		end = len(tasks)
	}

	return tasks[start:end], total, nil
}

func (s *MemoryStorage) filterTasks(tasks []models.Task, query models.TaskQuery) []models.Task {
	var filtered []models.Task

	for _, task := range tasks {
		// Фильтр по статусу выполнения
		if query.Completed != nil && task.Completed != *query.Completed {
			continue
		}

		// Полнотекстовый поиск
		if query.Search != "" {
			searchLower := strings.ToLower(query.Search)
			titleLower := strings.ToLower(task.Title)
			descLower := strings.ToLower(task.Description)

			if !strings.Contains(titleLower, searchLower) &&
				!strings.Contains(descLower, searchLower) {
				continue
			}
		}

		filtered = append(filtered, task)
	}

	return filtered
}

func (s *MemoryStorage) sortTasks(tasks []models.Task, query models.TaskQuery) []models.Task {
	sort.Slice(tasks, func(i, j int) bool {
		switch query.SortBy {
		case "created_at":
			if query.SortOrder == "desc" {
				return tasks[i].CreatedAt.After(tasks[j].CreatedAt)
			}
			return tasks[i].CreatedAt.Before(tasks[j].CreatedAt)
		case "completed":
			if query.SortOrder == "desc" {
				return tasks[i].Completed && !tasks[j].Completed
			}
			return !tasks[i].Completed && tasks[j].Completed
		default:
			return tasks[i].CreatedAt.After(tasks[j].CreatedAt)
		}
	})

	return tasks
}

func (s *MemoryStorage) Update(id string, updatedTask models.Task) (models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, exists := s.tasks[id]
	if !exists {
		return models.Task{}, ErrTaskNotFound
	}

	updatedTask.ID = id
	updatedTask.CreatedAt = existing.CreatedAt
	updatedTask.UpdatedAt = time.Now()

	s.tasks[id] = updatedTask
	return updatedTask, nil
}

func (s *MemoryStorage) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[id]; !exists {
		return ErrTaskNotFound
	}

	delete(s.tasks, id)
	return nil
}

func (s *MemoryStorage) CompleteTask(id string) (models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[id]
	if !exists {
		return models.Task{}, ErrTaskNotFound
	}

	task.Completed = true
	task.UpdatedAt = time.Now()
	s.tasks[id] = task

	return task, nil
}
