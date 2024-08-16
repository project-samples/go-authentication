package user

import (
	"context"
	q "github.com/core-go/mongo/query"
	"github.com/core-go/search"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewUserTransport(db *mongo.Database, logError func(context.Context, string, ...map[string]interface{})) UserTransport {
	query := q.NewQuery[User, string, *UserFilter](db, "user", BuildQuery, search.GetSort)
	return NewUserHandler(query, logError)
}
