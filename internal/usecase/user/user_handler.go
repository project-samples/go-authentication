package user

import (
	"context"
	sv "github.com/core-go/core"
	"github.com/core-go/search"
	"net/http"
	"reflect"
)

type UserHandler interface {
	Search(w http.ResponseWriter, r *http.Request)
	Load(w http.ResponseWriter, r *http.Request)
}

func NewUserHandler(find func(context.Context, interface{}, interface{}, int64, ...int64) (int64, string, error), load func(ctx context.Context, id interface{}, result interface{}) (bool, error), logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error) UserHandler {
	searchModelType := reflect.TypeOf(UserFilter{})
	modelType := reflect.TypeOf(User{})
	searchHandler := search.NewSearchHandler(find, modelType, searchModelType, logError, writeLog)
	return &userHandler{load: load, SearchHandler: searchHandler, Error: logError, Log: writeLog}
}

type userHandler struct {
	load func(ctx context.Context, id interface{}, result interface{}) (bool, error)
	*search.SearchHandler
	Error func(context.Context, string, ...map[string]interface{})
	Log   func(context.Context, string, string, bool, string) error
}

func (h *userHandler) Load(w http.ResponseWriter, r *http.Request) {
	id := sv.GetRequiredParam(w, r)
	if len(id) > 0 {
		var user User
		ok, err := h.load(r.Context(), id, &user)
		sv.RespondIfFound(w, r, user, ok, err, h.Error, nil)
	}
}
