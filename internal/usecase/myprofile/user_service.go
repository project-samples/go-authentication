package myprofile

import "context"

type UserService interface {
	GetMyProfile(ctx context.Context, id string) (*User, error)
	SaveMyProfile(ctx context.Context, user map[string]interface{}) (int64, error)
	GetMySettings(ctx context.Context, id string) (*Settings, error)
	SaveMySettings(ctx context.Context, id string, settings *Settings) (int64, error)
}
type Repository interface {
	Get(ctx context.Context, id string, result interface{}) (bool, error)
	Exist(ctx context.Context, id string) (bool, error)
	Insert(ctx context.Context, model interface{}) (int64, error)
	Update(ctx context.Context, model interface{}) (int64, error)
	Patch(ctx context.Context, model map[string]interface{}) (int64, error)
	Delete(ctx context.Context, id string) (int64, error)
}
func NewUserService(repository Repository) UserService {
	return &userService{repository: repository}
}

type userService struct {
	repository Repository
}

func (s *userService) SaveMyProfile(ctx context.Context, user map[string]interface{}) (int64, error) {
	return s.repository.Patch(ctx, user)
}
func (s *userService) GetMyProfile(ctx context.Context, id string) (*User, error) {
	var user User
	ok, err := s.repository.Get(ctx, id, &user)
	if !ok {
		return nil, err
	} else {
		user.Settings = nil
		return &user, err
	}
}
func (s *userService) GetMySettings(ctx context.Context, id string) (*Settings, error) {
	var user User
	ok, err := s.repository.Get(ctx, id, &user)
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
