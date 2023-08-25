package database

import (
	"context"
	"go-todo/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TodoRepo interface {
	AddOne(ctx context.Context, document models.Todo) (*mongo.InsertOneResult, error)
	FindById(ctx context.Context, id string) (*models.Todo, error)
	UpdateById(ctx context.Context, id string, requestBody struct{ Text string }) (*mongo.UpdateResult, error)
	ReplaceById(ctx context.Context, id string, todo models.Todo) (*mongo.UpdateResult, error)
}

type TodoRepository struct {
	todoCollection *mongo.Collection
}

func NewTodoRepository(mongoClient *mongo.Client) *TodoRepository {
	collection := mongoClient.Database("go-todo").Collection("todos")
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

func (r *TodoRepository) FindById(ctx context.Context, id string) (*models.Todo, error) {
	var result models.Todo

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

	return &result, nil
}

func (r *TodoRepository) UpdateById(ctx context.Context, id string, requestBody struct{ Text string }) (*mongo.UpdateResult, error) {
	updateFields := bson.M{
		"text": requestBody.Text,
	}

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	result, err := r.todoCollection.UpdateOne(
		ctx,
		bson.M{"_id": objId},
		bson.M{"$set": updateFields},
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *TodoRepository) ReplaceById(ctx context.Context, id string, todo models.Todo) (*mongo.UpdateResult, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	result, err := r.todoCollection.ReplaceOne(
		ctx,
		bson.M{"_id": objId},
		todo,
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}
