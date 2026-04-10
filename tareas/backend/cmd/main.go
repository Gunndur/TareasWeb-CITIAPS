package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"tareas/internal/config"
	"tareas/internal/routes"
)

func main() {

	// Carga las variables de entorno de la aplicación.
	env := config.LoadEnvConfig()

	// Configura tiempos de espera con la conexión con MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Conecta a MongoDB
	client, err := config.ConnectMongo(ctx, env.MongoURI)
	if err != nil {
		log.Fatalf("error conectando a MongoDB: %v", err)
	}
	// Error al cerrar conexión con MongoDB
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Printf("error cerrando MongoDB: %v", err)
		}
	}()

	// Configura el router con las rutas de la API
	router := routes.NewRouter(client.Database(env.DBName))

	// Tiempos de espera para el servidor HTTP
	server := &http.Server{
		Addr:              ":" + env.Port,
		Handler:           router,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Inicia el servidor HTTP
	log.Printf("backend ejecutándose en http://localhost:%s", env.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("error iniciando servidor: %v", err)
	}
}
