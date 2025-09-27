package database

import (
	"concurrency/internal/entity"
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserEntityMongo struct {
	ID   string `bson:"_id"`
	Name string `bson:"name"`
}

type UserRepository struct {
	Collection mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		Collection: *db.Collection("users"),
	}
}

func (u *UserRepository) FindAll(ctx context.Context) ([]entity.User, error) {
	cur, err := u.Collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var users []entity.User
	for cur.Next(ctx) {
		var d UserEntityMongo
		if err := cur.Decode(&d); err != nil {
			return nil, err
		}
		users = append(users, entity.User{ID: d.ID, Name: d.Name})
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (u *UserRepository) FindById(ctx context.Context, id string) (*entity.User, error) {

	var user UserEntityMongo
	err := u.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &entity.User{ID: user.ID, Name: user.Name}, nil
}
