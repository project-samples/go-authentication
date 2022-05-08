package user

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type User struct {
	Id           string        `json:"id,omitempty" gorm:"column:id;primary_key" bson:"_id,omitempty" dynamodbav:"id,omitempty" firestore:"id,omitempty" validate:"required,max=40"`
	Username     string        `json:"username,omitempty" gorm:"column:username" bson:"username,omitempty" dynamodbav:"username,omitempty" firestore:"username,omitempty" validate:"required,username,max=100"`
	Email        string        `json:"email,omitempty" gorm:"column:email" bson:"email,omitempty" dynamodbav:"email,omitempty" firestore:"email,omitempty" validate:"email,max=100"`
	Phone        string        `json:"phone,omitempty" gorm:"column:phone" bson:"phone,omitempty" dynamodbav:"phone,omitempty" firestore:"phone,omitempty" validate:"required,phone,max=18"`
	Bio          string        `json:"bio,omitempty" gorm:"column:bio" bson:"bio,omitempty" dynamodbav:"bio,omitempty" firestore:"bio,omitempty" validate:"max=500"`
	DateOfBirth  *time.Time    `json:"dateOfBirth,omitempty" gorm:"column:date_of_birth" bson:"dateOfBirth,omitempty" dynamodbav:"dateOfBirth,omitempty" firestore:"dateOfBirth,omitempty"`
	Interests    []string      `json:"interests,omitempty" gorm:"column:interests" bson:"interests,omitempty" dynamodbav:"interests,omitempty" firestore:"interests,omitempty"`
	Skills       []Skills      `json:"skills,omitempty" gorm:"column:skills" bson:"skills,omitempty" dynamodbav:"skills,omitempty" firestore:"skills,omitempty"`
	Achievements []Achievement `json:"achievements,omitempty" gorm:"column:achievements" bson:"achievements,omitempty" dynamodbav:"achievements,omitempty" firestore:"achievements,omitempty"`
	Settings     *Settings     `json:"settings,omitempty" gorm:"column:settings" bson:"settings,omitempty" dynamodbav:"settings,omitempty" firestore:"settings,omitempty"`
}

type Skills struct {
	Skill   string `json:"skill,omitempty" gorm:"column:skill" bson:"skill,omitempty" dynamodbav:"skill,omitempty" firestore:"skill,omitempty" validate:"required"`
	Hirable bool   `json:"hirable,omitempty" gorm:"column:hirable" bson:"hirable,omitempty" dynamodbav:"hirable,omitempty" firestore:"hirable,omitempty"`
}

func (c Skills) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *Skills) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &c)
}

type Achievement struct {
	Subject     string `json:"subject,omitempty" gorm:"column:subject" bson:"subject,omitempty" dynamodbav:"subject,omitempty" firestore:"subject,omitempty" validate:"required"`
	Description string `json:"description,omitempty" gorm:"column:description" bson:"description,omitempty" dynamodbav:"description,omitempty" firestore:"description,omitempty" validate:"required"`
}

func (c Achievement) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *Achievement) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &c)
}

type Settings struct {
	UserId         string `json:"-" gorm:"-" bson:"userId,omitempty" dynamodbav:"userId,omitempty" firestore:"userId,omitempty"`
	Language       string `json:"language,omitempty" gorm:"column:language" bson:"language,omitempty" dynamodbav:"language,omitempty" firestore:"language,omitempty"`
	DateFormat     string `json:"dateFormat,omitempty" gorm:"column:date_format" bson:"dateFormat,omitempty" dynamodbav:"dateFormat,omitempty" firestore:"dateFormat,omitempty"`
	DateTimeFormat string `json:"dateTimeFormat,omitempty" gorm:"column:date_time_format" bson:"dateTimeFormat,omitempty" dynamodbav:"dateTimeFormat,omitempty" firestore:"dateTimeFormat,omitempty"`
	TimeFormat     string `json:"timeFormat,omitempty" gorm:"column:time_format" bson:"timeFormat,omitempty" dynamodbav:"timeFormat,omitempty" firestore:"timeFormat,omitempty"`
	Notification   bool   `json:"notification,omitempty" gorm:"column:notification" bson:"notification,omitempty" dynamodbav:"notification,omitempty" firestore:"notification,omitempty"`
}

func (c Settings) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *Settings) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &c)
}
