package upload

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/core-go/service/shortid"
	"github.com/core-go/storage"
)

const contentTypeHeader = "Content-Type"

type UploadHandler interface {
	UploadFile(fileName string, contentType string, data []byte, r *http.Request) (rs string, errorRespone error)
	DeleteFile(url string, r *http.Request) (bool, error)
	UploadGallery(id string, source string, name string, contentType string, data []byte, r *http.Request) ([]UploadInfo, error)
	DeleteGallery(url string, id string, r *http.Request) (bool, error)
	UploadCover(id string, data []FileInfo, contentType string, r *http.Request) (string, error)
	UploadImage(id string, data []FileInfo, contentType string, r *http.Request) (string, error)
}

type uploadHandler struct {
	Service          storage.StorageService
	Provider         string
	GeneralDirectory string
	Directory        string
	KeyFile          string
	SizesImage       []int
	SizesCover       []int
	loadData         func(ctx context.Context, id string) (interface{}, error)
	patchData        func(ctx context.Context, user map[string]interface{}) (int64, error)
}
type UploadInfo struct {
	Source string `json:"source,omitempty" gorm:"column:source" bson:"source,omitempty" dynamodbav:"source,omitempty" firestore:"source,omitempty"`
	Url    string `json:"url,omitempty" gorm:"column:url" bson:"url,omitempty" dynamodbav:"url,omitempty" firestore:"url,omitempty"`
	Type   string `json:"category,omitempty" gorm:"column:category" bson:"category,omitempty" dynamodbav:"category,omitempty" firestore:"category,omitempty"`
}

type FileInfo struct {
	Name string `json:"name,omitempty" gorm:"column:name" bson:"name,omitempty" dynamodbav:"name,omitempty" firestore:"name,omitempty"`
	Data []byte `json:"data,omitempty" gorm:"column:data" bson:"data,omitempty" dynamodbav:"data,omitempty" firestore:"data,omitempty"`
}

func NewUploadHandler(service storage.StorageService, provider string, generalDirectory string,
	keyFile string, directory string, loadData func(ctx context.Context, id string) (interface{}, error), patchData func(ctx context.Context, user map[string]interface{}) (int64, error)) UploadHandler {
	return &uploadHandler{Service: service, Provider: provider, GeneralDirectory: generalDirectory, KeyFile: keyFile, Directory: directory, loadData: loadData, patchData: patchData}
}

func (u *uploadHandler) UploadCover(id string, data []FileInfo, contentType string, r *http.Request) (string, error) {
	//delete

	result, _ := u.loadData(r.Context(), id)
	val := reflect.ValueOf(result).Elem()
	cover := val.FieldByName("CoverURL").Interface().(string)

	sizes := []int{576, 768}
	if cover != "" {
		_, err := u.DeleteFileUpload(sizes, cover, r)
		if err != nil {
			return "", errors.New("internal server error")
		}
	}
	//upload
	var newUrl string
	for i := range data {
		file := data[i]
		rs, errorRespone := u.UploadFile(file.Name, contentType, file.Data, r)
		if errorRespone != nil {
			return "", errorRespone
		}
		if i != 0 {
			continue
		}
		newUrl = rs
	}
	user := make(map[string]interface{})
	user["id"] = id
	user["coverURL"] = newUrl
	_, err1 := u.patchData(r.Context(), user)
	if err1 != nil {
		return "", err1
	}
	return newUrl, nil
}

func (u *uploadHandler) UploadImage(id string, data []FileInfo, contentType string, r *http.Request) (string, error) {
	//delete

	result, _ := u.loadData(r.Context(), id)
	val := reflect.ValueOf(result).Elem()
	url := val.FieldByName("ImageURL").Interface().(string)

	sizes := []int{40, 400}
	if url != "" {
		_, err := u.DeleteFileUpload(sizes, url, r)
		if err != nil {
			return "", errors.New("internal server error")
		}
	}
	//upload
	var newUrl string
	for i := range data {
		file := data[i]
		rs, errorRespone := u.UploadFile(file.Name, contentType, file.Data, r)
		if errorRespone != nil {
			return "", errorRespone
		}
		if i != 0 {
			continue
		}
		newUrl = rs
	}
	user := make(map[string]interface{})
	user["id"] = id
	user["imageURL"] = newUrl
	_, err1 := u.patchData(r.Context(), user)
	if err1 != nil {
		return "", err1
	}
	return newUrl, nil
}

