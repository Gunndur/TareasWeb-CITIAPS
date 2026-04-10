# Tareas Backend

API REST en Go para gestionar tareas, con MongoDB.

## Requisitos

- Go 1.26+
- MongoDB corriendo local o remoto

## Variables de entorno

Archivo `.env` en la raíz:
(Entiendo que no hay que mostrarlo, pero ya que es de ejemplo...)
```env
PORT=8080
MONGODB_URI=mongodb://localhost:27017
MONGODB_DB=taskmanager
```

## Ejecutar

Desde la raíz del proyecto:
```bash
go mod tidy
go run ./cmd/main.go
```

La API queda en: `http://localhost:8080`

## Endpoints

- `GET /health`
- `POST /tasks`
- `POST /tasks/bulk`
- `GET /tasks`
- `GET /tasks/{id}`
- `PUT /tasks/{id}/complete`
- `DELETE /tasks/{id}`

`GET /tasks` acepta `page` y `limit` como query params opcionales.
Si no envías `limit`, por defecto devuelve hasta 100 tareas.

### Crear varias tareas

`POST /tasks/bulk`

Body (array JSON):

```json
[
	{
		"title": "Tarea 1",
		"description": "Primera tarea",
		"tags": ["trabajo", "urgente"]
	},
	{
		"title": "Tarea 2",
		"description": "Segunda tarea",
		"tags": ["personal"]
	}
]
```
