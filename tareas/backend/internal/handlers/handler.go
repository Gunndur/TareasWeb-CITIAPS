package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"tareas/internal/models"
	"tareas/internal/services"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskHandler struct {
	service *services.TaskService
}

// Crea un controlador de tareas.
func NewTaskHandler(db *mongo.Database) *TaskHandler {
	return &TaskHandler{service: services.NewTaskService(db)}
}

// Crea una tarea
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	req.Normalize()
	if validationErrors := req.Validate(); len(validationErrors) > 0 {
		respondValidationError(w, validationErrors)
		return
	}

	task := models.Task{
		ID:          primitive.NewObjectID(),
		Title:       req.Title,
		Description: req.Description,
		Completed:   false,
		Tags:        req.Tags,
		CreatedAt:   time.Now().UTC(),
	}

	createdTask, err := h.service.CreateTask(r.Context(), task)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "no se pudo crear la tarea")
		return
	}

	respondJSON(w, http.StatusCreated, createdTask)
}

// Crea múltiples tareas en una sola petición.
func (h *TaskHandler) CreateTasksBulk(w http.ResponseWriter, r *http.Request) {
	var reqs []models.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&reqs); err != nil {
		respondError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	if len(reqs) == 0 {
		respondError(w, http.StatusBadRequest, "debes enviar al menos una tarea")
		return
	}

	if len(reqs) > 100 {
		respondError(w, http.StatusBadRequest, "solo puedes crear hasta 100 tareas por petición")
		return
	}

	tasks := make([]models.Task, 0, len(reqs))
	type itemValidationError struct {
		Index   int               `json:"index"`
		Details map[string]string `json:"details"`
	}
	validationErrors := make([]itemValidationError, 0)

	for i := range reqs {
		req := reqs[i]
		req.Normalize()
		if errs := req.Validate(); len(errs) > 0 {
			validationErrors = append(validationErrors, itemValidationError{
				Index:   i,
				Details: errs,
			})
			continue
		}

		tasks = append(tasks, models.Task{
			ID:          primitive.NewObjectID(),
			Title:       req.Title,
			Description: req.Description,
			Completed:   false,
			Tags:        req.Tags,
			CreatedAt:   time.Now().UTC(),
		})
	}

	if len(validationErrors) > 0 {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"error":   "error de validación",
			"details": validationErrors,
		})
		return
	}

	createdTasks, err := h.service.CreateTasks(r.Context(), tasks)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "no se pudieron crear las tareas")
		return
	}

	respondJSON(w, http.StatusCreated, map[string]any{
		"count": len(createdTasks),
		"items": createdTasks,
	})
}

// Lista tareas paginadas.
func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	page := parsePositiveInt(r.URL.Query().Get("page"), 1)
	limit := parsePositiveInt(r.URL.Query().Get("limit"), 100)

	result, err := h.service.ListTasks(r.Context(), page, limit)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "no se pudo leer las tareas")
		return
	}

	respondJSON(w, http.StatusOK, result)
}

// Obtiene una tarea por ID.
func (h *TaskHandler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseObjectID(mux.Vars(r)["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "id inválido")
		return
	}

	task, err := h.service.GetTaskByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			respondError(w, http.StatusNotFound, "tarea no encontrada")
			return
		}
		respondError(w, http.StatusInternalServerError, "no se pudo obtener la tarea")
		return
	}

	respondJSON(w, http.StatusOK, task)
}

// Marca una tarea como completada
func (h *TaskHandler) CompleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := parseObjectID(mux.Vars(r)["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "id inválido")
		return
	}

	task, err := h.service.CompleteTask(r.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			respondError(w, http.StatusNotFound, "tarea no encontrada")
			return
		}
		respondError(w, http.StatusInternalServerError, "no se pudo actualizar la tarea")
		return
	}

	respondJSON(w, http.StatusOK, task)
}

// Elimina una tarea por ID.
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := parseObjectID(mux.Vars(r)["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "id inválido")
		return
	}

	err = h.service.DeleteTask(r.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			respondError(w, http.StatusNotFound, "tarea no encontrada")
			return
		}
		respondError(w, http.StatusInternalServerError, "no se pudo eliminar la tarea")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "tarea eliminada correctamente"})
}

// Convierte un string a ObjectID
func parseObjectID(id string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(id)
}

// Convierte string a entero positivo
func parsePositiveInt(raw string, fallback int) int {
	if raw == "" {
		return fallback
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return fallback
	}
	return value
}

// Escribe una respuesta JSON
func respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// Escribe una respuesta de error
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]any{
		"error": message,
	})
}

// Escribe una respuesta de validación
func respondValidationError(w http.ResponseWriter, validationErrors map[string]string) {
	respondJSON(w, http.StatusBadRequest, map[string]any{
		"error":   "error de validación",
		"details": validationErrors,
	})
}
