package app

import (
	"context"
	"reflect"
	"strings"
	"time"

	. "github.com/core-go/auth"
	ac "github.com/core-go/auth/cassandra"
	h "github.com/core-go/auth/handler"
	oa2 "github.com/core-go/auth/oauth2"
	oc "github.com/core-go/auth/oauth2/cassandra"
	"github.com/core-go/cassandra"
	"github.com/core-go/cassandra/passcode"
	. "github.com/core-go/core/crypto"
	. "github.com/core-go/core/jwt"
	"github.com/core-go/core/shortid"
	. "github.com/core-go/health"
	"github.com/core-go/log"
	. "github.com/core-go/mail"
	. "github.com/core-go/mail/sendgrid"
	. "github.com/core-go/mail/smtp"
	. "github.com/core-go/password"
	pc "github.com/core-go/password/cassandra"
	"github.com/core-go/redis/v8"
	. "github.com/core-go/signup"
	sc "github.com/core-go/signup/cassandra"
	. "github.com/core-go/signup/mail"
	s "github.com/core-go/sql"
	q "github.com/core-go/sql/query"
	"github.com/gocql/gocql"
	_ "github.com/lib/pq"

	"go-service/internal/usecase/myprofile"
)

type ApplicationContext struct {
	Health         *Handler
	Authentication *h.AuthenticationHandler
	SignOut        *h.SignOutHandler
	Password       *PasswordHandler
	SignUp         *SignUpHandler
	OAuth2         *oa2.OAuth2Handler
	MyProfile      myprofile.MyProfileHandler
	Skill          *q.QueryHandler
	Interest       *q.QueryHandler
	LookingFor     *q.QueryHandler
}

