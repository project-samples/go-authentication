package user

import (
	"context"
	"net/http"

	"github.com/core-go/search"
	sv "github.com/core-go/service"
)

type UserHandler interface {
	Search(w http.ResponseWriter, r *http.Request)
	Load(w http.ResponseWriter, r *http.Request)
}

type userHandler struct {
	load func(ctx context.Context, id interface{}, result interface{}) (bool, error)
	*search.SearchHandler
	Error func(context.Context, string)
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
