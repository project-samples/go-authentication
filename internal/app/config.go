package app

import (
	. "github.com/core-go/auth"
	. "github.com/core-go/auth/mail"
	as "github.com/core-go/auth/sql"
	. "github.com/core-go/mail/smtp"
	"github.com/core-go/mongo"
	. "github.com/core-go/oauth2"
	. "github.com/core-go/password/mail"
	"github.com/core-go/redis"
	sv "github.com/core-go/service"
	"github.com/core-go/service/builder"
	. "github.com/core-go/signup/mail"
	"github.com/core-go/sql"
)

type Config struct {
	Server   ServerConfig `mapstructure:"server"`
	Mongo    mongo.MongoConfig  `mapstructure:"mongo"`
	Location mongo.MongoConfig  `mapstructure:"location"`
	Sql      sql.Config   `mapstructure:"sql"`
	Redis    redis.Config `mapstructure:"redis"`

	MaxPasswordFailed     int                           `mapstructure:"max_password_failed"`
	LockedMinutes         int                           `mapstructure:"locked_minutes"`
	MaxPasswordAge        int32                         `mapstructure:"max_password_age"`
	PasswordEncryptionKey string                        `mapstructure:"password_encryption_key"`
	Token                 TokenConfig                   `mapstructure:"token"`
	Payload               PayloadConfig                 `mapstructure:"payload"`
	Status                *StatusConfig                 `mapstructure:"status"`
	UserStatus            UserStatusConfig              `mapstructure:"user_status"`
	Auth                  AuthMailConfig                `mapstructure:"auth"`
	AuthSqlConfig         as.SqlConfig                  `mapstructure:"auth_sql"`
	Password              PasswordMailConfig            `mapstructure:"password"`
	SignUp                SignUpConfigWithEmailTemplate `mapstructure:"sign_up"`
	OAuth2                OAuth2Config                  `mapstructure:"oauth2"`
	Mail                  MailConfig                    `mapstructure:"mail"`
	CallBackURL           CallbackURL                   `mapstructure:"callback_url"`

	Tracking   builder.TrackingConfig `mapstructure:"action"`
	ModelStatus     *sv.StatusConfig       `mapstructure:"model_status"`
	Action     *sv.ActionConfig       `mapstructure:"action"`
}

type ServerConfig struct {
	Name string `mapstructure:"name" json:"name,omitempty" gorm:"column:name" bson:"name,omitempty" dynamodbav:"name,omitempty" firestore:"name,omitempty"`
	Port *int64 `mapstructure:"port" json:"port,omitempty" gorm:"column:port" bson:"port,omitempty" dynamodbav:"port,omitempty" firestore:"port,omitempty"`
}
