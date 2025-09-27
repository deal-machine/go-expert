package mongodb

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	DATABASE_NAME     = "DATABASE_NAME"
	DATABASE_HOST     = "DATABASE_HOST"
	DATABASE_USER     = "DATABASE_USER"
	DATABASE_PASSWORD = "DATABASE_PASSWORD"
	DATABASE_PORT     = "DATABASE_PORT"
)

func NewMongoDBConnection(ctx context.Context) (*mongo.Database, error) {
	dbUsername := os.Getenv(DATABASE_USER)
	dbPassword := os.Getenv(DATABASE_PASSWORD)
	dbHost := os.Getenv(DATABASE_HOST)
	dbPort := os.Getenv(DATABASE_PORT)
	dbName := os.Getenv(DATABASE_NAME)

	url := "mongodb://" + dbUsername + ":" + dbPassword + "@" + dbHost + ":" + dbPort + "/"

	client, err := mongo.Connect(options.Client().ApplyURI(url))
	if err != nil {
		log.Println("Error on connecting with mongodb")
		return nil, err
	}
	if err = client.Ping(ctx, nil); err != nil {
		log.Println("Error on connecting with mongodb")
		return nil, err
	}

	return client.Database(dbName), nil
}
