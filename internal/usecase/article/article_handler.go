package article

import (
	"context"
	"net/http"
	"reflect"

	"github.com/core-go/search"
	sv "github.com/core-go/service"
)

type ArticleHandler interface {
	Search(w http.ResponseWriter, r *http.Request)
	Load(w http.ResponseWriter, r *http.Request)
}

func NewArticleHandler(find func(context.Context, interface{}, interface{}, int64, ...int64) (int64, string, error), service ArticleService, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error) ArticleHandler {
	searchModelType := reflect.TypeOf(ArticleFilter{})
	modelType := reflect.TypeOf(Article{})
	searchHandler := search.NewSearchHandler(find, modelType, searchModelType, logError, writeLog)
	return &articleHandler{service: service, SearchHandler: searchHandler, Error: logError, Log: writeLog}
}

type articleHandler struct {
	service ArticleService
	*search.SearchHandler
	Error func(context.Context, string)
	Log   func(context.Context, string, string, bool, string) error
}

func (h *articleHandler) Load(w http.ResponseWriter, r *http.Request) {
	id := sv.GetRequiredParam(w, r)
	if len(id) > 0 {
		result, err := h.service.Load(r.Context(), id)
		sv.RespondModel(w, r, result, err, h.Error, nil)
	}
}
