package models

import (
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Completed   bool               `bson:"completed" json:"completed"`
	Tags        []string           `bson:"tags" json:"tags"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
}

type CreateTaskRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

// Normaliza los campos de la solicitud de creación de tarea, eliminando espacios y etiquetas vacías
func (r *CreateTaskRequest) Normalize() {
	r.Title = strings.TrimSpace(r.Title)
	r.Description = strings.TrimSpace(r.Description)
	cleanTags := make([]string, 0, len(r.Tags))
	seen := map[string]bool{}
	for _, tag := range r.Tags {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			continue
		}
		if !seen[tag] {
			cleanTags = append(cleanTags, tag)
			seen[tag] = true
		}
	}
	r.Tags = cleanTags
}

// Valida los campos de la solicitud de creación de tarea y devuelve errores si hay problemas
func (r CreateTaskRequest) Validate() map[string]string {
	errors := map[string]string{}
	if r.Title == "" {
		errors["title"] = "el título es obligatorio"
	}
	if len(r.Title) > 120 {
		errors["title"] = "el título no puede superar 120 caracteres"
	}
	if len(r.Description) > 500 {
		errors["description"] = "la descripción no puede superar 500 caracteres"
	}
	return errors
}