func NewApp(ctx context.Context, conf Config) (*ApplicationContext, error) {
	/*
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conf.Mongo.Uri))
	if err != nil {
		return nil, err
	}
	mongoDb := client.Database(conf.Mongo.Database)
	 */
	cluster := gocql.NewCluster(conf.Cassandra.PublicIp)
	cluster.Consistency = gocql.Quorum
	cluster.ProtoVersion = 4
	cluster.Timeout = time.Second * 1000
	cluster.ConnectTimeout = time.Second * 1000
	cluster.Authenticator = gocql.PasswordAuthenticator{Username: conf.Cassandra.UserName, Password: conf.Cassandra.Password}
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	db, err := s.OpenByConfig(conf.Sql)
	if err != nil {
		return nil, err
	}
	logError := log.LogError

	generateId := shortid.Generate

	userCollection := "user"
	authentication := "authentication"

	redisPort, err := v8.NewRedisAdapterByConfig(conf.Redis)
	if err != nil {
		return nil, err
	}
	tokenBlacklistChecker := NewTokenBlacklistChecker("blacklist:", conf.Token.Expires, redisPort)

	mailService := NewMailService(conf.Mail)

	userPort := ac.NewUserRepositoryByConfig(session, userCollection, authentication, conf.SignUp.UserStatus.Activated, conf.UserStatus, conf.Auth.Schema)
	bcryptComparator := &BCryptStringComparator{}
	tokenPort := NewTokenAdapter()
	verifiedCodeSender := NewPasscodeSender(mailService, conf.Mail.From, NewTemplateLoaderByConfig(conf.Auth.Template))
	passcodeRepository := passcode.NewPasscodeRepository(cluster, "authenpasscode")
	status := InitStatus(conf.Status)
	authenticator := NewAuthenticatorWithTwoFactors(status, userPort, bcryptComparator, tokenPort.GenerateToken, conf.Token, conf.Payload, nil, verifiedCodeSender.Send, passcodeRepository, conf.Auth.Expires)
	authenticationHandler := h.NewAuthenticationHandler(authenticator.Authenticate, status.Error, status.Timeout, logError)
	signOutHandler := h.NewSignOutHandler(tokenPort.VerifyToken, conf.Token.Secret, tokenBlacklistChecker.Revoke, logError)

	passwordResetCode := "passwordResetCode"
	passwordRepository := pc.NewPasswordRepositoryByConfig(session, userCollection, authentication, userCollection, "userId", conf.Password.Schema)
	passResetCodeRepository := passcode.NewPasscodeRepository(cluster, passwordResetCode)
	p := conf.Password
	exps := []string{p.Exp1, p.Exp2, p.Exp3, p.Exp4, p.Exp5, p.Exp6}
	signupSender := NewVerifiedEmailSender(mailService, *conf.SignUp.UserVerified, conf.Mail.From, NewTemplateLoaderByConfig(*conf.SignUp.Template))
	passwordResetSender := NewPasscodeSender(mailService, conf.Mail.From, NewTemplateLoaderByConfig(conf.Password.Template.ResetTemplate))
	passwordChangeSender := NewPasscodeSender(mailService, conf.Mail.From, NewTemplateLoaderByConfig(conf.Password.Template.ChangeTemplate))
	passwordService := NewPasswordService(bcryptComparator, passwordRepository, conf.Password.ResetExpires, passResetCodeRepository, passwordResetSender.Send, tokenBlacklistChecker.RevokeAllTokens, exps, 5, nil, conf.Password.ChangeExpires, passResetCodeRepository, passwordChangeSender.Send)
	passwordHandler := NewPasswordHandler(passwordService, log.LogError, nil)

	signupCode := "signupCode"
	signupRepository := sc.NewSignUpRepositoryByConfig(session, userCollection, authentication, conf.SignUp.UserStatus, conf.MaxPasswordAge, conf.SignUp.Schema, nil)
	signupCodeRepository := passcode.NewPasscodeRepository(cluster, signupCode)
	signupStatus := InitSignUpStatus(conf.SignUp.Status)
	emailValidator := NewEmailValidator(true, "")
	signupService := NewSignUpService(signupStatus, true, signupRepository, generateId, bcryptComparator.Hash, bcryptComparator, signupCodeRepository, signupSender.Send, conf.SignUp.Expires, emailValidator.Validate, exps)
	signupHandler := NewSignUpHandler(signupService, signupStatus.Error, log.LogError, conf.SignUp.Action)

	integrationConfiguration := "integrationconfiguration"
	sources := []string{oa2.SourceGoogle, oa2.SourceFacebook, oa2.SourceLinkedIn, oa2.SourceAmazon, oa2.SourceMicrosoft, oa2.SourceDropbox}
	oauth2UserRepositories := make(map[string]oa2.OAuth2UserRepository)
	oauth2UserRepositories[oa2.SourceGoogle] = oa2.NewGoogleUserRepository()
	oauth2UserRepositories[oa2.SourceFacebook] = oa2.NewFacebookUserRepository()
	oauth2UserRepositories[oa2.SourceLinkedIn] = oa2.NewLinkedInUserRepository()
	oauth2UserRepositories[oa2.SourceAmazon] = oa2.NewAmazonUserRepository(conf.CallBackURL.Amazon)
	oauth2UserRepositories[oa2.SourceMicrosoft] = oa2.NewMicrosoftUserRepository(conf.CallBackURL.Microsoft)
	oauth2UserRepositories[oa2.SourceDropbox] = oa2.NewDropboxUserRepository()

	activatedStatus := conf.SignUp.UserStatus.Activated
	schema := conf.OAuth2.Schema
	services := strings.Split(conf.OAuth2.Services, ",")
	userRepositories := make(map[string]oa2.UserRepository)
	for _, source := range sources {
		userRepository := oc.NewUserRepositoryByConfig(session, userCollection, source, activatedStatus, services, schema, &conf.UserStatus)
		userRepositories[source] = userRepository
	}
	configurationRepository, err := oc.NewConfigurationRepository(session, integrationConfiguration, oauth2UserRepositories, "status", "A")
	if err != nil {
		return nil, err
	}

	oauth2Service := oa2.NewOAuth2Service(status, oauth2UserRepositories, userRepositories, configurationRepository, generateId, tokenPort, conf.Token, nil)
	oauth2Handler := oa2.NewDefaultOAuth2Handler(oauth2Service, status.Error, log.LogError)

	cassandraHealthChecker := cassandra.NewHealthChecker(cluster)
	redisHealthChecker := v8.NewHealthChecker(redisPort.Client)

	skillService := q.NewStringService(db, "skills", "skill")
	skillHandler := q.NewQueryHandler(skillService.Load, log.LogError)
	interestService := q.NewStringService(db, "interests", "interest")
	interestHandler := q.NewQueryHandler(interestService.Load, log.LogError)
	lookingForService := q.NewStringService(db, "searchs", "item")
	lookingForHandler := q.NewQueryHandler(lookingForService.Load, log.LogError)

	myprofileType := reflect.TypeOf(myprofile.User{})
	userRepository, err := cassandra.NewRepository(cluster, "user", myprofileType)
	myProfileService := myprofile.NewUserService(userRepository)
	myProfileHandler, err := myprofile.NewMyProfileHandler(myProfileService, log.LogError, skillService.Save, interestService.Save, lookingForService.Save)

	healthHandler := NewHandler(redisHealthChecker, cassandraHealthChecker)

	appCtx := ApplicationContext{
		Health:         healthHandler,
		Authentication: authenticationHandler,
		SignOut:        signOutHandler,
		Password:       passwordHandler,
		SignUp:         signupHandler,
		OAuth2:         oauth2Handler,
		MyProfile:      myProfileHandler,
		Skill:          skillHandler,
		Interest:       interestHandler,
		LookingFor:     lookingForHandler,
	}
	return &appCtx, nil
}

func NewMailService(mailConfig MailConfig) SimpleMailSender {
	if mailConfig.Provider == "sendgrid" {
		return NewSimpleMailSender(NewSendGridMailSender(mailConfig.APIkey))
	}
	return NewSimpleMailSender(NewSmtpMailSender(mailConfig.SMTP))
}
