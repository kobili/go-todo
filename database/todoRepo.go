package database

import (
	"context"
	"go-todo/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TodoRepository struct {
	todoCollection *mongo.Collection
}

func NewTodoRepository(collection *mongo.Collection) *TodoRepository {
	return &TodoRepository{
		todoCollection: collection,
	}
}

func (r *TodoRepository) AddOne(ctx context.Context, document models.Todo) (*mongo.InsertOneResult, error) {
	result, err := r.todoCollection.InsertOne(ctx, document)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *TodoRepository) FindById(ctx context.Context, id string) (primitive.M, error) {
	var result bson.M

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = r.todoCollection.FindOne(
		ctx,
		bson.D{{Key: "_id", Value: objId}},
	).Decode(&result)

	if err != nil {
		return nil, err
	}

	return result, nil
}
