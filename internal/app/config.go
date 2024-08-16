package app

import (
	. "github.com/core-go/authentication"
	. "github.com/core-go/authentication/mail"
	o2 "github.com/core-go/authentication/oauth2"
	. "github.com/core-go/mail/smtp"
	. "github.com/core-go/password/mail"
	"github.com/core-go/redis/v8"
	. "github.com/core-go/signup/mail"
	"github.com/core-go/sql"
)

type Config struct {
	Server   ServerConfig `mapstructure:"server"`
	Mongo    MongoConfig  `mapstructure:"mongo"`
	Location MongoConfig  `mapstructure:"location"`
	Sql      sql.Config   `mapstructure:"sql"`
	Redis    v8.Config    `mapstructure:"redis"`

	MaxPasswordFailed     int                           `mapstructure:"max_password_failed"`
	LockedMinutes         int                           `mapstructure:"locked_minutes"`
	MaxPasswordAge        int32                         `mapstructure:"max_password_age"`
	PasswordEncryptionKey string                        `mapstructure:"password_encryption_key"`
	Token                 TokenConfig                   `mapstructure:"token"`
	Payload               PayloadConfig                 `mapstructure:"payload"`
	Status                *StatusConfig                 `mapstructure:"status"`
	UserStatus            UserStatusConfig              `mapstructure:"user_status"`
	Auth                  AuthMailConfig                `mapstructure:"auth"`
	Password              PasswordMailConfig            `mapstructure:"password"`
	SignUp                SignUpConfigWithEmailTemplate `mapstructure:"sign_up"`
	OAuth2                o2.OAuth2Config               `mapstructure:"oauth2"`
	Mail                  MailConfig                    `mapstructure:"mail"`
	CallBackURL           o2.CallbackURL                `mapstructure:"callback_url"`
}

type ServerConfig struct {
	Name string `mapstructure:"name" json:"name,omitempty" gorm:"column:name" bson:"name,omitempty" dynamodbav:"name,omitempty" firestore:"name,omitempty"`
	Port *int64 `mapstructure:"port" json:"port,omitempty" gorm:"column:port" bson:"port,omitempty" dynamodbav:"port,omitempty" firestore:"port,omitempty"`
}

type MongoConfig struct {
	Uri      string `yaml:"uri" mapstructure:"uri" json:"uri,omitempty" gorm:"column:uri" bson:"uri,omitempty" dynamodbav:"uri,omitempty" firestore:"uri,omitempty"`
	Database string `yaml:"database" mapstructure:"database" json:"database,omitempty" gorm:"column:database" bson:"database,omitempty" dynamodbav:"database,omitempty" firestore:"database,omitempty"`
}