func (u *uploadHandler) UploadGallery(id string, source string, name string, contentType string, data []byte, r *http.Request) ([]UploadInfo, error) {
	sid, _ := shortid.Generate(r.Context())
	fileName := removeExt(name) + "_" + sid + getExt(name)

	rs, errorRespone := u.UploadFile(fileName, contentType, data, r)
	if errorRespone != nil {
		return nil, errorRespone
	}
	result, _ := u.loadData(r.Context(), id)
	val := reflect.ValueOf(result).Elem()
	fmt.Print(val.FieldByName("gallery").IsValid())
	gallery := []UploadInfo{}
	if val.FieldByName("gallery").IsValid() {
		gallery = val.FieldByName("gallery").Interface().([]UploadInfo)
	}

	gallery = append(gallery, UploadInfo{Source: source, Type: strings.Split(contentType, "/")[0], Url: rs})
	user := make(map[string]interface{})
	user["userId"] = id
	user["gallery"] = gallery

	_, err := u.patchData(r.Context(), user)
	if err != nil {
		return nil, err
	}
	return gallery, nil
}

func (u *uploadHandler) UploadFile(fileName string, contentType string, data []byte, r *http.Request) (rs string, errorRespone error) {
	directory := u.Directory
	rs, err2 := u.Service.Upload(r.Context(), directory, fileName, data, contentType)
	if err2 != nil {
		errorRespone = err2
		return
	}
	return
}

func (u *uploadHandler) DeleteGallery(url string, id string, r *http.Request) (bool, error) {
	// a = struct {
	// 	UserId string
	// }{}

	test, _ := u.loadData(r.Context(), id)
	val := reflect.ValueOf(test).Elem()
	gallery := val.FieldByName("Gallery").Interface().([]UploadInfo)
	// if find url in gallery
	idx := -1
	for i := range gallery {
		if gallery[i].Url == url {
			idx = i
		}
	}
	if idx == -1 {
		return false, nil
	}
	_, err := u.DeleteFile(url, r)
	if err != nil {
		return false, err
	}
	gallery = append(gallery[:idx], gallery[idx+1:]...)
	user := make(map[string]interface{})
	user["userId"] = id
	user["gallery"] = gallery
	_, err2 := u.patchData(r.Context(), user)
	if err2 != nil {
		return false, err2
	}
	return true, err2
}
func (u *uploadHandler) DeleteFileUpload(sizes []int, url string, r *http.Request) (bool, error) {
	// i := strings.LastIndex(url, "/")
	// filename := ""
	// if i <= 0 {
	// 	return false, errors.New("internal server error")
	// }
	// filename = url[i+1:]
	// i = strings.LastIndex(filename, "?")
	// filename = filename[:i]
	arrOrigin := strings.Split(url, "/")
	delOriginUrl := arrOrigin[len(arrOrigin)-2] + "/" + arrOrigin[len(arrOrigin)-1]
	rs, err := u.DeleteFile(delOriginUrl, r)
	fmt.Print(rs, err)
	// if err != nil {
	// 	return false, errors.New("internal server error")
	// }
	for i := range sizes {
		oldUrl := removeExt(url) + "_" + strconv.Itoa(sizes[i]) + getExt(url)
		arr := strings.Split(oldUrl, "/")
		delUrl := arr[len(arr)-2] + "/" + arr[len(arr)-1]
		rss, err := u.DeleteFile(delUrl, r)
		fmt.Print(rss, err)
	}
	return true, nil
}

func (u *uploadHandler) DeleteFile(url string, r *http.Request) (bool, error) {
	rs, err := u.Service.Delete(r.Context(), url)
	return rs, err
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

func removeExt(file string) string {
	return file[:len(file)-len(filepath.Ext(file))]
}
