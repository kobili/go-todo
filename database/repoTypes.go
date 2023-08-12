package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"go-todo/models"
)

type TodoRepo interface {
	AddOne(ctx context.Context, document models.Todo) (*mongo.InsertOneResult, error)
	FindById(ctx context.Context, id string) (primitive.M, error)
}
