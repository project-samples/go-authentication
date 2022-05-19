package upload

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/core-go/storage"
	"io"
	"net/http"
	"path/filepath"
	"strings"
)

const contentTypeHeader = "Content-Type"

type UploadHandler interface {
	UploadFile(w http.ResponseWriter, r *http.Request)
	DeleteFile(w http.ResponseWriter, r *http.Request)
}

type uploadHandler struct {
	Service          storage.StorageService
	Provider         string
	GeneralDirectory string
	Directory        string
	KeyFile          string
}

func NewUploadHandler(service storage.StorageService, provider string, generalDirectory string, keyFile string, directory string) UploadHandler {
	return &uploadHandler{Service: service, Provider: provider, GeneralDirectory: generalDirectory, KeyFile: keyFile, Directory: directory}
}

func (u *uploadHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, "not avalable", http.StatusInternalServerError)
		return
	}

	file, handler, err0 := r.FormFile(u.KeyFile)
	if err0 != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	bufferFile := bytes.NewBuffer(nil)
	_, err1 := io.Copy(bufferFile, file)
	if err1 != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	defer file.Close()

	bytes := bufferFile.Bytes()
	contentType := handler.Header.Get(contentTypeHeader)
	if len(contentType) == 0 {
		contentType = getExt(handler.Filename)
	}

	var directory string
	if u.Provider == "google-storage" {
		directory = u.Directory
	} else {
		directory = u.GeneralDirectory
	}

	rs, err2 := u.Service.Upload(r.Context(), directory, handler.Filename, bytes, contentType)
	if err2 != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	respond(w, http.StatusOK, rs)
}

func (u *uploadHandler) DeleteFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("asd")
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