package user_usecase

import (
	"concurrency/internal/entity"
	"context"
	"log"
	"time"
)

type FindUserUserInput struct {
	ID string `json:"id"`
}

type FindUserUserOutput struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type IFindUserByIdUseCase interface {
	Execute(ctx context.Context, input FindUserUserInput) (*FindUserUserOutput, error)
}

type FindUserById struct {
	UserRepository entity.IUserRepository
}

func (u *FindUserById) Execute(ctx context.Context, input FindUserUserInput) (*FindUserUserOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	user, err := u.UserRepository.FindById(ctx, input.ID)
	if err != nil {
		log.Println("User not found")
		return nil, err
	}
	return &FindUserUserOutput{
		ID:   user.ID,
		Name: user.Name,
	}, nil
}
