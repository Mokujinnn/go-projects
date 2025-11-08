package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"todo-api/internal/models"
	"todo-api/internal/service"
)

type TodoHandler struct {
	service *service.TodoService
}

func NewTodoHandler(service *service.TodoService) *TodoHandler {
	return &TodoHandler{
		service: service,
	}
}

// CreateTask создает новую задачу
// @Summary Создать новую задачу
// @Description Создает новую задачу с указанным заголовком и описанием
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body models.CreateTaskRequest true "Данные для создания задачи"
// @Success 201 {object} models.Task
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks [post]
func (h *TodoHandler) CreateTask(c *gin.Context) {
	var req models.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.service.CreateTask(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// GetTasks возвращает список задач с пагинацией и фильтрацией
// @Summary Получить список задач
// @Description Возвращает список задач с поддержкой пагинации, фильтрации и поиска
// @Tags tasks
// @Accept json
// @Produce json
// @Param limit query int false "Лимит (по умолчанию 10)"
// @Param offset query int false "Смещение (по умолчанию 0)"
// @Param completed query bool false "Фильтр по статусу выполнения"
// @Param search query string false "Поиск по заголовку и описанию"
// @Param sort_by query string false "Поле для сортировки (created_at, completed)"
// @Param sort_order query string false "Порядок сортировки (asc, desc)"
// @Success 200 {object} models.TasksResponse
// @Failure 500 {object} map[string]string
// @Router /tasks [get]
func (h *TodoHandler) GetTasks(c *gin.Context) {
	query := models.TaskQuery{}

	// Параметры пагинации
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			query.Limit = limit
		}
	} else {
		query.Limit = 10
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			query.Offset = offset
		}
	} else {
		query.Offset = 0
	}

	// Фильтр по статусу выполнения
	if completedStr := c.Query("completed"); completedStr != "" {
		if completed, err := strconv.ParseBool(completedStr); err == nil {
			query.Completed = &completed
		}
	}

	// Поиск
	if search := c.Query("search"); search != "" {
		query.Search = search
	}

	// Сортировка
	if sortBy := c.Query("sort_by"); sortBy != "" {
		query.SortBy = sortBy
	}
	if sortOrder := c.Query("sort_order"); sortOrder != "" {
		query.SortOrder = sortOrder
	}

	response, err := h.service.GetAllTasks(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tasks"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetTask возвращает задачу по ID
// @Summary Получить задачу по ID
// @Description Возвращает задачу по её уникальному идентификатору
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "ID задачи"
// @Success 200 {object} models.Task
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/{id} [get]
func (h *TodoHandler) GetTask(c *gin.Context) {
	id := c.Param("id")

	task, err := h.service.GetTask(id)
	if err != nil {
		if err.Error() == "task not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, task)
}

// UpdateTask обновляет задачу
// @Summary Обновить задачу
// @Description Обновляет данные задачи по её ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "ID задачи"
// @Param task body models.UpdateTaskRequest true "Данные для обновления"
// @Success 200 {object} models.Task
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/{id} [put]
func (h *TodoHandler) UpdateTask(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.service.UpdateTask(id, req)
	if err != nil {
		if err.Error() == "task not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, task)
}

// DeleteTask удаляет задачу
// @Summary Удалить задачу
// @Description Удаляет задачу по её ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "ID задачи"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/{id} [delete]
func (h *TodoHandler) DeleteTask(c *gin.Context) {
	id := c.Param("id")

	err := h.service.DeleteTask(id)
	if err != nil {
		if err.Error() == "task not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// CompleteTask отмечает задачу как выполненную
// @Summary Отметить задачу как выполненную
// @Description Устанавливает статус выполнения задачи в true
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "ID задачи"
// @Success 200 {object} models.Task
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/{id}/complete [patch]
func (h *TodoHandler) CompleteTask(c *gin.Context) {
	id := c.Param("id")

	task, err := h.service.CompleteTask(id)
	if err != nil {
		if err.Error() == "task not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, task)
}
