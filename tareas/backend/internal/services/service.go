package services

import (
	"context"

	"tareas/internal/models"
	"tareas/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskService struct {
	repo *repository.TaskRepository
}

type TaskListResult struct {
	Items      []models.Task  `json:"items"`
	Pagination PaginationMeta `json:"pagination"`
}

type PaginationMeta struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

// Crea un servicio de tareas.
func NewTaskService(db *mongo.Database) *TaskService {
	return &TaskService{repo: repository.NewTaskRepository(db)}
}

// Normaliza parámetros de paginación.
func normalizePagination(page, limit int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 100
	}
	if limit > 100 {
		limit = 100
	}
	return page, limit
}

// Inserta una tarea.
func (s *TaskService) CreateTask(parent context.Context, task models.Task) (models.Task, error) {
	return s.repo.CreateTask(parent, task)
}

// Inserta múltiples tareas.
func (s *TaskService) CreateTasks(parent context.Context, tasks []models.Task) ([]models.Task, error) {
	return s.repo.CreateTasks(parent, tasks)
}

// Obtiene tareas paginadas y metadatos.
func (s *TaskService) ListTasks(parent context.Context, page, limit int) (TaskListResult, error) {
	page, limit = normalizePagination(page, limit)
	tasks, total, err := s.repo.ListTasks(parent, page, limit)
	if err != nil {
		return TaskListResult{}, err
	}

	return TaskListResult{
		Items: tasks,
		Pagination: PaginationMeta{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	}, nil
}

// Busca una tarea por ID
func (s *TaskService) GetTaskByID(parent context.Context, id primitive.ObjectID) (models.Task, error) {
	return s.repo.GetTaskByID(parent, id)
}

// Marca una tarea como completada
func (s *TaskService) CompleteTask(parent context.Context, id primitive.ObjectID) (models.Task, error) {
	return s.repo.CompleteTask(parent, id)
}

// Elimina una tarea por ID
func (s *TaskService) DeleteTask(parent context.Context, id primitive.ObjectID) error {
	return s.repo.DeleteTask(parent, id)
}
