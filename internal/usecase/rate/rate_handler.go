package rate

import (
	"context"
	"net/http"
	"reflect"

	sv "github.com/core-go/core"
	"github.com/core-go/search"
)

type RateHandler interface {
	Search(w http.ResponseWriter, r *http.Request)
	Load(w http.ResponseWriter, r *http.Request)
}

func NewRateHandler(find func(context.Context, interface{}, interface{}, int64, int64) (int64, error), load func(ctx context.Context, id interface{}, result interface{}) (bool, error), logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error) RateHandler {
	searchModelType := reflect.TypeOf(RateFilter{})
	modelType := reflect.TypeOf(Rate{})
	searchHandler := search.NewSearchHandler(find, modelType, searchModelType, logError, writeLog)
	return &rateHandler{load: load, SearchHandler: searchHandler, Error: logError, Log: writeLog}
}

type rateHandler struct {
	load func(ctx context.Context, id interface{}, result interface{}) (bool, error)
	*search.SearchHandler
	Error func(context.Context, string, ...map[string]interface{})
	Log   func(context.Context, string, string, bool, string) error
}

func (h *rateHandler) Load(w http.ResponseWriter, r *http.Request) {
	id := sv.GetRequiredParam(w, r)
	if len(id) > 0 {
		var rate Rate
		ok, err := h.load(r.Context(), id, &rate)
		sv.RespondIfFound(w, r, rate, ok, err, h.Error, nil)
	}
}
