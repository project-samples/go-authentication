package user

import "github.com/core-go/search"

type UserFilter struct {
	*search.Filter
	Id           string            `json:"id" gorm:"column:id;primary_key" bson:"_id" dynamodbav:"id" firestore:"id" match:"equal" validate:"required,max=40"`
	Username     string            `json:"username" gorm:"column:username" bson:"username" dynamodbav:"username" firestore:"username" match:"prefix" validate:"required,username,max=100"`
	Email        string            `json:"email" gorm:"column:email" bson:"email" dynamodbav:"email" firestore:"email" match:"prefix" validate:"email,max=100" q:"prefix"`
	Phone        string            `json:"phone" gorm:"column:phone" bson:"phone" dynamodbav:"phone" firestore:"phone" validate:"required,phone,max=18" q:"prefix"`
	FirstName    string            `json:"firstName" gorm:"column:firstName" bson:"firstName" dynamodbav:"firstName" q:"prefix" firestore:"firstName" validate:"required,firstName,max=100"`
	LastName     string            `json:"last_name" gorm:"column:last_name" bson:"last_name" dynamodbav:"last_name" firestore:"last_name" q:"prefix" validate:"required,last_name,max=100"`
	Occupation   string            `json:"occupation" gorm:"column:occupation" bson:"occupation" dynamodbav:"occupation" q:"prefix" firestore:"occupation" validate:"required,occupation,max=100"`
	Company      string            `json:"company" gorm:"column:company" bson:"company" dynamodbav:"company" firestore:"company" validate:"required,company,max=100"`
	DateOfBirth  *search.TimeRange `json:"dateOfBirth" gorm:"column:date_of_birth" bson:"dateOfBirth" dynamodbav:"dateOfBirth" firestore:"dateOfBirth"`
	Interests    []string          `json:"interests" gorm:"column:interests" bson:"interests" dynamodbav:"interests" firestore:"interests" validate:""`
	Skills       []Skills          `json:"skills" gorm:"column:skills" bson:"skills" dynamodbav:"skills" firestore:"skills" validate:""`
	Achievements []Achievement     `json:"achievements" gorm:"column:achievements" bson:"achievements" dynamodbav:"achievements" firestore:"achievements" validate:""`
}
