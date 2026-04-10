# Deployment con Docker Compose

Este directorio usa imagenes publicadas en Docker Hub para levantar toda la app sin compilar codigo local.

## Imagenes usadas

- `tgustafsson/tareas-backend:latest`
- `tgustafsson/tareas-frontend:latest`
- `mongo:7`

## Requisitos

- Docker Desktop (o Docker Engine + Docker Compose)

## Levantar servicios

Desde esta carpeta:

```bash
docker compose pull
docker compose up -d
```

## Ver estado

```bash
docker compose ps
docker compose logs -f
```

## URLs

- Frontend: http://localhost:3000
- Backend: http://localhost:8080
- MongoDB: mongodb://localhost:27017

## Detener

```bash
docker compose down
```

Para borrar tambien el volumen de MongoDB:

```bash
docker compose down -v
```
