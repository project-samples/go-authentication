package myprofile

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go-service/internal/usecase/upload"

	"github.com/core-go/service/shortid"

	// "go-service/internal/usecase/upload"
	"io"
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
	UploadImage(w http.ResponseWriter, r *http.Request)
	UploadGallery(w http.ResponseWriter, r *http.Request)
	UploadCover(w http.ResponseWriter, r *http.Request)
	DeleteGallery(w http.ResponseWriter, r *http.Request)
	DeleteFile(w http.ResponseWriter, r *http.Request)
}

const contentTypeHeader = "Content-Type"

func NewMyProfileHandler(service UserService, logError func(context.Context, string), status *sv.StatusConfig,
	serviceStorage storage.StorageService, provider string, generalDirectory string, keyFile string, directory string,
	saveSkills func(ctx context.Context, values []string) (int64, error),
	saveInterests func(ctx context.Context, values []string) (int64, error),
	saveLookingFor func(ctx context.Context, values []string) (int64, error), uploadHandler upload.UploadHandler) MyProfileHandler {
	var user User
	userType := reflect.TypeOf(user)
	keys, indexes, _ := sv.BuildMapField(userType)
	validator := v.NewValidator()
	s := sv.InitializeStatus(status)
	return &myProfileHandler{service: service, Validate: validator.Validate, LogError: logError, Keys: keys,
		Indexes: indexes, Status: s, ServiceStorage: serviceStorage,
		Provider: provider, GeneralDirectory: generalDirectory, KeyFile: keyFile,
		Directory: directory, SaveSkills: saveSkills, SaveInterests: saveInterests,
		SaveLookingFor: saveLookingFor, uploadHandler: uploadHandler}
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
	uploadHandler    upload.UploadHandler
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
	err := r.ParseMultipartForm(200000)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	formdata := r.MultipartForm // ok, no problem so far, read the Form data

	//get the *fileheaders
	files := formdata.File[u.KeyFile] // grab the filenames
	_, handler, _ := r.FormFile(u.KeyFile)
	contentType := handler.Header.Get(contentTypeHeader)
	if len(contentType) == 0 {
		contentType = getExt(handler.Filename)
	}
	generateStr, _ := shortid.Generate(r.Context())
	var list []upload.FileInfo
	for i, _ := range files { // loop through the files one by one
		file, err := files[i].Open()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer file.Close()
		out := bytes.NewBuffer(nil)

		_, err = io.Copy(out, file) // file not files[i] !

		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		bytes := out.Bytes()
		name := generateStr + "_" + files[i].Filename
		list = append(list, upload.FileInfo{name, bytes})

	}

	id := sv.GetRequiredParam(w, r, 1)
	rs, err := u.uploadHandler.UploadImage(id, list, contentType, r)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	respond(w, http.StatusOK, rs)
}

func (u *myProfileHandler) UploadGallery(w http.ResponseWriter, r *http.Request) {

	err1 := r.ParseMultipartForm(32 << 20)
	if err1 != nil {
		http.Error(w, "not avalable", http.StatusInternalServerError)
		return
	}

	file, handler, err2 := r.FormFile(u.KeyFile)
	if err2 != nil {
		http.Error(w, "not avalable", http.StatusInternalServerError)
		return
	}

	bufferFile := bytes.NewBuffer(nil)
	_, err3 := io.Copy(bufferFile, file)
	if err3 != nil {
		http.Error(w, "not avalable", http.StatusInternalServerError)
		return
	}

	defer file.Close()
	bytes := bufferFile.Bytes()
	contentType := handler.Header.Get(contentTypeHeader)
	if len(contentType) == 0 {
		contentType = getExt(handler.Filename)
	}

	id := sv.GetRequiredParam(w, r, 1)
	rs, err5 := u.uploadHandler.UploadGallery(id, r.FormValue("source"), handler.Filename, contentType, bytes, r)

	if err5 != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	sv.HandleResult(w, r, rs, 1, err5, u.Status, u.LogError, nil)
	//respond(w, http.StatusOK, res)

}

