package myprofile

import (
	"context"
	"encoding/json"
	sv "github.com/core-go/core"
	v "github.com/core-go/core/v10"
	"net/http"
	"reflect"
)

type MyProfileHandler interface {
	GetMyProfile(w http.ResponseWriter, r *http.Request)
	SaveMyProfile(w http.ResponseWriter, r *http.Request)
	GetMySettings(w http.ResponseWriter, r *http.Request)
	SaveMySettings(w http.ResponseWriter, r *http.Request)
}

func NewMyProfileHandler(service UserService, logError func(context.Context, string, ...map[string]interface{}),
		saveSkills func(ctx context.Context, values []string) (int64, error),
		saveInterests func(ctx context.Context, values []string) (int64, error),
		saveLookingFor func(ctx context.Context, values []string) (int64, error)) (MyProfileHandler, error) {
	var user User
	userType := reflect.TypeOf(user)
	keys, indexes, _ := sv.BuildMapField(userType)
	validator, err := v.NewValidator()
	// s := sv.InitializeStatus(status)
	return &myProfileHandler{service: service, Validate: validator.Validate, LogError: logError, Keys: keys, Indexes: indexes, SaveSkills: saveSkills, SaveInterests: saveInterests, SaveLookingFor: saveLookingFor}, err
}

type myProfileHandler struct {
	service        UserService
	Validate       func(ctx context.Context, model interface{}) ([]sv.ErrorMessage, error)
	LogError       func(context.Context, string, ...map[string]interface{})
	Keys           []string
	Indexes        map[string]int
	SaveSkills     func(ctx context.Context, values []string) (int64, error)
	SaveInterests  func(ctx context.Context, values []string) (int64, error)
	SaveLookingFor func(ctx context.Context, values []string) (int64, error)
}

func (h *myProfileHandler) GetMyProfile(w http.ResponseWriter, r *http.Request) {
	id := sv.GetRequiredParam(w, r)
	if len(id) > 0 {
		user, err := h.service.GetMyProfile(r.Context(), id)
		sv.Return(w, r, user, err, h.LogError, nil)
	}
}
func (h *myProfileHandler) SaveMyProfile(w http.ResponseWriter, r *http.Request) {
	var user User
	r, json, er1 := sv.BuildMapAndCheckId(w, r, &user, h.Keys, h.Indexes)
	if er1 == nil {
		errors, er2 := h.Validate(r.Context(), &user)
		if !sv.HasError(w, r, errors, er2, h.LogError, nil) {
			if h.SaveSkills != nil && len(user.Skills) > 0 {
				skills := make([]string, 0)
				for _, s := range user.Skills {
					skills = append(skills, s.Skill)
				}
				h.SaveSkills(r.Context(), skills)
			}
			if h.SaveInterests != nil && len(user.Interests) > 0 {
				h.SaveInterests(r.Context(), user.Interests)
			}
			if h.SaveLookingFor != nil && len(user.LookingFor) > 0 {
				h.SaveLookingFor(r.Context(), user.LookingFor)
			}
			res, er3 := h.service.SaveMyProfile(r.Context(), json)
			sv.HandleResult(w, r, json, res, er3, h.LogError, nil)
		}
	}
}
func (h *myProfileHandler) GetMySettings(w http.ResponseWriter, r *http.Request) {
	id := sv.GetRequiredParam(w, r, 1)
	if len(id) > 0 {
		settings, err := h.service.GetMySettings(r.Context(), id)
		sv.Return(w, r, settings, err, h.LogError, nil)
	}
}
func (h *myProfileHandler) SaveMySettings(w http.ResponseWriter, r *http.Request) {
	id := sv.GetRequiredParam(w, r, 1)
	if len(id) > 0 {
		var settings Settings
		err := json.NewDecoder(r.Body).Decode(&settings)
		if err != nil {
			http.Error(w, "Invalid Data", http.StatusBadRequest)
			return
		}
		res, err := h.service.SaveMySettings(r.Context(), id, &settings)
		sv.Return(w, r, res, err, h.LogError, nil)
	}
}
