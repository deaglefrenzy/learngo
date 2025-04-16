package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepository struct {
	Client     *mongo.Client
	Database   *mongo.Database
	Collection *mongo.Collection
}

func NewMongoRepository(client *mongo.Client, dbName, collectionName string) *MongoRepository {
	db := client.Database(dbName)
	return &MongoRepository{
		Client:     client,
		Database:   db,
		Collection: db.Collection(collectionName),
	}
}

func (r *MongoRepository) WithCollection(collectionName string) *MongoRepository {
	return &MongoRepository{
		Client:     r.Client,
		Database:   r.Database,
		Collection: r.Database.Collection(collectionName),
	}
}

func (r *MongoRepository) Create(document interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.Collection.InsertOne(ctx, document)
	return err
}

func (r *MongoRepository) FindOne(filter interface{}, dest interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.Collection.FindOne(ctx, filter).Decode(dest)
}

func (r *MongoRepository) FindMany(filter interface{}, dest interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	return cursor.All(ctx, dest)
}

func (r *MongoRepository) UpdateOne(filter interface{}, update interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.Collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *MongoRepository) CountDocuments(filter interface{}) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := r.Collection.CountDocuments(ctx, filter)
	return count, err
}

func (r *MongoRepository) DeleteOne(filter interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.Collection.DeleteOne(ctx, filter)
	return err
}
