package entity

import "context"

type User struct {
	ID   string
	Name string
}

type IUserRepository interface {
	FindAll(ctx context.Context) ([]User, error)
	FindById(ctx context.Context, id string) (*User, error)
}
