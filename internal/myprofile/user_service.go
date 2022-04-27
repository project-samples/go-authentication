package myprofile

import (
	"context"
	sv "github.com/core-go/service"
)

type UserService interface {
	Load(ctx context.Context, id string) (*User, error)
	Create(ctx context.Context, user *User) (int64, error)
	Patch(ctx context.Context, user map[string]interface{}) (int64, error)
	Update(ctx context.Context, user *User) (int64, error)
	Delete(ctx context.Context, id string) (int64, error)
}

func NewUserService(repository sv.Repository) UserService {
	return &userService{repository: repository}
}

type userService struct {
	repository sv.Repository
}

func (s *userService) Patch(ctx context.Context, user map[string]interface{}) (int64, error) {
	return s.repository.Patch(ctx, user)
}

func (s *userService) Update(ctx context.Context, user *User) (int64, error) {
	return s.repository.Update(ctx, user)
}

func (s *userService) Delete(ctx context.Context, id string) (int64, error) {
	return s.repository.Delete(ctx, id)
}

func (s *userService) Load(ctx context.Context, id string) (*User, error) {
	var user User
	ok, err := s.repository.LoadAndDecode(ctx, id, &user)
	if !ok {
		return nil, err
	} else {
		return &user, err
	}
}

func (s *userService) Create(ctx context.Context, user *User) (int64, error) {
	return s.repository.Insert(ctx, user)
}
