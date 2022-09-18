package app

import (
	"context"

	. "github.com/core-go/security"
	"github.com/gorilla/mux"
)

func Route(r *mux.Router, context context.Context, root Config) error {
	app, err := NewApp(context, root)
	if err != nil {
		return err
	}
	r.HandleFunc("/health", app.Health.Check).Methods(GET)

	r.HandleFunc("/authenticate", app.Authentication.Authenticate).Methods(POST)
	r.HandleFunc("/authenticate/signout/{username}", app.SignOut.SignOut).Methods(GET)

	r.HandleFunc("/password/change", app.Password.ChangePassword).Methods(POST)
	r.HandleFunc("/password/forgot", app.Password.ForgotPassword).Methods(POST)
	r.HandleFunc("/password/reset", app.Password.ResetPassword).Methods(POST)

	r.HandleFunc("/signup/signup", app.SignUp.SignUp).Methods(POST)
	r.HandleFunc("/signup/verify/{userId}/{code}", app.SignUp.VerifyUser).Methods(GET)

	r.HandleFunc("/oauth2/configurations/{type}", app.OAuth2.Configuration).Methods(GET)
	r.HandleFunc("/oauth2/configurations", app.OAuth2.Configurations).Methods(GET)
	r.HandleFunc("/oauth2/authenticate", app.OAuth2.Authenticate).Methods(POST)

	r.HandleFunc("/skills", app.Skill.Query).Methods(GET)
	r.HandleFunc("/interests", app.Interest.Query).Methods(GET)
	r.HandleFunc("/looking-for", app.LookingFor.Query).Methods(GET)

	r.HandleFunc("/my-profile/{id}", app.MyProfile.GetMyProfile).Methods(GET)
	r.HandleFunc("/my-profile/{id}", app.MyProfile.SaveMyProfile).Methods(PATCH)
	r.HandleFunc("/my-profile/{id}/settings", app.MyProfile.GetMySettings).Methods(GET)
	r.HandleFunc("/my-profile/{id}/settings", app.MyProfile.SaveMySettings).Methods(PATCH)
	r.HandleFunc("/my-profile/{id}/upload", app.MyProfile.UploadImage).Methods(POST)
	r.HandleFunc("/my-profile/{id}/gallery", app.MyProfile.UploadGallery).Methods(POST)
	r.HandleFunc("/my-profile/{id}/cover", app.MyProfile.UploadCover).Methods(POST)
	r.HandleFunc("/my-profile/{id}/gallery", app.MyProfile.DeleteGalleryFile).Methods(DELETE)
	// r.HandleFunc("/my-profile/delete/{id}", app.MyProfile.DeleteFile).Methods(DELETE)
	//r.HandleFunc("/my-profile/image/{id}", app.MyProfile.GetAvt).Methods(GET)

	return err
}
