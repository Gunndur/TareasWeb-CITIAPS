package routes

import (
	"net/http"
	"tareas/internal/handlers"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewRouter(db *mongo.Database) http.Handler {
	taskHandler := handlers.NewTaskHandler(db)

	r := mux.NewRouter()

	// Middleware CORS
	r.Use(corsMiddleware)
	r.Use(jsonMiddleware)

	// Ruta de health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}).Methods(http.MethodGet)

	// Rutas de tareas
	r.PathPrefix("/").Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	r.HandleFunc("/tasks", taskHandler.CreateTask).Methods(http.MethodPost)
	r.HandleFunc("/tasks/bulk", taskHandler.CreateTasksBulk).Methods(http.MethodPost)
	r.HandleFunc("/tasks", taskHandler.ListTasks).Methods(http.MethodGet)
	r.HandleFunc("/tasks/{id}", taskHandler.GetTaskByID).Methods(http.MethodGet)
	r.HandleFunc("/tasks/{id}/complete", taskHandler.CompleteTask).Methods(http.MethodPut)
	r.HandleFunc("/tasks/{id}", taskHandler.DeleteTask).Methods(http.MethodDelete)

	return r
}

// CORS Middleware mejorado
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Para desarrollo: permite todo
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept")
		w.Header().Set("Access-Control-Max-Age", "86400") // 24 horas

		// Manejo del preflight (OPTIONS)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
