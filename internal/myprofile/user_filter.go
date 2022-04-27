package myprofile

import "github.com/core-go/search"

type UserFilter struct {
	Id           string            `json:"id" gorm:"column:id;primary_key" bson:"_id" dynamodbav:"id" firestore:"id" match:"equal" validate:"required,max=40"`
	Username     string            `json:"username" gorm:"column:username" bson:"username" dynamodbav:"username" firestore:"username" match:"prefix" validate:"required,username,max=100"`
	Email        string            `json:"email" gorm:"column:email" bson:"email" dynamodbav:"email" firestore:"email" match:"prefix" validate:"email,max=100" q:"prefix"`
	Phone        string            `json:"phone" gorm:"column:phone" bson:"phone" dynamodbav:"phone" firestore:"phone" validate:"required,phone,max=18" q:"prefix"`
	DateOfBirth  *search.TimeRange `json:"dateOfBirth" gorm:"column:date_of_birth" bson:"dateOfBirth" dynamodbav:"dateOfBirth" firestore:"dateOfBirth"`
	Interests    []string          `json:"interests" gorm:"column:interests" bson:"interests" dynamodbav:"interests" firestore:"interests" validate:""`
	Skills       []Skills          `json:"skills" gorm:"column:skills" bson:"skills" dynamodbav:"skills" firestore:"skills" validate:""`
	Achievements []Achievement     `json:"achievements" gorm:"column:achievements" bson:"achievements" dynamodbav:"achievements" firestore:"achievements" validate:""`
	Settings     *Settings         `json:"settings" gorm:"column:settings" bson:"settings" dynamodbav:"settings" firestore:"settings" validate:""`
}
