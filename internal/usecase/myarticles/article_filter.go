package myarticles

import (
	"github.com/core-go/search"
	"time"
)

type ArticleFilter struct {
	*search.Filter
	Id          string     `json:"id" gorm:"primary_key;column:id" bson:"_id" dynamodbav:"id,omitempty" firestore:"-"`
	Title       string     `json:"title,omitempty" gorm:"column:title" bson:"title,omitempty" dynamodbav:"title,omitempty" firestore:"title,omitempty"`
	Description string     `json:"description,omitempty" gorm:"column:description" bson:"description" dynamodbav:"description,omitempty" firestore:"description,omitempty"`
	PublishedAt *time.Time `json:"publishedAt,omitempty" gorm:"column:publishedAt" bson:"publishedAt,omitempty" dynamodbav:"publishedAt,omitempty" firestore:"publishedAt,omitempty"`
	Content     string     `json:"content,omitempty" gorm:"column:content" bson:"content,omitempty" dynamodbav:"content,omitempty" firestore:"content,omitempty"`
	Tags        []string   `json:"tags,omitempty" gorm:"column:tags" bson:"tags,omitempty" dynamodbav:"tags,omitempty" firestore:"tags,omitempty"`
	Status      string     `json:"status,omitempty" gorm:"column:status" bson:"status" dynamodbav:"status,omitempty" firestore:"status,omitempty"`
}
