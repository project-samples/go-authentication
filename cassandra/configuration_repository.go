package cassandra

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/core-go/oauth2"
	"github.com/gocql/gocql"
)

type ConfigurationRepository struct {
	Cassandra              *gocql.ClusterConfig
	TableName              string
	OAuth2UserRepositories map[string]oauth2.OAuth2UserRepository
	Status                 string
	Active                 string
	Driver                 string
}

func NewConfigurationRepository(db *gocql.ClusterConfig, tableName string, oauth2UserRepositories map[string]oauth2.OAuth2UserRepository, status string, active string) *ConfigurationRepository {
	if len(status) == 0 {
		status = "status"
	}
	if len(active) == 0 {
		active = "A"
	}
	// collection := db.Collection(collectionName)
	return &ConfigurationRepository{Cassandra: db, TableName: tableName, OAuth2UserRepositories: oauth2UserRepositories, Status: status, Active: active}
}

func (s *ConfigurationRepository) GetConfiguration(ctx context.Context, id string) (*oauth2.Configuration, string, error) {
	session, er0 := s.Cassandra.CreateSession()
	if er0 != nil {
		return nil, "", er0
	}
	model := oauth2.Configuration{}
	query := fmt.Sprintf(`select %s from %s where %s = ? ALLOW FILTERING`, "clientid, clientsecret ", s.TableName, "sourcetype")
	log.Println(query)
	rows := session.Query(query, id)
	if rows.Exec() != nil {
		return nil, "", rows.Exec()

	}
	err2 := rows.Scan(&model.ClientId, &model.ClientSecret)
	if err2 != nil {
		return nil, "", err2

	}
	defer session.Close()
	clientId := model.ClientId
	clientId, err := s.OAuth2UserRepositories[id].GetRequestTokenOAuth(ctx, model.ClientId, model.ClientSecret)
	return &model, clientId, err
}
func (s *ConfigurationRepository) GetConfigurations(ctx context.Context) (*[]oauth2.Configuration, error) {
	session, er0 := s.Cassandra.CreateSession()
	if er0 != nil {
		return nil, er0
	}
	query := fmt.Sprintf(`select * from %s where %s = ? `, s.TableName, s.Status)
	rows := session.Query(query, s.Active)
	if rows.Exec() != nil {
		return nil, rows.Exec()
	}
	model := oauth2.Configuration{}
	models := make([]oauth2.Configuration, 0)
	modelType := reflect.TypeOf(model)
	fieldsIndex, er1 := getColumnIndexes(modelType)
	if er1 != nil {
		return nil, er1
	}
	defer session.Close()
	err1 := rows.Scan(rows, &models, fieldsIndex)
	if err1 != nil {
		return nil, err1
	}
	return &models, nil
}

func getColumnIndexes(modelType reflect.Type) (map[string]int, error) {
	mapp := make(map[string]int, 0)
	if modelType.Kind() != reflect.Struct {
		return mapp, errors.New("bad type")
	}
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		ormTag := field.Tag.Get("gorm")
		column, ok := findTag(ormTag, "column")
		if ok {
			mapp[column] = i
		}
	}
	return mapp, nil
}

func findTag(tag string, key string) (string, bool) {
	if has := strings.Contains(tag, key); has {
		str1 := strings.Split(tag, ";")
		num := len(str1)
		for i := 0; i < num; i++ {
			str2 := strings.Split(str1[i], ":")
			for j := 0; j < len(str2); j++ {
				if str2[j] == key {
					return str2[j+1], true
				}
			}
		}
	}
	return "", false
}
