package myprofile

import (
	"context"
	"fmt"
	sv "github.com/core-go/service"
	"net/http"
)

type UserHandler interface {
	//Search(w http.ResponseWriter, r *http.Request)
	Load(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Patch(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	GetMySetting(w http.ResponseWriter, r *http.Request)
	SaveMySetting(w http.ResponseWriter, r *http.Request)
}

func NewUserHandler(service UserService, validate func(ctx context.Context, model interface{}) ([]sv.ErrorMessage, error)) UserHandler {
	//searchModelType := reflect.TypeOf(UserFilter{})
	//modelType := reflect.TypeOf(User{})
	//searchHandler := search.NewSearchHandler(find, modelType, searchModelType, logError, params.Log)
	return &userHandler{service: service}
}

type userHandler struct {
	service UserService
}

func (u *userHandler) Load(w http.ResponseWriter, r *http.Request) {
	id := sv.GetRequiredParam(w, r)
	if len(id) > 0 {
		result, err := u.service.Load(r.Context(), id)
		sv.RespondModel(w, r, result, err, nil, nil)
	}
}

func (u *userHandler) Create(w http.ResponseWriter, r *http.Request) {
	var user User
	user.Phone = ""
	user.Username = ""
	er1 := sv.Decode(w, r, &user)
	if er1 == nil {
		//errors, er2 := u.Validate(r.Context(), &user)
		if !sv.HasError(w, r, nil, nil, 0, nil, nil, "", "") {
			result, er3 := u.service.Create(r.Context(), &user)
			sv.AfterCreated(w, r, &user, result, er3, sv.StatusConfig{}, nil, nil, "", "")
		}
	}
}

func (h *userHandler) Patch(w http.ResponseWriter, r *http.Request) {
	var user User
	r, json, er1 := sv.BuildMapAndCheckId(w, r, &user, nil, nil)
	if er1 == nil {
		id := sv.GetRequiredParam(w, r)
		if len(id) > 0 {
			result, err := h.service.Load(r.Context(), id)
			if err != nil {
				sv.RespondModel(w, r, result, err, nil, nil)
			}
			if result != nil {
				//errors, er2 := h.Validate(r.Context(), &user)
				if !sv.HasError(w, r, nil, nil, 0, nil, nil, "", "") {
					result1, er3 := h.service.Update(r.Context(), &user)
					sv.HandleResult(w, r, json, result1, er3, sv.StatusConfig{}, nil, nil, "", "")
				}
			} else {
				//errors, er2 := h.Validate(r.Context(), &user)
				if !sv.HasError(w, r, nil, nil, 0, nil, nil, "", "") {
					result1, er3 := h.service.Create(r.Context(), &user)
					sv.HandleResult(w, r, json, result1, er3, sv.StatusConfig{}, nil, nil, "", "")
				}
			}
		}

	}
}
func (h *userHandler) Update(w http.ResponseWriter, r *http.Request) {
	var user User
	er1 := sv.DecodeAndCheckId(w, r, &user, nil, nil)
	if er1 == nil {
		//errors, er2 := h.Validate(r.Context(), &user)
		if !sv.HasError(w, r, nil, nil, 0, nil, nil, "", "") {
			result, er3 := h.service.Update(r.Context(), &user)
			sv.HandleResult(w, r, &user, result, er3, sv.StatusConfig{}, nil, nil, "", "")
		}
	}
}
func (h *userHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := sv.GetRequiredParam(w, r)
	if len(id) > 0 {
		result, err := h.service.Delete(r.Context(), id)
		sv.HandleDelete(w, r, result, err, nil, nil, "", "")
	}
}

func (h *userHandler) GetMySetting(w http.ResponseWriter, r *http.Request) {
	id := sv.GetRequiredParam(w, r)
	fmt.Println(id)
	if len(id) > 0 {
		result, err := h.service.Load(r.Context(), id)
		if result.Settings == nil {
			result.Settings.DateFormat = ""
			result.Settings.Language = "English"
			result.Settings.DateTimeFormat = ""
			result.Settings.Notification = true
			sv.RespondModel(w, r, result.Settings, err, nil, nil)
		} else {
			sv.RespondModel(w, r, result.Settings, err, nil, nil)
		}
	}
}
func (h *userHandler) SaveMySetting(w http.ResponseWriter, r *http.Request) {
	//var settings Settings
	var user User
	r, json, er1 := sv.BuildMapAndCheckId(w, r, &user, nil, nil)
	if er1 == nil {
		//errors, er2 := h.Validate(r.Context(), &user.Settings)
		if !sv.HasError(w, r, nil, nil, 0, nil, nil, "", "") {
			result, er3 := h.service.Patch(r.Context(), json)
			sv.HandleResult(w, r, json, result, er3, sv.StatusConfig{}, nil, nil, "", "")
		}
	}

}
