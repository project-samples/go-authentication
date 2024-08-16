package myprofile

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"

	sv "github.com/core-go/core"
	v "github.com/core-go/core/validator"
)

func NewMyProfileHandler(service UserService, logError func(context.Context, string, ...map[string]interface{}),
	saveSkills func(ctx context.Context, values []string) (int64, error),
	saveInterests func(ctx context.Context, values []string) (int64, error),
	saveLookingFor func(ctx context.Context, values []string) (int64, error)) (*MyProfileHandler, error) {
	keys, indexes, _ := sv.BuildMapField(reflect.TypeOf(User{}))
	validator, err := v.NewValidator[*User]()
	if err != nil {
		return nil, err
	}
	return &MyProfileHandler{service: service, Validate: validator.Validate, LogError: logError, Keys: keys, Indexes: indexes, SaveSkills: saveSkills, SaveInterests: saveInterests, SaveLookingFor: saveLookingFor}, err
}

type MyProfileHandler struct {
	service        UserService
	Validate       func(ctx context.Context, user *User) ([]sv.ErrorMessage, error)
	LogError       func(context.Context, string, ...map[string]interface{})
	Keys           []string
	Indexes        map[string]int
	SaveSkills     func(ctx context.Context, values []string) (int64, error)
	SaveInterests  func(ctx context.Context, values []string) (int64, error)
	SaveLookingFor func(ctx context.Context, values []string) (int64, error)
}

func (h *MyProfileHandler) GetMyProfile(w http.ResponseWriter, r *http.Request) {
	id := sv.GetRequiredParam(w, r)
	if len(id) > 0 {
		user, err := h.service.GetMyProfile(r.Context(), id)
		sv.Return(w, r, user, err, h.LogError, nil)
	}
}
func (h *MyProfileHandler) SaveMyProfile(w http.ResponseWriter, r *http.Request) {
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
func (h *MyProfileHandler) GetMySettings(w http.ResponseWriter, r *http.Request) {
	id := sv.GetRequiredParam(w, r, 1)
	if len(id) > 0 {
		settings, err := h.service.GetMySettings(r.Context(), id)
		sv.Return(w, r, settings, err, h.LogError, nil)
	}
}
func (h *MyProfileHandler) SaveMySettings(w http.ResponseWriter, r *http.Request) {
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
