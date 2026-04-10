
FROM golang:1.26.2-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/tareas ./cmd/main.go

FROM alpine:latest
WORKDIR /app
RUN addgroup -S app && adduser -S app -G app
COPY --from=builder /app/tareas ./tareas

EXPOSE 8080

USER app
CMD ["./tareas"]
