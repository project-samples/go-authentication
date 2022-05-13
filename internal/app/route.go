package app

import (
	"context"
	. "github.com/core-go/security"
	"github.com/gorilla/mux"
)

func Route(r *mux.Router, context context.Context, root Root) error {
	app, err := NewApp(context, root)
	if err != nil {
		return err
	}
	r.HandleFunc("/health", app.HealthHandler.Check).Methods(GET)

	r.HandleFunc("/authentication", app.AuthenticationHandler.Authenticate).Methods(POST)
	r.HandleFunc("/authentication/signout/{username}", app.SignOutHandler.SignOut).Methods(GET)

	r.HandleFunc("/password/change", app.PasswordHandler.ChangePassword).Methods(POST)
	r.HandleFunc("/password/forgot", app.PasswordHandler.ForgotPassword).Methods(POST)
	r.HandleFunc("/password/reset", app.PasswordHandler.ResetPassword).Methods(POST)

	r.HandleFunc("/signup/signup", app.SignUpHandler.SignUp).Methods(POST)
	r.HandleFunc("/signup/verify/{userId}/{code}", app.SignUpHandler.VerifyUser).Methods(GET)

	r.HandleFunc("/oauth2/configurations/{type}", app.OAuth2Handler.Configuration).Methods(GET)
	r.HandleFunc("/oauth2/configurations", app.OAuth2Handler.Configurations).Methods(GET)
	r.HandleFunc("/oauth2/authenticate", app.OAuth2Handler.Authenticate).Methods(POST)

	r.HandleFunc("/my-profile/{id}", app.User.Patch).Methods(PATCH)
	r.HandleFunc("/my-profile/{id}", app.User.Load).Methods(GET)
	//r.HandleFunc("/my-profile/search", app.User.Search).Methods(GET, POST)
	r.HandleFunc("/my-profile/{id}", app.User.Delete).Methods(DELETE)
	r.HandleFunc("/my-profile/settings/{id}", app.User.GetMySetting).Methods(GET)
	r.HandleFunc("/my-profile/settings/{id}", app.User.SaveMySetting).Methods(PATCH)
	//r.HandleFunc("/users/search", app.User.Search).Methods(GET)

	user := "/users"
	r.HandleFunc(user+"/search", app.User.Search).Methods(GET, POST)
	r.HandleFunc(user+"/{id}", app.User.Load).Methods(GET)
	return err
}
