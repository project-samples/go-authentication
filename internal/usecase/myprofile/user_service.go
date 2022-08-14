package myprofile

import (
	"context"
	sv "github.com/core-go/core"
)

type UserService interface {
	GetMyProfile(ctx context.Context, id string) (*User, error)
	SaveMyProfile(ctx context.Context, user map[string]interface{}) (int64, error)
	GetMySettings(ctx context.Context, id string) (*Settings, error)
	SaveMySettings(ctx context.Context, id string, settings *Settings) (int64, error)
}

func NewUserService(repository sv.Repository) UserService {
	return &userService{repository: repository}
}

type userService struct {
	repository sv.Repository
}

func (s *userService) SaveMyProfile(ctx context.Context, user map[string]interface{}) (int64, error) {
	return s.repository.Patch(ctx, user)
}
func (s *userService) GetMyProfile(ctx context.Context, id string) (*User, error) {
	var user User
	ok, err := s.repository.LoadAndDecode(ctx, id, &user)
	if !ok {
		return nil, err
	} else {
		user.Settings = nil
		return &user, err
	}
}
func (s *userService) GetMySettings(ctx context.Context, id string) (*Settings, error) {
	var user User
	ok, err := s.repository.LoadAndDecode(ctx, id, &user)
	if !ok {
		return nil, err
	} else {
		return user.Settings, nil
	}
}
func (s *userService) SaveMySettings(ctx context.Context, id string, settings *Settings) (int64, error) {
	user := make(map[string]interface{})
	user["userId"] = id
	user["settings"] = settings
	return s.repository.Patch(ctx, user)
}
