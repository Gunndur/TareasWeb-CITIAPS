package repository

import (
	"context"
	"time"

	"tareas/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TaskRepository struct {
	collection *mongo.Collection
}

// Crea un repositorio de tareas.
func NewTaskRepository(db *mongo.Database) *TaskRepository {
	return &TaskRepository{collection: db.Collection("tasks")}
}

// Inserta una tarea.
func (r *TaskRepository) CreateTask(parent context.Context, task models.Task) (models.Task, error) {
	ctx, cancel := context.WithTimeout(parent, 5*time.Second)
	defer cancel()

	_, err := r.collection.InsertOne(ctx, task)
	if err != nil {
		return models.Task{}, err
	}

	return task, nil
}

// Inserta múltiples tareas.
func (r *TaskRepository) CreateTasks(parent context.Context, tasks []models.Task) ([]models.Task, error) {
	ctx, cancel := context.WithTimeout(parent, 5*time.Second)
	defer cancel()

	docs := make([]any, len(tasks))
	for i := range tasks {
		docs[i] = tasks[i]
	}

	_, err := r.collection.InsertMany(ctx, docs)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

// Obtiene tareas paginadas.
func (r *TaskRepository) ListTasks(parent context.Context, page, limit int) ([]models.Task, int64, error) {
	skip := int64((page - 1) * limit)

	ctx, cancel := context.WithTimeout(parent, 5*time.Second)
	defer cancel()

	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	findOptions := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}).SetSkip(skip).SetLimit(int64(limit))
	cursor, err := r.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var tasks []models.Task
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

// Busca una tarea por ID.
func (r *TaskRepository) GetTaskByID(parent context.Context, id primitive.ObjectID) (models.Task, error) {
	return r.findTaskByID(parent, id)
}

// Marca una tarea como completada.
func (r *TaskRepository) CompleteTask(parent context.Context, id primitive.ObjectID) (models.Task, error) {
	ctx, cancel := context.WithTimeout(parent, 5*time.Second)
	defer cancel()

	result, err := r.collection.UpdateByID(ctx, id, bson.M{"$set": bson.M{"completed": true}})
	if err != nil {
		return models.Task{}, err
	}
	if result.MatchedCount == 0 {
		return models.Task{}, mongo.ErrNoDocuments
	}

	return r.findTaskByID(parent, id)
}

// Elimina una tarea por ID.
func (r *TaskRepository) DeleteTask(parent context.Context, id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(parent, 5*time.Second)
	defer cancel()

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

// Recupera una tarea por ID desde MongoDB.
func (r *TaskRepository) findTaskByID(parent context.Context, id primitive.ObjectID) (models.Task, error) {
	ctx, cancel := context.WithTimeout(parent, 5*time.Second)
	defer cancel()

	var task models.Task
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&task)
	return task, err
}
