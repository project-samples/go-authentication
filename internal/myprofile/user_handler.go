package myprofile

import (
	"context"
	"encoding/json"
	"github.com/core-go/search"
	sv "github.com/core-go/service"
	"io/ioutil"
	"net/http"
	"reflect"
)

type UserHandler interface {
	Search(w http.ResponseWriter, r *http.Request)
	Load(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Patch(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	GetMySetting(w http.ResponseWriter, r *http.Request)
	SaveMySetting(w http.ResponseWriter, r *http.Request)
}

func NewUserHandler(find func(context.Context, interface{}, interface{}, int64, ...int64) (int64, string, error), service UserService, status sv.StatusConfig, logError func(context.Context, string), validate func(ctx context.Context, model interface{}) ([]sv.ErrorMessage, error), action *sv.ActionConfig) UserHandler {
	searchModelType := reflect.TypeOf(UserFilter{})
	modelType := reflect.TypeOf(User{})
	params := sv.CreateParams(modelType, &status, logError, validate, action)
	searchHandler := search.NewSearchHandler(find, modelType, searchModelType, logError, params.Log)
	return &userHandler{service: service, SearchHandler: searchHandler, Params: params}
}

type userHandler struct {
	service UserService
	*search.SearchHandler
	*sv.Params
}

func (u *userHandler) Load(w http.ResponseWriter, r *http.Request) {
	id := sv.GetRequiredParam(w, r)
	if len(id) > 0 {
		result, err := u.service.Load(r.Context(), id)
		sv.RespondModel(w, r, result, err, nil, nil)
	}
}

func (h *userHandler) Create(w http.ResponseWriter, r *http.Request) {
	var user User
	user.Phone = ""
	user.Username = ""
	er1 := sv.Decode(w, r, &user)
	if er1 == nil {
		errors, er2 := h.Validate(r.Context(), &user)
		if !sv.HasError(w, r, errors, er2, *h.Status.ValidationError, h.Error, h.Log, h.Resource, h.Action.Create) {
			result, er3 := h.service.Create(r.Context(), &user)
			sv.AfterCreated(w, r, &user, result, er3, h.Status, h.Error, h.Log, h.Resource, h.Action.Create)
		}
	}
}

func (h *userHandler) Patch(w http.ResponseWriter, r *http.Request) {
	var user User
	r, json, er1 := sv.BuildMapAndCheckId(w, r, &user, h.Keys, h.Indexes)
	if er1 == nil {
		errors, er2 := h.Validate(r.Context(), &user)
		if !sv.HasError(w, r, errors, er2, *h.Status.ValidationError, h.Error, h.Log, h.Resource, h.Action.Patch) {
			result, er3 := h.service.Update(r.Context(), &user)
			sv.HandleResult(w, r, json, result, er3, h.Status, h.Error, h.Log, h.Resource, h.Action.Patch)
		}
	}
}
func (h *userHandler) Update(w http.ResponseWriter, r *http.Request) {
	var user User
	er1 := sv.DecodeAndCheckId(w, r, &user, h.Keys, h.Indexes)
	if er1 == nil {
		errors, er2 := h.Validate(r.Context(), &user)
		if !sv.HasError(w, r, errors, er2, *h.Status.ValidationError, h.Error, h.Log, h.Resource, h.Action.Patch) {
			result, er3 := h.service.Update(r.Context(), &user)
			sv.HandleResult(w, r, &user, result, er3, h.Status, h.Error, h.Log, h.Resource, h.Action.Patch)
		}
	}
}
func (h *userHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := sv.GetRequiredParam(w, r)
	if len(id) > 0 {
		result, err := h.service.Delete(r.Context(), id)
		sv.HandleDelete(w, r, result, err, h.Error, h.Log, h.Resource, h.Action.Delete)
	}
}

func (h *userHandler) GetMySetting(w http.ResponseWriter, r *http.Request) {
	id := sv.GetRequiredParam(w, r)
	if len(id) > 0 {
		result, err := h.service.Load(r.Context(), id)
		sv.RespondModel(w, r, result.Settings, err, h.Error, h.Log)
	}
}
func (h *userHandler) SaveMySetting(w http.ResponseWriter, r *http.Request) {
	var settings Settings
	var user User
	id := sv.GetRequiredParam(w, r)
	if len(id) == 0 {
		sv.HasError(w, r, nil, nil, 0, h.Error, h.Log, h.Resource, h.Action.Update)
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		sv.HasError(w, r, nil, err, 0, h.Error, h.Log, h.Resource, h.Action.Update)
	}
	err = json.Unmarshal(body, &settings)
	if err == nil {
		result1, _ := h.service.Load(r.Context(), id)
		user = *result1
		settings.UserId = id
		user.Settings = &settings
		errors, er2 := h.Validate(r.Context(), &settings)
		if !sv.HasError(w, r, errors, er2, *h.Status.ValidationError, h.Error, h.Log, h.Resource, h.Action.Update) {
			result, er3 := h.service.Update(r.Context(), &user)
			sv.HandleResult(w, r, user, result, er3, h.Status, h.Error, h.Log, h.Resource, h.Action.Update)
		}
	}
}
