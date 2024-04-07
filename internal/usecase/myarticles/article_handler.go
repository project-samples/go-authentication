package myarticles

import (
	"context"
	sv "github.com/core-go/core"
	"github.com/core-go/core/builder"
	"github.com/core-go/search"
	"net/http"
	"reflect"
)

type ArticleHandler interface {
	Search(w http.ResponseWriter, r *http.Request)
	Load(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Patch(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

func NewArticleHandler(find func(context.Context, interface{}, interface{}, int64, int64) (int64, error), service ArticleService, generateId func(context.Context) (string, error), logError func(context.Context, string, ...map[string]interface{}), validate func(ctx context.Context, model interface{}) ([]sv.ErrorMessage, error), tracking builder.TrackingConfig, action *sv.ActionConfig, writeLog func(context.Context, string, string, bool, string) error) ArticleHandler {
	searchModelType := reflect.TypeOf(ArticleFilter{})
	modelType := reflect.TypeOf(Article{})
	builder := builder.NewBuilderWithIdAndConfig(generateId, modelType, tracking)
	patchHandler, params := sv.CreatePatchAndParams(modelType, logError, service.Patch, validate, builder.Patch, action, writeLog)
	searchHandler := search.NewSearchHandler(find, modelType, searchModelType, logError, params.Log)
	return &articleHandler{service: service, builder: builder, PatchHandler: patchHandler, SearchHandler: searchHandler, Params: params}
}

type articleHandler struct {
	service ArticleService
	builder sv.Builder
	*sv.PatchHandler
	*search.SearchHandler
	*sv.Params
}

func (h *articleHandler) Load(w http.ResponseWriter, r *http.Request) {
	id := sv.GetRequiredParam(w, r)
	if len(id) > 0 {
		result, err := h.service.Load(r.Context(), id)
		sv.Return(w, r, result, err, h.Error, nil)
	}
}
func (h *articleHandler) Create(w http.ResponseWriter, r *http.Request) {
	var article Article
	er1 := sv.Decode(w, r, &article, h.builder.Create)
	if er1 == nil {
		errors, er2 := h.Validate(r.Context(), &article)
		if !sv.HasError(w, r, errors, er2, h.Error, h.Log, h.Resource, h.Action.Create) {
			result, er3 := h.service.Create(r.Context(), &article)
			sv.AfterCreated(w, r, &article, result, er3, h.Error, h.Log, h.Resource, h.Action.Create)
		}
	}
}
func (h *articleHandler) Update(w http.ResponseWriter, r *http.Request) {
	var article Article
	er1 := sv.DecodeAndCheckId(w, r, &article, h.Keys, h.Indexes, h.builder.Update)
	if er1 == nil {
		errors, er2 := h.Validate(r.Context(), &article)
		if !sv.HasError(w, r, errors, er2, h.Error, h.Log, h.Resource, h.Action.Update) {
			result, er3 := h.service.Update(r.Context(), &article)
			sv.HandleResult(w, r, &article, result, er3, h.Error, h.Log, h.Resource, h.Action.Update)
		}
	}
}
func (h *articleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := sv.GetRequiredParam(w, r)
	if len(id) > 0 {
		result, err := h.service.Delete(r.Context(), id)
		sv.HandleDelete(w, r, result, err, h.Error, h.Log, h.Resource, h.Action.Delete)
	}
}
