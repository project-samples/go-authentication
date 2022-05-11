package app

import (
	. "github.com/core-go/auth"
	. "github.com/core-go/auth/mail"
	as "github.com/core-go/auth/sql"
	. "github.com/core-go/mail/smtp"
	. "github.com/core-go/oauth2"
	. "github.com/core-go/password/mail"
	"github.com/core-go/redis"
	sv "github.com/core-go/service"
	. "github.com/core-go/signup/mail"
)

type Root struct {
	Server                ServerConfig `mapstructure:"server"`
	Mongo                 MongoConfig  `mapstructure:"mongo"`
	Redis                 redis.Config `mapstructure:"redis"`
	StatusUser            *sv.StatusConfig
	Action                *sv.ActionConfig
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
	Template              bool                          `mapstructure:"template"`
}

type ServerConfig struct {
	Name string `mapstructure:"name" json:"name,omitempty" gorm:"column:name" bson:"name,omitempty" dynamodbav:"name,omitempty" firestore:"name,omitempty"`
	Port *int64 `mapstructure:"port" json:"port,omitempty" gorm:"column:port" bson:"port,omitempty" dynamodbav:"port,omitempty" firestore:"port,omitempty"`
}
type MongoConfig struct {
	Uri      string `mapstructure:"uri" json:"uri,omitempty" gorm:"column:uri" bson:"uri,omitempty" dynamodbav:"uri,omitempty" firestore:"uri,omitempty"`
	Database string `mapstructure:"database" json:"database,omitempty" gorm:"column:database" bson:"database,omitempty" dynamodbav:"database,omitempty" firestore:"database,omitempty"`
}
