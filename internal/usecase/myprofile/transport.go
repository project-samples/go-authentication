package myprofile

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"

	repo "github.com/core-go/mongo/repository"
)

type MyProfileTransport interface {
	GetMyProfile(w http.ResponseWriter, r *http.Request)
	SaveMyProfile(w http.ResponseWriter, r *http.Request)
	GetMySettings(w http.ResponseWriter, r *http.Request)
	SaveMySettings(w http.ResponseWriter, r *http.Request)
}

func NewMyProfileTransport(db *mongo.Database, logError func(context.Context, string, ...map[string]interface{}),
	saveSkills func(ctx context.Context, values []string) (int64, error),
	saveInterests func(ctx context.Context, values []string) (int64, error),
	saveLookingFor func(ctx context.Context, values []string) (int64, error)) (*MyProfileHandler, error) {

	repository := repo.NewRepository[User, string](db, "user")
	service := NewUserService(repository)
	return NewMyProfileHandler(service, logError, saveSkills, saveInterests, saveLookingFor)
}
