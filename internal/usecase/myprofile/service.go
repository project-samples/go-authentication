package myprofile

import "context"

type UserService interface {
	GetMyProfile(ctx context.Context, id string) (*User, error)
	SaveMyProfile(ctx context.Context, user map[string]interface{}) (int64, error)
	GetMySettings(ctx context.Context, id string) (*Settings, error)
	SaveMySettings(ctx context.Context, id string, settings *Settings) (int64, error)
}

func NewUserService(repository UserRepository) *UserUseCase {
	return &UserUseCase{repository: repository}
}

type UserUseCase struct {
	repository UserRepository
}

func (s *UserUseCase) SaveMyProfile(ctx context.Context, user map[string]interface{}) (int64, error) {
	return s.repository.Patch(ctx, user)
}
func (s *UserUseCase) GetMyProfile(ctx context.Context, id string) (*User, error) {
	user, err := s.repository.Load(ctx, id)
	if err != nil {
		return nil, err
	}
	if user != nil {
		user.Settings = nil
	}
	return user, nil
}
func (s *UserUseCase) GetMySettings(ctx context.Context, id string) (*Settings, error) {
	user, err := s.repository.Load(ctx, id)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return user.Settings, nil
	}
	return nil, nil
}
func (s *UserUseCase) SaveMySettings(ctx context.Context, id string, settings *Settings) (int64, error) {
	user := make(map[string]interface{})
	user["userId"] = id
	user["settings"] = settings
	return s.repository.Patch(ctx, user)
}
