package myprofile

import (
	"context"
	"encoding/json"
	"go-service/internal/usecase/upload"

	// "go-service/internal/usecase/upload"

	"net/http"
	"path/filepath"
	"reflect"
	"strings"

	sv "github.com/core-go/service"
	v "github.com/core-go/service/v10"
	"github.com/core-go/storage"
)

type MyProfileHandler interface {
	GetMyProfile(w http.ResponseWriter, r *http.Request)
	SaveMyProfile(w http.ResponseWriter, r *http.Request)
	GetMySettings(w http.ResponseWriter, r *http.Request)
	SaveMySettings(w http.ResponseWriter, r *http.Request)
	upload.UploadHander
}

const contentTypeHeader = "Content-Type"

func NewMyProfileHandler(service UserService, logError func(context.Context, string), status *sv.StatusConfig,
	serviceStorage storage.StorageService, provider string, generalDirectory string, keyFile string, directory string,
	saveSkills func(ctx context.Context, values []string) (int64, error),
	saveInterests func(ctx context.Context, values []string) (int64, error),
	saveLookingFor func(ctx context.Context, values []string) (int64, error), uploadService upload.UploadService) MyProfileHandler {
	var user User
	userType := reflect.TypeOf(user)
	keys, indexes, _ := sv.BuildMapField(userType)
	validator := v.NewValidator()
	s := sv.InitializeStatus(status)
	uploadHandler := upload.NewUploadHandler(uploadService, logError, status, provider, generalDirectory, keyFile)

	return &myProfileHandler{service: service, Validate: validator.Validate, LogError: logError, Keys: keys,
		Indexes: indexes, Status: s, ServiceStorage: serviceStorage,
		Provider: provider, GeneralDirectory: generalDirectory, KeyFile: keyFile,
		Directory: directory, SaveSkills: saveSkills, SaveInterests: saveInterests,
		SaveLookingFor: saveLookingFor, uploadHandler: uploadHandler,
	}
}

type myProfileHandler struct {
	service        UserService
	Validate       func(ctx context.Context, model interface{}) ([]sv.ErrorMessage, error)
	LogError       func(context.Context, string)
	Keys           []string
	Indexes        map[string]int
	Status         sv.StatusConfig
	SaveSkills     func(ctx context.Context, values []string) (int64, error)
	SaveInterests  func(ctx context.Context, values []string) (int64, error)
	SaveLookingFor func(ctx context.Context, values []string) (int64, error)

	ServiceStorage   storage.StorageService
	uploadHandler    upload.UploadHander
	Provider         string
	GeneralDirectory string
	Directory        string
	KeyFile          string
}

func (h *myProfileHandler) GetMyProfile(w http.ResponseWriter, r *http.Request) {
	id := sv.GetRequiredParam(w, r)
	if len(id) > 0 {
		user, err := h.service.GetMyProfile(r.Context(), id)
		sv.RespondModel(w, r, user, err, h.LogError, nil)
	}
}
func (h *myProfileHandler) SaveMyProfile(w http.ResponseWriter, r *http.Request) {
	var user User
	r, json, er1 := sv.BuildMapAndCheckId(w, r, &user, h.Keys, h.Indexes)
	if er1 == nil {
		errors, er2 := h.Validate(r.Context(), &user)
		if !sv.HasError(w, r, errors, er2, *h.Status.ValidationError, h.LogError, nil) {
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
			sv.HandleResult(w, r, json, res, er3, h.Status, h.LogError, nil)
		}
	}
}
func (h *myProfileHandler) GetMySettings(w http.ResponseWriter, r *http.Request) {
	id := sv.GetRequiredParam(w, r, 1)
	if len(id) > 0 {
		settings, err := h.service.GetMySettings(r.Context(), id)
		sv.RespondModel(w, r, settings, err, h.LogError, nil)
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
		sv.RespondModel(w, r, res, err, h.LogError, nil)
	}
}

func (u *myProfileHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	u.uploadHandler.UploadImage(w, r)
}

func (u *myProfileHandler) UploadGallery(w http.ResponseWriter, r *http.Request) {
	u.uploadHandler.UploadGallery(w, r)
}

func (u *myProfileHandler) UploadCover(w http.ResponseWriter, r *http.Request) {
	u.uploadHandler.UploadCover(w, r)
}

func (u *myProfileHandler) DeleteGalleryFile(w http.ResponseWriter, r *http.Request) {
	u.uploadHandler.DeleteGalleryFile(w, r)
}

func respond(w http.ResponseWriter, code int, result interface{}) {
	res, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(res)
}
func getExt(file string) string {
	ext := filepath.Ext(file)
	if strings.HasPrefix(ext, ":") {
		ext = ext[1:]
		return ext
	}
	return ext
}
