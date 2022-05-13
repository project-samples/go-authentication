package myprofile

import (
	"fmt"
	"github.com/core-go/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
)

var user User

func BuildQuery(param interface{}) (bson.D, bson.M) {
	userType := reflect.TypeOf(user)
	var query = bson.D{}
	var fields = bson.M{}
	filter, ok := param.(*UserFilter)
	if !ok {
		return query, fields
	}
	if len(filter.Fields) > 0 {
		for _, key := range filter.Fields {
			_, _, columnName := mongo.GetFieldByJson(userType, key)
			if len(columnName) < 0 {
				fields = bson.M{}
				break
			}
			fields[columnName] = 1
		}
	}
	if len(filter.Username) > 0 {
		query = append(query, bson.E{Key: "username", Value: primitive.Regex{Pattern: fmt.Sprintf("^%v", filter.Username)}})
	}
	if len(filter.Email) > 0 {
		query = append(query, bson.E{Key: "email", Value: primitive.Regex{Pattern: fmt.Sprintf("^%v", filter.Email)}})
	}
	if len(filter.Phone) > 0 {
		query = append(query, bson.E{Key: "phone", Value: primitive.Regex{Pattern: fmt.Sprintf("\\w*%v\\w*", filter.Phone)}})
	}
	if filter.Interests != nil && len(filter.Interests) > 0 {
		for _, value := range filter.Interests {
			query = append(query, bson.E{Key: "interests", Value: primitive.Regex{Pattern: fmt.Sprintf("^%v", value)}})
		}
		//query = append(query, bson.E{Key: "interests", Value: primitive.Regex{Pattern: fmt.Sprintf("\\w*%v\\w*", filter.Interests), Options: "i"}})
	}
	if len(filter.FirstName) > 0 {
		query = append(query, bson.E{Key: "firstName", Value: primitive.Regex{Pattern: fmt.Sprintf("\\w*%v\\w*", filter.FirstName)}})
	}
	return query, fields
}
