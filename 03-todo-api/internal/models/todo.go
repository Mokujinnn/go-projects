package models

import (
	"time"
)

type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description,omitempty"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required,max=200"`
	Description string `json:"description,omitempty"`
}

type UpdateTaskRequest struct {
	Title       string `json:"title" binding:"max=200"`
	Description string `json:"description,omitempty"`
	Completed   *bool  `json:"completed,omitempty"`
}

type TasksResponse struct {
	Tasks  []Task `json:"tasks"`
	Total  int    `json:"total"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

type TaskQuery struct {
	Limit     int
	Offset    int
	Completed *bool
	Search    string
	SortBy    string
	SortOrder string
}