func (u *myProfileHandler) UploadCover(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(200000)
	if err != nil {
		http.Error(w, "not avalable", http.StatusInternalServerError)
		return
	}
	formdata := r.MultipartForm // ok, no problem so far, read the Form data

	//get the *fileheaders
	files := formdata.File[u.KeyFile] // grab the filenames
	_, handler, _ := r.FormFile(u.KeyFile)
	contentType := handler.Header.Get(contentTypeHeader)
	if len(contentType) == 0 {
		contentType = getExt(handler.Filename)
	}
	generateStr, _ := shortid.Generate(r.Context())
	var list []upload.FileInfo
	for i, _ := range files { // loop through the files one by one
		file, err := files[i].Open()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer file.Close()
		out := bytes.NewBuffer(nil)

		_, err = io.Copy(out, file) // file not files[i] !

		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		bytes := out.Bytes()
		name := generateStr + "_" + files[i].Filename
		list = append(list, upload.FileInfo{name, bytes})

	}

	id := sv.GetRequiredParam(w, r, 1)
	rs, err := u.uploadHandler.UploadCover(id, list, contentType, r)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	respond(w, http.StatusOK, rs)
}

func (u *myProfileHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(200000)
	if err != nil {
		http.Error(w, "not avalable", http.StatusInternalServerError)
		return
	}
	formdata := r.MultipartForm // ok, no problem so far, read the Form data

	//get the *fileheaders
	files := formdata.File[u.KeyFile] // grab the filenames
	var directory string
	if u.Provider == "google-storage" {
		directory = u.Directory
	} else {
		directory = u.GeneralDirectory
	}
	_, handler, _ := r.FormFile(u.KeyFile)
	contentType := handler.Header.Get(contentTypeHeader)
	if len(contentType) == 0 {
		contentType = getExt(handler.Filename)
	}
	var listURL []string
	for i, _ := range files { // loop through the files one by one
		file, err := files[i].Open()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer file.Close()
		out := bytes.NewBuffer(nil)

		_, err = io.Copy(out, file) // file not files[i] !

		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		bytes := out.Bytes()
		rs, err2 := u.ServiceStorage.Upload(r.Context(), directory, files[i].Filename, bytes, contentType)
		if err2 != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		listURL = append(listURL, rs)

	}

	// res, _ := json.Marshal(rs)

	id := sv.GetRequiredParam(w, r, 1)
	if len(id) > 0 {
		result, err4 := u.service.insertImage(r.Context(), id, listURL[0])
		if err4 != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		sv.HandleResult(w, r, listURL, result, err4, u.Status, u.LogError, nil)
		//respond(w, http.StatusOK, res)
	}

}

func (u *myProfileHandler) DeleteGallery(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	fmt.Println(url)
	if len(url) == 0 {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	i := strings.LastIndex(url, "/")
	filename := ""
	if i <= 0 {
		http.Error(w, "require id", http.StatusBadRequest)
		return
	}
	filename = url[i+1:]
	i = strings.LastIndex(filename, "?")
	filename = filename[:i]
	if len(filename) == 0 {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	rs, err := u.uploadHandler.DeleteFile(filename, r)
	fmt.Print(rs)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	id := sv.GetRequiredParam(w, r, 1)
	result, err4 := u.service.deleteGallery(r.Context(), id, url)
	if err4 != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	sv.HandleResult(w, r, result, result, err4, u.Status, u.LogError, nil)
}

func (u *myProfileHandler) DeleteFile(w http.ResponseWriter, r *http.Request) {
	i := strings.LastIndex(r.RequestURI, "/")
	filename := ""
	if i <= 0 {
		http.Error(w, "require id", http.StatusBadRequest)
		return
	}
	filename = r.RequestURI[i+1:]

	var filepath string
	if u.Provider == "drop-box" {
		filepath = fmt.Sprintf("/%s/%s", u.GeneralDirectory, filename)
	} else {
		filepath = filename
	}

	rs, err := u.ServiceStorage.Delete(r.Context(), filepath)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	var msg string
	if rs {
		msg = fmt.Sprintf("file '%s' has been deleted successfully!!!", filename)
	} else {
		msg = fmt.Sprintf("delete file '%s' failed!!!", filename)
	}
	respond(w, http.StatusOK, msg)
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
