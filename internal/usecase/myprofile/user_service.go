package myprofile

import (
	"context"

	sv "github.com/core-go/service"
)

type UserService interface {
	GetMyProfile(ctx context.Context, id string) (*User, error)
	LoadData(ctx context.Context, id string) (interface{}, error)
	SaveMyProfile(ctx context.Context, user map[string]interface{}) (int64, error)
	GetMySettings(ctx context.Context, id string) (*Settings, error)
	SaveMySettings(ctx context.Context, id string, settings *Settings) (int64, error)
	insertImage(ctx context.Context, id string, image string) (int64, error)
	insertCover(ctx context.Context, id string, image string) (int64, error)
	insertGallery(ctx context.Context, id string, image UploadInfo) (int64, error)
	deleteGallery(ctx context.Context, id string, url string) (int64, error)
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
func (s *userService) insertImage(ctx context.Context, id string, image string) (int64, error) {
	user := make(map[string]interface{})
	user["userId"] = id
	user["imageURL"] = image
	return s.repository.Patch(ctx, user)
}

func (s *userService) insertCover(ctx context.Context, id string, image string) (int64, error) {
	user := make(map[string]interface{})
	user["userId"] = id
	user["coverURL"] = image
	return s.repository.Patch(ctx, user)
}

func (s *userService) insertGallery(ctx context.Context, id string, image UploadInfo) (int64, error) {
	var user User
	ok, err := s.repository.LoadAndDecode(ctx, id, &user)
	if !ok {
		return 0, err
	}
	gallery := user.Gallery
	gallery = append(gallery, image)
	user2 := make(map[string]interface{})
	user2["userId"] = id
	user2["gallery"] = gallery
	return s.repository.Patch(ctx, user2)

}
func (s *userService) deleteGallery(ctx context.Context, id string, url string) (int64, error) {
	var user User
	ok, err := s.repository.LoadAndDecode(ctx, id, &user)
	if !ok {
		return 0, err
	}
	gallery := user.Gallery
	idx := 0
	for i := range gallery {
		if gallery[i].Url == url {
			idx = i
		}
	}
	gallery = append(gallery[:idx], gallery[idx+1:]...)
	user2 := make(map[string]interface{})
	user2["userId"] = id
	user2["gallery"] = gallery
	return s.repository.Patch(ctx, user2)

}
func (s *userService) LoadData(ctx context.Context, id string) (interface{}, error) {
	var user User
	ok, err := s.repository.LoadAndDecode(ctx, id, &user)
	if !ok {
		return nil, err
	} else {
		user.Settings = nil
		return &user, err
	}
}
