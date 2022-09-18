package upload

import (
	"bytes"
	"context"

	"github.com/core-go/service/shortid"

	// "go-service/internal/usecase/upload"
	"io"
	"net/http"
	"reflect"

	sv "github.com/core-go/service"
	v "github.com/core-go/service/v10"
)

type Upload struct {
	UserId string `json:"id,omitempty" gorm:"column:id;primary_key" bson:"_id,omitempty" validate:"required,max=40"`
	Source string `json:"source,omitempty" gorm:"column:source" bson:"source,omitempty" dynamodbav:"source,omitempty" firestore:"source,omitempty"`
	Type   string `json:"category,omitempty" gorm:"column:category" bson:"category,omitempty" dynamodbav:"category,omitempty" firestore:"category,omitempty"`
	Name   string `json:"name,omitempty" gorm:"column:name" bson:"name,omitempty" dynamodbav:"name,omitempty" firestore:"name,omitempty"`
	Data   []byte `json:"data,omitempty" gorm:"column:data" bson:"data,omitempty" dynamodbav:"data,omitempty" firestore:"data,omitempty"`
}

type User struct {
	UserId     string       `json:"userId,omitempty" gorm:"column:userid;primary_key" bson:"_id,omitempty" dynamodbav:"userId,omitempty" firestore:"userId,omitempty" validate:"required,max=40"`
	Username   string       `json:"username,omitempty" gorm:"column:username" bson:"username,omitempty" dynamodbav:"username,omitempty" firestore:"username,omitempty" validate:"required,username,max=100"`
	Email      string       `json:"email,omitempty" gorm:"column:email" bson:"email,omitempty" dynamodbav:"email,omitempty" firestore:"email,omitempty" validate:"email,max=100"`
	Phone      string       `json:"phone,omitempty" gorm:"column:phone" bson:"phone,omitempty" dynamodbav:"phone,omitempty" firestore:"phone,omitempty" validate:"required,phone,max=18"`
	FirstName  string       `json:"firstName,omitempty" gorm:"column:firstName" bson:"firstName,omitempty" dynamodbav:"firstName,omitempty" firestore:"firstName,omitempty" `
	LastName   string       `json:"lastName,omitempty" gorm:"column:lastName" bson:"lastName,omitempty" dynamodbav:"lastName,omitempty" firestore:"lastName,omitempty" `
	ImageURL   string       `json:"imageURL,omitempty" gorm:"column:imageurl" bson:"imageURL,omitempty" dynamodbav:"imageURL,omitempty" firestore:"imageURL,omitempty" `
	Occupation string       `json:"occupation,omitempty" gorm:"column:occupation" bson:"occupation,omitempty" dynamodbav:"occupation,omitempty" firestore:"occupation,omitempty" `
	Company    string       `json:"company,omitempty" gorm:"column:company" bson:"company,omitempty" dynamodbav:"company,omitempty" firestore:"company,omitempty" `
	Interests  []string     `json:"interests,omitempty" gorm:"column:interests" bson:"interests,omitempty" dynamodbav:"interests,omitempty" firestore:"interests,omitempty"`
	LookingFor []string     `json:"lookingFor,omitempty" gorm:"column:lookingfor" bson:"lookingFor,omitempty" dynamodbav:"lookingFor,omitempty" firestore:"lookingFor,omitempty"`
	Gallery    []UploadInfo `json:"gallery,omitempty" gorm:"column:gallery"`
	CoverURL   string       `json:"coverURL,omitempty" gorm:"column:coverurl" bson:"coverURL,omitempty" dynamodbav:"coverURL,omitempty" firestore:"coverURL,omitempty" `
}

func NewUploadHandler(service UploadService, logError func(context.Context, string), status *sv.StatusConfig,
	provider string, generalDirectory string, keyFile string,
) *uploadHander {
	var user User
	userType := reflect.TypeOf(user)
	keys, indexes, _ := sv.BuildMapField(userType)
	validator := v.NewValidator()
	s := sv.InitializeStatus(status)
	return &uploadHander{service: service, Validate: validator.Validate, LogError: logError, Keys: keys,
		Indexes: indexes, Status: s, KeyFile: keyFile,
	}
}

type UploadHander interface {
	UploadImage(w http.ResponseWriter, r *http.Request)
	UploadGallery(w http.ResponseWriter, r *http.Request)
	UploadCover(w http.ResponseWriter, r *http.Request)
	DeleteGalleryFile(w http.ResponseWriter, r *http.Request)
}
type uploadHander struct {
	service  UploadService
	Validate func(ctx context.Context, model interface{}) ([]sv.ErrorMessage, error)
	LogError func(context.Context, string)
	Keys     []string
	Indexes  map[string]int
	Status   sv.StatusConfig
	// uploadHandler    upload.UploadService
	// Provider         string
	// GeneralDirectory string
	// Directory        string
	KeyFile    string
	getUploads func(id string) []UploadInfo
	generateId func() string
	sizesCover []int
	sizesImage []int
}

func (u *uploadHander) UploadImage(w http.ResponseWriter, r *http.Request) {
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
	var list []UploadData
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
		list = append(list, UploadData{name, bytes})

	}

	id := sv.GetRequiredParam(w, r, 1)
	rs, err := u.service.UploadImage(id, list, contentType, r, []int{40, 400})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	respond(w, http.StatusOK, rs)
}

func (u *uploadHander) UploadGallery(w http.ResponseWriter, r *http.Request) {

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
	uploadFile := Upload{Source: r.FormValue("source"), Data: bytes, Name: handler.Filename, UserId: id, Type: contentType}
	rs, err5 := u.service.UploadGallery(uploadFile, r)

	if err5 != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	sv.HandleResult(w, r, rs, 1, err5, u.Status, u.LogError, nil)
	//respond(w, http.StatusOK, res)

}

func (u *uploadHander) UploadCover(w http.ResponseWriter, r *http.Request) {
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
	var list []UploadData
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
		list = append(list, UploadData{name, bytes})

	}

	id := sv.GetRequiredParam(w, r, 1)

	rs, err := u.service.UploadCover(id, list, contentType, r, []int{576, 768})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	respond(w, http.StatusOK, rs)
}

func (u *uploadHander) UploadAvatar(w http.ResponseWriter, r *http.Request) {
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
	var list []UploadData
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
		list = append(list, UploadData{name, bytes})

	}

	id := sv.GetRequiredParam(w, r, 1)
	rs, err := u.service.UploadImage(id, list, contentType, r, []int{40, 400})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	respond(w, http.StatusOK, rs)
}

func (u *uploadHander) DeleteGalleryFile(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	id := sv.GetRequiredParam(w, r, 1)

	result, err4 := u.service.DeleteGalleryFile(id, url, r)
	if err4 != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	sv.HandleResult(w, r, result, result, err4, u.Status, u.LogError, nil)
}

// func respond(w http.ResponseWriter, code int, result interface{}) {
// 	res, _ := json.Marshal(result)
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(code)
// 	w.Write(res)
// }
// func getExt(file string) string {
// 	ext := filepath.Ext(file)
// 	if strings.HasPrefix(ext, ":") {
// 		ext = ext[1:]
// 		return ext
// 	}
// 	return ext
// }
