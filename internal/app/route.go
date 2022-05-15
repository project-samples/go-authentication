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

	r.HandleFunc("/authentication/authenticate", app.Authentication.Authenticate).Methods(POST)
	r.HandleFunc("/authentication/signout/{username}", app.SignOut.SignOut).Methods(GET)

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

	user := "/users"
	r.HandleFunc(user+"/search", app.User.Search).Methods(GET, POST)
	r.HandleFunc(user+"/{id}", app.User.Load).Methods(GET)

	r.HandleFunc("/my-profile/{id}", app.MyProfile.GetMyProfile).Methods(GET)
	r.HandleFunc("/my-profile/{id}", app.MyProfile.SaveMyProfile).Methods(PATCH)
	r.HandleFunc("/my-profile/{id}/settings", app.MyProfile.GetMySettings).Methods(GET)
	r.HandleFunc("/my-profile/{id}/settings", app.MyProfile.SaveMySettings).Methods(PATCH)
	return err
}
