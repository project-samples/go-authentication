package app

import (
	"context"
	"reflect"
	"strings"

	. "github.com/core-go/auth"
	am "github.com/core-go/auth/mongo"
	. "github.com/core-go/health"
	"github.com/core-go/log"
	. "github.com/core-go/mail"
	. "github.com/core-go/mail/sendgrid"
	. "github.com/core-go/mail/smtp"
	mgo "github.com/core-go/mongo"
	. "github.com/core-go/oauth2"
	om "github.com/core-go/oauth2/mongo"
	. "github.com/core-go/password"
	pm "github.com/core-go/password/mongo"
	"github.com/core-go/redis"
	"github.com/core-go/search"
	. "github.com/core-go/security/crypto"
	. "github.com/core-go/security/jwt"
	"github.com/core-go/service/shortid"
	. "github.com/core-go/signup"
	. "github.com/core-go/signup/mail"
	sm "github.com/core-go/signup/mongo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	usr "go-service/internal/usecase/user"
)

type ApplicationContext struct {
	Health         *Handler
	Authentication *AuthenticationHandler
	SignOut        *SignOutHandler
	Password       *PasswordHandler
	SignUp         *SignUpHandler
	OAuth2         *OAuth2Handler
	User           usr.UserHandler
}

func NewApp(ctx context.Context, root Root) (*ApplicationContext, error) {
	client, er1 := mongo.Connect(ctx, options.Client().ApplyURI(root.Mongo.Uri))
	if er1 != nil {
		return nil, er1
	}
	mongoDb := client.Database(root.Mongo.Database)

	logError := log.ErrorMsg

	generateId := shortid.Generate

	user := "user"
	authentication := "authentication"

	redisService, er2 := redis.NewRedisServiceByConfig(root.Redis)
	if er2 != nil {
		return nil, er2
	}
	tokenBlacklistChecker := NewTokenBlacklistChecker("blacklist:", root.Token.Expires, redisService)

	mailService := NewMailService(root.Mail)

	authenticationRepository := am.NewAuthenticationRepositoryByConfig(mongoDb, user, authentication, root.SignUp.UserStatus.Activated, root.UserStatus, root.Auth.Schema)
	userInfoService := NewUserInfoService(authenticationRepository, root.MaxPasswordAge, root.MaxPasswordFailed, root.LockedMinutes)
	bcryptComparator := &BCryptStringComparator{}
	tokenService := NewTokenService()
	verifiedCodeSender := NewPasscodeSender(mailService, root.Mail.From, NewTemplateLoaderByConfig(root.Auth.Template))
	passCodeService := mgo.NewPasscodeRepository(mongoDb, "authenpasscode")
	status := InitStatus(root.Status)
	authenticator := NewAuthenticatorWithTwoFactors(status, userInfoService, bcryptComparator, tokenService.GenerateToken, root.Token, root.Payload, nil, verifiedCodeSender.Send, passCodeService, root.Auth.Expires)
	authenticationHandler := NewAuthenticationHandler(authenticator.Authenticate, status.Error, status.Timeout, logError)
	signOutHandler := NewSignOutHandler(tokenService.VerifyToken, root.Token.Secret, tokenBlacklistChecker.Revoke, logError)

	passwordResetCode := "passwordResetCode"
	passwordRepository := pm.NewPasswordRepositoryByConfig(mongoDb, user, authentication, user, "userId", root.Password.Schema)
	passResetCodeRepository := mgo.NewPasscodeRepository(mongoDb, passwordResetCode)
	p := root.Password
	exps := []string{p.Exp1, p.Exp2, p.Exp3, p.Exp4, p.Exp5, p.Exp6}
	signupSender := NewVerifiedEmailSender(mailService, *root.SignUp.UserVerified, root.Mail.From, NewTemplateLoaderByConfig(*root.SignUp.Template))
	passwordResetSender := NewPasscodeSender(mailService, root.Mail.From, NewTemplateLoaderByConfig(root.Password.Template.ResetTemplate))
	passwordChangeSender := NewPasscodeSender(mailService, root.Mail.From, NewTemplateLoaderByConfig(root.Password.Template.ChangeTemplate))
	passwordService := NewPasswordService(bcryptComparator, passwordRepository, root.Password.ResetExpires, passResetCodeRepository, passwordResetSender.Send, tokenBlacklistChecker.RevokeAllTokens, exps, 5, nil, root.Password.ChangeExpires, passResetCodeRepository, passwordChangeSender.Send)
	passwordHandler := NewPasswordHandler(passwordService, logError, nil)

	signUpCode := "signupCode"
	signUpRepository := sm.NewSignUpRepositoryByConfig(mongoDb, user, authentication, root.SignUp.UserStatus, root.MaxPasswordAge, root.SignUp.Schema, nil)
	signUpCodeRepository := mgo.NewPasscodeRepository(mongoDb, signUpCode)
	signupStatus := InitSignUpStatus(root.SignUp.Status)
	emailValidator := NewEmailValidator(true, "")
	signUpService := NewSignUpService(signupStatus, true, signUpRepository, generateId, bcryptComparator.Hash, bcryptComparator, signUpCodeRepository, signupSender.Send, root.SignUp.Expires, emailValidator.Validate, exps)
	signupHandler := NewSignUpHandler(signUpService, signupStatus.Error, logError, root.SignUp.Action)

	integrationConfiguration := "integrationconfiguration"
	sources := []string{SourceGoogle, SourceFacebook, SourceLinkedIn, SourceAmazon, SourceMicrosoft, SourceDropbox}
	oauth2UserRepositories := make(map[string]OAuth2UserRepository)
	oauth2UserRepositories[SourceGoogle] = NewGoogleUserRepository()
	oauth2UserRepositories[SourceFacebook] = NewFacebookUserRepository()
	oauth2UserRepositories[SourceLinkedIn] = NewLinkedInUserRepository()
	oauth2UserRepositories[SourceAmazon] = NewAmazonUserRepository(root.CallBackURL.Amazon)
	oauth2UserRepositories[SourceMicrosoft] = NewMicrosoftUserRepository(root.CallBackURL.Microsoft)
	oauth2UserRepositories[SourceDropbox] = NewDropboxUserRepository()

	activatedStatus := root.SignUp.UserStatus.Activated
	schema := root.OAuth2.Schema
	services := strings.Split(root.OAuth2.Services, ",")
	userRepositories := make(map[string]UserRepository)
	for _, source := range sources {
		userRepository := om.NewUserRepositoryByConfig(mongoDb, user, source, activatedStatus, services, schema, &root.UserStatus)
		userRepositories[source] = userRepository
	}
	configurationRepository := om.NewConfigurationRepository(mongoDb, integrationConfiguration, oauth2UserRepositories, "status", "A")

	oauth2Service := NewOAuth2Service(status, oauth2UserRepositories, userRepositories, configurationRepository, generateId, tokenService, root.Token, nil)
	oauth2Handler := NewDefaultOAuth2Handler(oauth2Service, status.Error, logError)

	mongoHealthChecker := mgo.NewHealthChecker(mongoDb)
	redisHealthChecker := redis.NewHealthChecker(redisService.Pool)

	userType := reflect.TypeOf(usr.User{})
	userSearchBuilder := mgo.NewSearchBuilder(mongoDb, "user", usr.BuildQuery, search.GetSort)
	getUser := mgo.UseGet(mongoDb, "user", userType)
	userHandler := usr.NewUserHandler(userSearchBuilder.Search, getUser, logError, nil)

	healthHandler := NewHandler(redisHealthChecker, mongoHealthChecker)

	app := ApplicationContext{
		Health:         healthHandler,
		Authentication: authenticationHandler,
		SignOut:        signOutHandler,
		Password:       passwordHandler,
		SignUp:         signupHandler,
		OAuth2:         oauth2Handler,
		User:           userHandler,
	}
	return &app, nil
}

func NewMailService(mailConfig MailConfig) SimpleMailSender {
	if mailConfig.Provider == "sendgrid" {
		return NewSimpleMailSender(NewSendGridMailSender(mailConfig.APIkey))
	}
	return NewSimpleMailSender(NewSmtpMailSender(mailConfig.SMTP))
}
