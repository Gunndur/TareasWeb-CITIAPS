package config

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EnvConfig struct {
	MongoURI string
	DBName   string
	Port     string
}

// Carga las variables de entorno de la aplicación.
func LoadEnvConfig() EnvConfig {
	loadDotEnv("../.env", ".env")

	return EnvConfig{
		MongoURI: getEnv("MONGODB_URI", ""),
		DBName:   getEnv("MONGODB_DB", ""),
		Port:     getEnv("PORT", ""),
	}
}

func loadDotEnv(paths ...string) {
	for _, path := range paths {
		if err := loadDotEnvFile(path); err == nil {
			return
		}
	}
}

func loadDotEnvFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, found := strings.Cut(line, "=")
		if !found {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		value = strings.Trim(value, `"'`)

		if key == "" {
			continue
		}

		if os.Getenv(key) == "" {
			if err := os.Setenv(key, value); err != nil {
				return err
			}
		}
	}

	return scanner.Err()
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// Conecta con MongoDB
func ConnectMongo(ctx context.Context, mongoURI string) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, fmt.Errorf("no se pudo conectar a mongo: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("no se pudo verificar la conexión con mongo: %w", err)
	}

	return client, nil
}
