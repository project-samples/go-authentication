package user

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	sv "github.com/core-go/core"
	s "github.com/core-go/search"
)

type UserTransport interface {
	Search(w http.ResponseWriter, r *http.Request)
	Load(w http.ResponseWriter, r *http.Request)
}

func NewUserHandler(query UserQuery, logError func(context.Context, string, ...map[string]interface{})) UserTransport {
	paramIndex, filterIndex := s.BuildParams(reflect.TypeOf(UserFilter{}))
	return &UserHandler{Query: query, LogError: logError, paramIndex: paramIndex, filterIndex: filterIndex}
}

type UserHandler struct {
	Query       UserQuery
	LogError    func(context.Context, string, ...map[string]interface{})
	paramIndex  map[string]int
	filterIndex int
}

func (h *UserHandler) Load(w http.ResponseWriter, r *http.Request) {
	id := sv.GetRequiredParam(w, r)
	if len(id) > 0 {
		user, err := h.Query.Load(r.Context(), id)
		if err != nil {
			h.LogError(r.Context(), fmt.Sprintf("Error to get user %s: %s", id, err.Error()))
			http.Error(w, sv.InternalServerError, http.StatusInternalServerError)
			return
		}
		if user == nil {
			sv.JSON(w, http.StatusNotFound, user)
		} else {
			sv.JSON(w, http.StatusOK, user)
		}
	}
}
func (h *UserHandler) Search(w http.ResponseWriter, r *http.Request) {
	filter := UserFilter{Filter: &s.Filter{}}
	s.Decode(r, &filter, h.paramIndex, h.filterIndex)

	offset := s.GetOffset(filter.Limit, filter.Page)
	users, total, err := h.Query.Search(r.Context(), &filter, filter.Limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sv.JSON(w, http.StatusOK, &s.Result{List: &users, Total: total})
}
